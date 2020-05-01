package queries

import (
	my "compare-tables/yaml"
	"fmt"
	"strings"
)

func getSrcName(env my.Env) string {
	if env.Driver == driverPsql {
		return fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			env.Host,
			env.Port,
			env.User,
			env.Pass,
			env.DB,
		)
	}

	if env.Driver == driverMysql {
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s",
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

func getData(data []*my.Table) []*Table {
	tables := make([]*Table, len(data))
	for i, v := range data {
		sep := ":"
		t := &Table{
			Columns: make([]*Column, len(v.Columns)),
		}
		t.Origin, t.Diff = splitWord(v.Table, sep)

		for i, v := range v.Columns {
			c := &Column{
				DisableMatch: v.DisableMatch,
			}
			c.Origin, c.Diff = splitWord(v.Target, sep)
			t.Columns[i] = c
		}

		t.JoinOn.Origin, t.JoinOn.Diff = getCondition(v.JoinOn, sep)
		t.Where.Origin, t.Where.Diff = getCondition(v.Where, sep)

		tables[i] = t
	}

	return tables
}

func getInnerJoinQuery(ins *Instance, i int) string {
	t := ins.Data[i]
	cols := make([]string, len(t.Columns)*2)
	for i := range t.Columns {
		c := t.Columns[i]
		s := "%s.%s"

		i *= 2
		cols[i] = fmt.Sprintf(s, t.Origin, c.Origin)
		cols[i+1] = fmt.Sprintf(s, t.Diff, c.Diff)
	}

	q := fmt.Sprintf(
		"SELECT %s FROM %s INNER JOIN %s ON %s WHERE %s",
		strings.Join(cols, ", "),
		t.Origin,
		t.Diff,
		t.JoinOn.Origin,
		t.Where.Origin,
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
