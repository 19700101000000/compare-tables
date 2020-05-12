package queries

import (
	"compare-tables/yaml"
	"database/sql"
	"fmt"

	// sql drivers
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func getSqlDB(env *yaml.Env) (left, right *sql.DB) {
	if env == nil {
		return
	}

	srcL, srcR, isSame := env.GetSrcNames()
	driverL, driverR := env.GetDriverNames()

	left, err := sql.Open(driverL, srcL)
	if err != nil {
		panic(
			fmt.Sprintf("coannot open %s\n%v", driverL, err),
		)
	}
	if isSame {
		right = left
	} else {
		right, err = sql.Open(driverR, srcR)
		if err != nil {
			panic(
				fmt.Sprintf("coannot open %s\n%v", driverR, err),
			)
		}
	}
	return
}

func getTables(cmps []*yaml.Compare) (left, right []*Table) {
	if cmps == nil {
		return
	}
	l := len(cmps)
	left = make([]*Table, l)
	right = make([]*Table, l)
	for i := range cmps {
		cmp := cmps[i]
		if cmp == nil {
			continue
		}

		var l, r Table
		l.Name, l.FullName = cmp.GetLeftTableNames()
		r.Name, r.FullName = cmp.GetRightTableNames()
		l.Where, r.Where = cmp.GetWheres()
		l.GroupBy, r.GroupBy = cmp.GetGroupBys()
		l.OrderBy, r.OrderBy = cmp.GetOrderBys()
		l.setColumns(cmp.Columns, true)
		r.setColumns(cmp.Columns, false)
		l.setJoins(cmp.Left.Joins)
		r.setJoins(cmp.Right.Joins)
		left[i], right[i] = &l, &r
	}
	return
}

func getLabels(cmps []*yaml.Compare) []string {
	if cmps == nil {
		return nil
	}
	ls := make([]string, len(cmps))
	for i := range cmps {
		ls[i] = cmps[i].Label
	}
	return ls
}

func (t *Table) setJoins(js []*yaml.Join) {
	if t == nil {
		return
	}

	t.Joins = make([]*Join, len(js))
	for i, j := range js {
		if j == nil {
			continue
		}
		n, fn := j.GetJoinNames()
		t.Joins[i] = &Join{
			Name:     n,
			FullName: fn,
			On:       j.GetJoinOn(),
		}
	}
}

func (t *Table) setColumns(cols []*yaml.Column, isLeft bool) {
	if t == nil {
		return
	}

	t.Columns = make([]*Column, len(cols))
	for i, col := range cols {
		t.Columns[i] = &Column{
			Name: col.GetColumnName(t.Name, isLeft),
		}
	}
}
