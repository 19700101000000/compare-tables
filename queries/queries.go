package queries

import (
	"compare-tables/types"
	my "compare-tables/yaml"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	// sql drivers
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

const (
	driverPsql  types.Driver = "postgres"
	driverMysql types.Driver = "mysql"
)

// tag heads
const (
	tagSQL = "[SQL]"
	tagCnt = "[CNT]"
	tagAll = "[ALL]"
)

// Instance struct
type Instance struct {
	DB     *sql.DB
	DiffDB *sql.DB
	Data   []*Query
}

// QueryResult in query result
type QueryResult struct {
	Header []string
	Data   [][]sql.NullString
}

// CheckResult check queries
type CheckResult struct {
	Match      map[int][]int
	IsMatchAll bool
}

// GetInstance do first
func GetInstance(env my.Env) *Instance {
	o, d := getEnvs(env)

	src := getSrcName(o)
	db, err := sql.Open(string(o.Driver), src)
	if err != nil {
		panic(
			fmt.Sprintf("coannot open %s: %v", o.Driver, err),
		)
	}
	var diffDb *sql.DB
	if env.NotSame {
		src = getSrcName(d)
		diffDb, err = sql.Open(string(d.Driver), src)
		if err != nil {
			panic(
				fmt.Sprintf("coannot open %s: %v", d.Driver, err),
			)
		}
	}
	return &Instance{
		DB:     db,
		DiffDB: diffDb,
	}
}

// Close set defer
func (ins *Instance) Close() {
	if ins == nil {
		return
	}

	if ins.DB != nil {
		ins.DB.Close()
	}
	if ins.DiffDB != nil {
		ins.DiffDB.Close()
	}
}

// Init init instance
func (ins *Instance) Init(data []*my.Table) {
	if ins == nil {
		return
	}

	if ins.DB != nil {
		if err := ins.DB.Ping(); err != nil {
			panic(
				fmt.Sprintf("cannot connect 1st db: %v", err),
			)
		}
		log.Printf("1st db connection ok.")
	}

	if ins.DiffDB != nil {
		if err := ins.DiffDB.Ping(); err != nil {
			panic(
				fmt.Sprintf("cannot connect 2nd db: %v", err),
			)
		}
		log.Printf("2nd db connection ok.")
	}

	ins.Data = getData(data)
}

// RunCompare do compare
func (ins *Instance) RunCompare() {
	if ins == nil {
		return
	}
	isMatchAll, results := true, ""

	fmt.Printf("[%s]\n", time.Now())
	for i, v := range ins.Data {
		fmt.Printf("----[%s:%s]----\n", v.Origin, v.Diff)

		if ins.DiffDB == nil {
			cntOrigin := ins.getCountOrigin(i)
			cntDiff := ins.getCountDiff(i)
			joinData := ins.getInnerJoin(i)
			matchData := ins.getInnerJoinWithMatch(i)

			cntJoin := len(joinData)
			cntMatch := len(matchData)
			isMatch := cntOrigin == cntDiff && cntOrigin == cntJoin && cntOrigin == cntMatch
			results += fmt.Sprintf(
				"\t[%s:%s]\t%v\n",
				v.Origin,
				v.Diff,
				isMatch,
			)
			if isMatchAll {
				isMatchAll = isMatch
			}
		} else {
			o := ins.findOrigin(i)
			d := ins.findDiff(i)
			r := ins.checkResults(&o, &d, i)
			results += fmt.Sprintf("%s:%s\n", v.Origin, v.Diff)
			c := 0
			for k := range r.Match {
				c += len(r.Match[k])
			}
			results += fmt.Sprintf("\t%s\t%d\t%v\n", tagAll, c, r.IsMatchAll)
			results += fmt.Sprintf(
				"\t[%s]\t%d/%d\t%v\n",
				v.Origin,
				len(r.Match),
				len(o.Data),
				len(r.Match) == len(o.Data),
			)
		}
	}
	fmt.Println("----[results]----")
	if ins.DiffDB == nil {
		results += fmt.Sprintf("\t%s\t%v\n", tagAll, isMatchAll)
		fmt.Printf("match\n%s", results)
	} else {
		fmt.Println(results)
	}

}

func (ins *Instance) getCountOrigin(i int) int {
	if ins == nil {
		return 0
	}
	d := ins.Data[i]
	fmt.Println("count", d.Origin)

	q := fmt.Sprintf("SELECT COUNT(*) FROM %s %s WHERE %s", d.Origin, d.Omit.Origin, d.Where.Origin)
	fmt.Printf("\t%s\t%s\n", tagSQL, q)

	var c int
	if err := ins.DB.QueryRow(q).Scan(&c); err != nil {
		panic(err)
	}
	fmt.Printf("\t%s\t%d\n", tagCnt, c)
	return c
}

func (ins *Instance) getCountDiff(i int) int {
	if ins == nil {
		return 0
	}
	d := ins.Data[i]
	fmt.Println("count", d.Diff)

	q := fmt.Sprintf("SELECT COUNT(*) FROM %s %s WHERE %s", d.Diff, d.Omit.Diff, d.Where.Diff)
	fmt.Printf("\t%s\t%s\n", tagSQL, q)

	var c int
	if ins.DiffDB != nil {
		if err := ins.DiffDB.QueryRow(q).Scan(&c); err != nil {
			panic(err)
		}
	} else {
		if err := ins.DB.QueryRow(q).Scan(&c); err != nil {
			panic(err)
		}
	}
	fmt.Printf("\t%s\t%d\n", tagCnt, c)
	return c
}

func (ins *Instance) getInnerJoin(i int) [][]sql.NullString {
	if ins == nil {
		return nil
	}
	d := ins.Data[i]
	fmt.Println("inner join")

	q := getInnerJoinQuery(ins, i)
	fmt.Printf("\t%s\t%s\n", tagSQL, q)

	rows, err := ins.DB.Query(q)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	v := [][]sql.NullString{}
	for rows.Next() {
		l := len(d.Columns) * 2
		r := make([]sql.NullString, l)
		p := make([]interface{}, l)
		for i := range r {
			p[i] = &r[i]
		}
		err = rows.Scan(p...)
		if err != nil {
			panic(err)
		}
		v = append(v, r)
	}
	fmt.Printf("\t%s\t%d\n", tagCnt, len(v))
	return v
}

func (ins *Instance) getInnerJoinWithMatch(i int) [][]sql.NullString {
	if ins == nil {
		return nil
	}
	data := ins.Data[i]
	fmt.Println("inner join with match")

	q := getInnerJoinQuery(ins, i)
	for _, v := range data.Columns {
		if v.DisableMatch {
			continue
		}

		if 0 < len(data.Where.Origin) {
			q += " AND "
		}
		s := "%s.%s"
		o := fmt.Sprintf(s, data.Omit.Origin, v.Origin)
		d := fmt.Sprintf(s, data.Omit.Diff, v.Diff)
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
		l := len(data.Columns) * 2
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

func (ins *Instance) findOrigin(i int) QueryResult {
	if ins == nil {
		return QueryResult{}
	}
	d := ins.Data[i]
	fmt.Println("find", d.Origin)

	cs := make([]string, len(d.Columns))
	for i, v := range d.Columns {
		if v == nil {
			continue
		}

		var s string
		if v.Distinct {
			s = "DISTINCT "
		}
		if v.IsRaw {
			s += "%s"
			cs[i] = fmt.Sprintf(s, v.Origin)
		} else {
			s += "%s.%s"
			cs[i] = fmt.Sprintf(s, d.Omit.Origin, v.Origin)
		}
	}

	q := fmt.Sprintf(
		"SELECT %s FROM %s %s WHERE %s%s",
		strings.Join(cs, ", "),
		d.Origin,
		d.Omit.Origin,
		d.Where.Origin,
		d.GetGroupByOrigin(),
	)
	fmt.Printf("\t%s\t%s\n", tagSQL, q)

	rows, err := ins.DB.Query(q)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	data := [][]sql.NullString{}
	for rows.Next() {
		l := len(cs)
		r := make([]sql.NullString, l)
		p := make([]interface{}, l)
		for i := range r {
			p[i] = &r[i]
		}
		err = rows.Scan(p...)
		if err != nil {
			panic(err)
		}
		data = append(data, r)
	}
	fmt.Printf("\t%s\t%d\n", tagCnt, len(data))

	return QueryResult{
		Header: cs,
		Data:   data,
	}
}

func (ins *Instance) findDiff(i int) QueryResult {
	if ins == nil {
		return QueryResult{}
	}
	d := ins.Data[i]
	fmt.Println("find", d.Diff)

	cs := make([]string, len(d.Columns))
	for i, v := range d.Columns {
		if v == nil {
			continue
		}
		var s string
		if v.Distinct {
			s = "DISTINCT "
		}
		if v.IsRaw {
			s += "%s"
			cs[i] = fmt.Sprintf(s, v.Diff)
		} else {
			s += "%s.%s"
			cs[i] = fmt.Sprintf(s, d.Omit.Diff, v.Diff)
		}
	}

	q := fmt.Sprintf(
		"SELECT %s FROM %s %s WHERE %s%s",
		strings.Join(cs, ", "),
		d.Diff,
		d.Omit.Diff,
		d.Where.Diff,
		d.GetGroupByDiff(),
	)
	fmt.Printf("\t%s\t%s\n", tagSQL, q)

	var (
		rows *sql.Rows
		err  error
	)
	if ins.DiffDB == nil {
		rows, err = ins.DB.Query(q)

	} else {
		rows, err = ins.DiffDB.Query(q)
	}
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	data := [][]sql.NullString{}
	for rows.Next() {
		l := len(cs)
		r := make([]sql.NullString, l)
		p := make([]interface{}, l)
		for i := range r {
			p[i] = &r[i]
		}
		err = rows.Scan(p...)
		if err != nil {
			panic(err)
		}
		data = append(data, r)
	}
	fmt.Printf("\t%s\t%d\n", tagCnt, len(data))

	return QueryResult{
		Header: cs,
		Data:   data,
	}
}

func (ins *Instance) checkResults(origin, diff *QueryResult, index int) CheckResult {
	if ins == nil {
		return CheckResult{}
	}
	d := ins.Data[index]

	r := CheckResult{
		Match:      make(map[int][]int),
		IsMatchAll: true,
	}

	basicCnt := len(d.Columns)
	for oy := range origin.Data {
		isMatch := false
		for dy := range diff.Data {
			cnt := 0
			for x := range d.Columns {
				if c := d.Columns[x]; c.DisableMatch {
					cnt++
					continue
				}
				/* 前提: origin と diff は column と同じ数 & 順 */
				o, d := origin.Data[oy][x], diff.Data[dy][x]
				if (o.Valid || d.Valid) && o.String != d.String {
					break
				}
				cnt++
			}
			if cnt == basicCnt {
				isMatch = true
				r.Match[oy] = append(r.Match[oy], dy)
			}
		}
		if r.IsMatchAll {
			r.IsMatchAll = isMatch
		}
	}
	return r
}
