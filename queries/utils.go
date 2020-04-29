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
		table := &Table{
			Columns: make([]*Target, len(v.Columns)),
		}

		if t := strings.Split(v.Table, sep); 1 < len(t) {
			table.Origin = t[0]
			table.Diff = t[1]
		} else {
			table.Origin = v.Table
			table.Diff = v.Table
		}

		for i, v := range v.Columns {
			column := &Target{}
			if c := strings.Split(v.Target, sep); 1 < len(c) {
				column.Origin = c[0]
				column.Diff = c[1]
			} else {
				column.Origin = v.Target
				column.Diff = v.Target
			}
			table.Columns[i] = column
		}

		for i, v := range v.JoinOn {
			if 0 < i {
				if 0 < len(v.And) {
					pref := " AND "

					if a := strings.Split(v.And, sep); 1 < len(a) {
						table.JoinOn.Origin += pref + a[0]
						table.JoinOn.Diff += pref + a[1]
					} else {
						table.JoinOn.Origin += pref + v.And
						table.JoinOn.Diff += pref + v.And
					}
				} else {
					pref := " OR "

					if o := strings.Split(v.Or, sep); 1 < len(o) {
						table.JoinOn.Origin += pref + o[0]
						table.JoinOn.Diff += pref + o[1]
					} else {
						table.JoinOn.Origin += pref + v.Or
						table.JoinOn.Diff += pref + v.Or
					}
				}
			} else {
				if 0 < len(v.And) {

					if a := strings.Split(v.And, sep); 1 < len(a) {
						table.JoinOn.Origin = a[0]
						table.JoinOn.Diff = a[1]
					} else {
						table.JoinOn.Origin = v.And
						table.JoinOn.Diff = v.And
					}
				} else {

					if o := strings.Split(v.Or, sep); 1 < len(o) {
						table.JoinOn.Origin = o[0]
						table.JoinOn.Diff = o[1]
					} else {
						table.JoinOn.Origin = v.Or
						table.JoinOn.Diff = v.Or
					}
				}
			}
		}

		for i, v := range v.Where {
			if 0 < i {
				if 0 < len(v.And) {
					pref := " AND "

					if a := strings.Split(v.And, sep); 1 < len(a) {
						table.Where.Origin += pref + a[0]
						table.Where.Diff += pref + a[1]
					} else {
						table.Where.Origin += pref + v.And
						table.Where.Diff += pref + v.And
					}
				} else {
					pref := " OR "

					if o := strings.Split(v.Or, sep); 1 < len(o) {
						table.Where.Origin += pref + o[0]
						table.Where.Diff += pref + o[1]
					} else {
						table.Where.Origin += pref + v.Or
						table.Where.Diff += pref + v.Or
					}
				}
			} else {
				if 0 < len(v.And) {

					if a := strings.Split(v.And, sep); 1 < len(a) {
						table.Where.Origin = a[0]
						table.Where.Diff = a[1]
					} else {
						table.Where.Origin = v.And
						table.Where.Diff = v.And
					}
				} else {

					if o := strings.Split(v.Or, sep); 1 < len(o) {
						table.Where.Origin = o[0]
						table.Where.Diff = o[1]
					} else {
						table.Where.Origin = v.Or
						table.Where.Diff = v.Or
					}
				}
			}
		}

		tables[i] = table
	}

	return tables
}
