package queries

// Target datatype
type Target struct {
	Origin string
	Diff   string
}

// Column column info
type Column struct {
	Target
	DisableMatch bool
}

// Table target table
type Table struct {
	Target
	Columns []*Column
	JoinOn  Target
	Where   Target
}
