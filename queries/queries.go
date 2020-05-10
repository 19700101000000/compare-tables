package queries

import (
	"compare-tables/yaml"
	"context"
	"database/sql"
	"errors"
	"time"
)

type Query struct {
	DB     *sql.DB
	Tables []*Table
}

type Instance struct {
	Left  Query
	Right Query
}

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
