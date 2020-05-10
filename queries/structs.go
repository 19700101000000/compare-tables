package queries

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
}
