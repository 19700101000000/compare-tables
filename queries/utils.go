package queries

import (
	"compare-tables/types"
	my "compare-tables/yaml"
	"fmt"
	"strings"
)

func getEnvs(env my.Env) (origin, diff my.Env) {
	sep := ":"

	od, dd := splitWord(string(env.Driver), sep)
	origin.Driver = types.Driver(od)
	diff.Driver = types.Driver(dd)
	origin.Host, diff.Host = splitWord(env.Host, sep)
	origin.Port, diff.Port = splitWord(env.Port, sep)
	origin.User, diff.User = splitWord(env.User, sep)
	origin.Pass, diff.Pass = splitWord(env.Pass, sep)
	origin.DB, diff.DB = splitWord(env.DB, sep)

	return
}

func getSrcName(env my.Env) string {
	if env.Driver == driverPsql {
		return fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			env.Host,
			env.Port,
			env.User,
			env.Pass,
			env.DB,
		)
	}

	if env.Driver == driverMysql {
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			env.User,
			env.Pass,
			env.Host,
			env.Port,
			env.DB,
		)
	}

	panic(
		fmt.Sprintf(
			"cannot found driver: %s\nused by %s | %s",
			env.Driver,
			driverMysql,
			driverPsql,
		),
	)

}

func getData(data []*my.Table) []*Query {
	queries := make([]*Query, len(data))
	sep := ":"
	for i, v := range data {
		q := &Query{
			Columns: make([]*Column, len(v.Columns)),
		}
		o, d := splitWord(v.Table, sep)
		s := " "
		q.Origin, q.Omit.Origin = splitWord(o, s)
		q.Diff, q.Omit.Diff = splitWord(d, s)

		for i, v := range v.Columns {
			c := &Column{
				DisableMatch: v.DisableMatch,
			}
			c.Origin, c.Diff = splitWord(v.Target, sep)
			q.Columns[i] = c
		}

		q.JoinOn.Origin, q.JoinOn.Diff = getCondition(v.JoinOn, sep)
		q.Where.Origin, q.Where.Diff = getCondition(v.Where, sep)

		queries[i] = q
	}

	return queries
}

func getInnerJoinQuery(ins *Instance, i int) string {
	d := ins.Data[i]
	cols := make([]string, len(d.Columns)*2)
	for i := range d.Columns {
		c := d.Columns[i]
		s := "%s.%s"

		i *= 2
		cols[i] = fmt.Sprintf(s, d.Omit.Origin, c.Origin)
		cols[i+1] = fmt.Sprintf(s, d.Omit.Diff, c.Diff)
	}

	q := fmt.Sprintf(
		"SELECT %s FROM %s %s INNER JOIN %s %s ON %s WHERE %s",
		strings.Join(cols, ", "),
		d.Origin,
		d.Omit.Origin,
		d.Diff,
		d.Omit.Diff,
		d.JoinOn.Origin,
		d.Where.Origin,
	)
	return q
}

func getCondition(c []*my.Condition, sep string) (o, d string) {
	for i, v := range c {
		isAnd := 0 < len(v.And)
		if 0 < i {
			var p string
			if isAnd {
				p = " AND "
			} else {
				p = " OR "
			}
			o, d = o+p, d+p
		}

		var s string
		if isAnd {
			s = v.And
		} else {
			s = v.Or
		}
		t, u := splitWord(s, sep)
		o, d = o+t, d+u
	}
	return
}

func splitWord(w, sep string) (o, d string) {
	if s := strings.Split(w, sep); 1 < len(s) {
		o, d = s[0], s[1]
	} else {
		o, d = w, w
	}
	return
}
