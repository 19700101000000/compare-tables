package queries

import (
	my "compare-tables/yaml"
	"database/sql"
	"fmt"
	"log"

	// sql drivers
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// sql drivers
const (
	driverPsql  = "postgres"
	driverMysql = "mysql"
)

// tag heads
const (
	tagSQL = "[SQL]"
	tagCnt = "[CNT]"
	tagAll = "[ALL]"
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
	matchAll, resultsMatch := true, ""

	for i, v := range ins.Data {
		fmt.Printf("----[%s:%s]----\n", v.Origin, v.Diff)

		cntOrigin := ins.getCountOrigin(i)
		cntDiff := ins.getCountDiff(i)
		joinData := ins.getInnerJoin(i)
		matchData := ins.getInnerJoinWithMatch(i)

		cntJoin := len(joinData)
		cntMatch := len(matchData)
		isMatchAll := cntOrigin == cntDiff && cntOrigin == cntJoin && cntOrigin == cntMatch
		resultsMatch += fmt.Sprintf(
			"\t[%s:%s]\t%v\n",
			v.Origin,
			v.Diff,
			isMatchAll,
		)
		if matchAll {
			matchAll = isMatchAll
		}
	}
	resultsMatch += fmt.Sprintf("\t%s\t%v\n", tagAll, matchAll)

	fmt.Println("----[results]----")
	fmt.Printf("match\n%s", resultsMatch)
}

func (ins *Instance) getCountOrigin(i int) int {
	t := ins.Data[i]
	fmt.Println("count", t.Origin)

	q := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", t.Origin, t.Where.Origin)
	fmt.Printf("\t%s\t%s: %s\n", tagSQL, t.Origin, q)

	var c int
	if err := ins.DB.QueryRow(q).Scan(&c); err != nil {
		panic(err)
	}
	fmt.Printf("\t%s\t%d\n", tagCnt, c)
	return c
}

func (ins *Instance) getCountDiff(i int) int {
	t := ins.Data[i]
	fmt.Println("count", t.Diff)

	q := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", t.Diff, t.Where.Diff)
	fmt.Printf("\t%s\t%s: %s\n", tagSQL, t.Diff, q)

	var c int
	if err := ins.DB.QueryRow(q).Scan(&c); err != nil {
		panic(err)
	}
	fmt.Printf("\t%s\t%d\n", tagCnt, c)
	return c
}

func (ins *Instance) getInnerJoin(i int) [][]sql.NullString {
	t := ins.Data[i]
	fmt.Println("inner join")

	q := getInnerJoinQuery(ins, i)
	fmt.Printf("\t%s\t%s\n", tagSQL, q)

	rows, err := ins.DB.Query(q)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	d := [][]sql.NullString{}
	for rows.Next() {
		l := len(t.Columns) * 2
		r := make([]sql.NullString, l)
		p := make([]interface{}, l)
		for i := range r {
			p[i] = &r[i]
		}
		err = rows.Scan(p...)
		if err != nil {
			panic(err)
		}
		d = append(d, r)
	}
	fmt.Printf("\t%s\t%d\n", tagCnt, len(d))
	return d
}

func (ins *Instance) getInnerJoinWithMatch(i int) [][]sql.NullString {
	t := ins.Data[i]
	fmt.Println("inner join with match")

	q := getInnerJoinQuery(ins, i)
	for _, v := range t.Columns {
		if v.DisableMatch {
			continue
		}

		if 0 < len(t.Where.Origin) {
			q += " AND "
		}
		s := "%s.%s"
		o := fmt.Sprintf(s, t.Origin, v.Origin)
		d := fmt.Sprintf(s, t.Diff, v.Diff)
		q += fmt.Sprintf("((%s IS NULL AND %s IS NULL) OR %s = %s)", o, d, o, d)
	}
	fmt.Printf("\t%s\t%s\n", tagSQL, q)

	rows, err := ins.DB.Query(q)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	d := [][]sql.NullString{}
	for rows.Next() {
		l := len(t.Columns) * 2
		r := make([]sql.NullString, l)
		p := make([]interface{}, l)
		for i := range r {
			p[i] = &r[i]
		}
		err = rows.Scan(p...)
		if err != nil {
			panic(err)
		}
		d = append(d, r)
	}
	fmt.Printf("\t%s\t%d\n", tagCnt, len(d))
	return d
}
