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
	for i, v := range ins.Data {
		fmt.Printf("---- %s : %s ----\n", v.Origin, v.Diff)
		ins.getCountOrigin(i)
		ins.getCountDiff(i)
		ins.getInnerJoin(i)
	}
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
		"\t[COUNT] %s: %d\n",
		table.Origin,
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
		"\t[COUNT] %s: %d\n",
		table.Diff,
		count,
	)
	return count
}

func (ins *Instance) getInnerJoin(i int) {
	table := ins.Data[i]
	fmt.Println("inner join")

	cols := make([]string, len(table.Columns)*2)
	for i := range table.Columns {
		col := table.Columns[i]

		i *= 2
		cols[i] = fmt.Sprintf("%s.%s", table.Origin, col.Origin)
		cols[i+1] = fmt.Sprintf("%s.%s", table.Diff, col.Diff)
	}

	sql := fmt.Sprintf(
		"SELECT %s FROM %s INNER JOIN %s ON %s WHERE %s",
		strings.Join(cols, ", "),
		table.Origin,
		table.Diff,
		table.JoinOn.Origin,
		table.Where.Origin,
	)
	fmt.Printf("\t[SQL]: %s\n", sql)

	rows, err := ins.DB.Query(sql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	data := [][]string{}

	for rows.Next() {
		results := make([]string, len(cols))
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

	fmt.Println("inner join column match all")
	for _, v := range table.Columns {
		if 0 < len(table.Where.Origin) {
			sql += " AND "
		}
		sql += fmt.Sprintf("%s.%s = %s.%s", table.Origin, v.Origin, table.Diff, v.Diff)
	}
	fmt.Printf("\t[SQL]: %s\n", sql)

	rows, err = ins.DB.Query(sql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	data = [][]string{}

	for rows.Next() {
		results := make([]string, len(cols))
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
}
