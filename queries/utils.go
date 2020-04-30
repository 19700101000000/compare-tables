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

		if s := strings.Split(v.Table, sep); 1 < len(s) {
			t.Origin = s[0]
			t.Diff = s[1]
		} else {
			t.Origin = v.Table
			t.Diff = v.Table
		}

		for i, v := range v.Columns {
			c := &Column{
				DisableMatch: v.DisableMatch,
			}
			if s := strings.Split(v.Target, sep); 1 < len(s) {
				c.Origin = s[0]
				c.Diff = s[1]
			} else {
				c.Origin = v.Target
				c.Diff = v.Target
			}
			t.Columns[i] = c
		}

		getCondition(&t.JoinOn.Origin, &t.JoinOn.Diff, v.JoinOn, sep)
		getCondition(&t.Where.Origin, &t.Where.Diff, v.Where, sep)

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

func getCondition(origin, diff *string, c []*my.Condition, sep string) {
	for i, v := range c {
		isAnd := 0 < len(v.And)
		if 0 < i {
			var p string
			if isAnd {
				p = " AND "
			} else {
				p = " OR "
			}

			*origin += p
			*diff += p
		}

		var o, d string
		if isAnd {
			if s := strings.Split(v.And, sep); 1 < len(s) {
				o, d = s[0], s[1]
			} else {
				o, d = v.And, v.And
			}
		} else {
			if s := strings.Split(v.Or, sep); 1 < len(s) {
				o, d = s[0], s[1]
			} else {
				o, d = v.Or, v.Or
			}
		}
		*origin += o
		*diff += d
	}
}
