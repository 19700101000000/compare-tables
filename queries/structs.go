package queries

// Target datatype
type Target struct {
	Origin string
	Diff   string
}

// Table Table info
type Table struct {
	Target
	Omit Target
}

// Column column info
type Column struct {
	Target
	DisableMatch bool
	Distinct     bool
}

// Query query info
type Query struct {
	Table
	Columns []*Column
	JoinOn  Target
	Where   Target
}
