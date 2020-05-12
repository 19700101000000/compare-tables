package queries

import (
	"compare-tables/yaml"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

// NewInstance is yaml.File to Instance
func NewInstance(f *yaml.File) (*Instance, error) {
	if f == nil {
		return nil, errors.New("file is nil")
	}

	ins := new(Instance)
	ins.Left.DB, ins.Right.DB = getSqlDB(&f.Env)
	ins.Left.Tables, ins.Right.Tables = getTables(f.Compares)
	ins.Labels = getLabels(f.Compares)
	return ins, nil
}

// Close close should use end.
func (ins *Instance) Close() {
	if ins == nil {
		return
	}
	ins.Left.DB.Close()
	ins.Right.DB.Close()
}

// Ping connection check db.
func (ins *Instance) Ping() error {
	if ins == nil {
		return errors.New("instance is nil!")
	}

	ch := make(chan error)
	f := func(ctx context.Context, db *sql.DB, ch chan error) {
		ch <- db.PingContext(ctx)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	go f(ctx, ins.Left.DB, ch)
	go f(ctx, ins.Right.DB, ch)

	err1, err2 := <-ch, <-ch
	if err1 != nil {
		return err1
	} else if err2 != nil {
		return err2
	}
	return nil
}

// Exec instance do
func (ins *Instance) Exec() *Results {
	if ins == nil {
		return nil
	}
	chInfL, chInfR := make(chan *Info), make(chan *Info)
	chL, chR := make(chan []*Info), make(chan []*Info)
	go exec(chInfL, &ins.Left)
	go serve(chL, chInfL)
	go exec(chInfR, &ins.Right)
	go serve(chR, chInfR)
	return &Results{
		Labels: ins.Labels,
		Left:   <-chL,
		Right:  <-chR,
	}
}

// Compare compare tables
func (rs *Results) Compare() {
	fmt.Println("result start ---->>")
	defer fmt.Println("<<---- end result")
	if rs == nil {
		return
	}

	var max int
	if l, r := len(rs.Left), len(rs.Right); l < r {
		max = l
	} else {
		max = r
	}
	for i := 0; i < max; i++ {
		fmt.Println("-", i, rs.Labels[i])
		l, r := rs.Left[i], rs.Right[i]
		if l == nil || r == nil {
			continue
		}
		fmt.Printf("\t[SQL]\t%s\n", l.Query)
		fmt.Printf("\t[CNT]\t%d\n", len(l.Data))
		fmt.Printf("\t[SQL]\t%s\n", r.Query)
		fmt.Printf("\t[CNT]\t%d\n", len(r.Data))

		var max int
		if l, r := len(l.Data), len(r.Data); l < r {
			max = l
		} else {
			max = r
		}

		var cnt, cntAll, cntY int
		var failed string
		for y := 0; y < max; y++ {
			l, r := l.Data[y], r.Data[y]
			var max int
			if l, r := len(l), len(r); l < r {
				max = l
			} else {
				max = r
			}
			var cntX int
			for x := 0; x < max; x++ {
				l, r := l[x], r[x]
				if ok := (!l.Valid && !r.Valid) || l.String == r.String; ok {
					cntX++
				} else {
					var vl, vr string
					if l.Valid {
						vl = l.String
					} else {
						vl = "NULL"
					}
					if r.Valid {
						vr = r.String
					} else {
						vr = "NULL"
					}
					failed += fmt.Sprintf("\t[%d:%d] %s <> %s", y, x, vl, vr)
				}
			}
			if max == cntX {
				cntY++
			}
			cnt += cntX
			cntAll += max
		}
		if max == cntY {
			fmt.Println("\t[RST]\tSUCCESS")
		} else {
			fmt.Printf("\t[RST]\tFAILED%s\n", failed)
		}
		fmt.Println("\t\tmatched", cnt, "/", cntAll)
	}
}

func serve(ch chan []*Info, chInf chan *Info) {
	infos := make([]*Info, 0)
	for {
		select {
		case i, ok := <-chInf:
			if !ok {
				ch <- infos
				return
			}
			infos = append(infos, i)
		}
	}
}

func exec(ch chan *Info, q *Query) {
	defer close(ch)
	if q == nil {
		return
	}

	for _, t := range q.Tables {
		if t == nil || len(t.Columns) == 0 {
			ch <- nil
			continue
		}

		query := "SELECT "
		for i, c := range t.Columns {
			if c == nil {
				continue
			}
			if i > 0 {
				query += ", "
			}
			query += c.Name
		}
		query += fmt.Sprintf(" FROM %s %s", t.FullName, t.Name)

		for _, j := range t.Joins {
			if j == nil {
				continue
			}
			query += fmt.Sprintf(" INNER JOIN %s %s", j.FullName, j.Name)
			if j.On != nil {
				query += fmt.Sprintf(" ON %s", *j.On)
			}
		}

		if t.Where != nil {
			query += fmt.Sprintf(" WHERE %s", *t.Where)
		}

		if t.GroupBy != nil {
			query += fmt.Sprintf(" GROUP BY %s", *t.GroupBy)
		}

		if t.OrderBy != nil {
			query += fmt.Sprintf(" ORDER BY %s", *t.OrderBy)
		}

		rows, err := q.DB.Query(query)
		if err != nil {
			log.Println(err)
			ch <- &Info{
				Query: query,
			}
			continue
		}
		defer rows.Close()
		data := [][]sql.NullString{}
		for rows.Next() {
			l := len(t.Columns)
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

		ch <- &Info{
			Query: query,
			Data:  data,
			Ok:    true,
		}
	}
}
