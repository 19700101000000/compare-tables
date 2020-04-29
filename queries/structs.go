package queries

// Target datatype
type Target struct {
	Origin string
	Diff   string
}

// Table target table
type Table struct {
	Target
	Columns []*Target
	JoinOn  Target
	Where   Target
}
