package queries

import (
	my "compare-tables/yaml"
	"database/sql"
	"fmt"
	"log"
	"strings"

	// sql drivers
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// sql drivers
const (
	driverPsql  = "postgres"
	driverMysql = "mysql"
)

// Instance struct
type Instance struct {
	DB   *sql.DB
	Data []*Table
}

// GetInstance do first
func GetInstance(env my.Env) *Instance {
	src := getSrcName(env)

	db, err := sql.Open(env.Driver, src)
	if err != nil {
		panic(
			fmt.Sprintf("coannot open %s: %v", env.Driver, err),
		)
	}
	return &Instance{
		DB: db,
	}
}

// Close set defer
func (ins *Instance) Close() {
	ins.DB.Close()
}

// Init init instance
func (ins *Instance) Init(data []*my.Table) {
	err := ins.DB.Ping()
	if err != nil {
		panic(
			fmt.Sprintf("cannot connect db: %v", err),
		)
	}
	log.Println("db connection ok.")

	ins.Data = getData(data)
}

// RunCompare do compare
func (ins *Instance) RunCompare() {
	matchAll, result := true, ""

	for i, v := range ins.Data {
		fmt.Printf("---- %s : %s ----\n", v.Origin, v.Diff)
		cntOrigin := ins.getCountOrigin(i)
		cntDiff := ins.getCountDiff(i)
		joinData := ins.getInnerJoin(i)
		cntJoin := len(joinData)
		matchAllData := ins.getInnerJoinAll(i)
		cntMatchAll := len(matchAllData)
		isMatch := cntOrigin > 0 && cntDiff > 0 && cntJoin > 0 && cntMatchAll > 0
		result += fmt.Sprintf(
			"%s=%s : %v\n",
			v.Origin,
			v.Diff,
			isMatch,
		)
		if matchAll {
			matchAll = isMatch
		}
	}
	result += fmt.Sprintf(
		"match all : %v\n",
		matchAll,
	)
	fmt.Print("---- Result ----\n", result)
}

func (ins *Instance) getCountOrigin(i int) int {
	table := ins.Data[i]
	fmt.Println("count", table.Origin)

	sql := fmt.Sprintf(
		"SELECT COUNT(*) FROM %s WHERE %s",
		table.Origin,
		table.Where.Origin,
	)
	fmt.Printf(
		"\t[SQL] %s: %s\n",
		table.Origin,
		sql,
	)

	var count int
	if err := ins.DB.QueryRow(sql).Scan(&count); err != nil {
		panic(err)
	}
	fmt.Printf(
		"\t[COUNT]: %d\n",
		count,
	)
	return count
}

func (ins *Instance) getCountDiff(i int) int {
	table := ins.Data[i]
	fmt.Println("count", table.Diff)

	sql := fmt.Sprintf(
		"SELECT COUNT(*) FROM %s WHERE %s",
		table.Diff,
		table.Where.Diff,
	)
	fmt.Printf(
		"\t[SQL] %s: %s\n",
		table.Diff,
		sql,
	)

	var count int
	if err := ins.DB.QueryRow(sql).Scan(&count); err != nil {
		panic(err)
	}
	fmt.Printf(
		"\t[COUNT]: %d\n",
		count,
	)
	return count
}

func (ins *Instance) getInnerJoin(i int) [][]sql.NullString {
	table := ins.Data[i]
	fmt.Println("inner join")

	cols := make([]string, len(table.Columns)*2)
	for i := range table.Columns {
		col := table.Columns[i]

		i *= 2
		cols[i] = fmt.Sprintf("%s.%s", table.Origin, col.Origin)
		cols[i+1] = fmt.Sprintf("%s.%s", table.Diff, col.Diff)
	}

	strSql := fmt.Sprintf(
		"SELECT %s FROM %s INNER JOIN %s ON %s WHERE %s",
		strings.Join(cols, ", "),
		table.Origin,
		table.Diff,
		table.JoinOn.Origin,
		table.Where.Origin,
	)
	fmt.Printf("\t[SQL]: %s\n", strSql)

	rows, err := ins.DB.Query(strSql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	data := [][]sql.NullString{}

	for rows.Next() {
		results := make([]sql.NullString, len(cols))
		pipe := make([]interface{}, len(cols))
		for i := range results {
			pipe[i] = &results[i]
		}
		err = rows.Scan(pipe...)
		if err != nil {
			panic(err)
		}
		data = append(data, results)
	}
	fmt.Printf("\t[COUNT]: %d\n", len(data))
	return data
}

func (ins *Instance) getInnerJoinAll(i int) [][]sql.NullString {
	table := ins.Data[i]
	fmt.Println("inner join column match all")

	cols := make([]string, len(table.Columns)*2)
	for i := range table.Columns {
		col := table.Columns[i]

		i *= 2
		cols[i] = fmt.Sprintf("%s.%s", table.Origin, col.Origin)
		cols[i+1] = fmt.Sprintf("%s.%s", table.Diff, col.Diff)
	}

	strSql := fmt.Sprintf(
		"SELECT %s FROM %s INNER JOIN %s ON %s WHERE %s",
		strings.Join(cols, ", "),
		table.Origin,
		table.Diff,
		table.JoinOn.Origin,
		table.Where.Origin,
	)

	for _, v := range table.Columns {
		if 0 < len(table.Where.Origin) {
			strSql += " AND "
		}
		origin := fmt.Sprintf("%s.%s", table.Origin, v.Origin)
		diff := fmt.Sprintf("%s.%s", table.Diff, v.Diff)
		strSql += fmt.Sprintf(
			"((%s IS NULL AND %s IS NULL) OR %s = %s)",
			origin,
			diff,
			origin,
			diff,
		)
	}
	fmt.Printf("\t[SQL]: %s\n", strSql)

	rows, err := ins.DB.Query(strSql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	data := [][]sql.NullString{}

	for rows.Next() {
		results := make([]sql.NullString, len(cols))
		pipe := make([]interface{}, len(cols))
		for i := range results {
			pipe[i] = &results[i]
		}
		err = rows.Scan(pipe...)
		if err != nil {
			panic(err)
		}
		data = append(data, results)
	}
	fmt.Printf("\t[COUNT]: %d\n", len(data))

	return data
}
