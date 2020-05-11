package queries

import (
	"compare-tables/yaml"
	"context"
	"database/sql"
	"errors"
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
	go serve(chw, chInfR)
	return &Results{
		Left: <-chL,
		Right: <-chR,
	}
}

func serve(ch chan []*Info, chInf chan *Info) {
	infos := make([]*Info)
	for {
		select {
		case i, ok := <- chInf:
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
	}
}
