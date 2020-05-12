package queries

import (
	"database/sql"
)

// Column type
type Column struct {
	Name string
}

// Join type
type Join struct {
	Name     string
	FullName string
	On       *string
}

// Table type
type Table struct {
	Name     string
	FullName string
	Columns  []*Column
	Joins    []*Join
	Where    *string
	GroupBy  *string
	OrderBy  *string
}

// Query query info
type Query struct {
	DB     *sql.DB
	Tables []*Table
}

// Instance queries
type Instance struct {
	Labels []string
	Left   Query
	Right  Query
}

// Info query result info
type Info struct {
	Query string
	Data  [][]sql.NullString
	Ok    bool
}

// Results queries results
type Results struct {
	Labels []string
	Left   []*Info
	Right  []*Info
}
