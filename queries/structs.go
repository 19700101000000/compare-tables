package queries

import (
	"fmt"
)

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
	IsRaw        bool
}

// Query query info
type Query struct {
	Table
	Columns []*Column
	JoinOn  Target
	Where   Target
	GroupBy []*Target
}

// GetGroupByOrigin get string group by
func (q *Query) GetGroupByOrigin() string {
	if q == nil {
		return ""
	}
	var groupby string
	if l := len(q.GroupBy); l > 0 {
		groupby = " GROUP BY "
		for i, g := range q.GroupBy {
			if g == nil {
				continue
			}

			if i > 0 {
				groupby += ", "
			}
			groupby += fmt.Sprintf("%s.%s", q.Omit.Origin, g.Origin)
		}
	}
	return groupby
}

// GetGroupByDiff get string group by
func (q *Query) GetGroupByDiff() string {
	if q == nil {
		return ""
	}
	var groupby string
	if l := len(q.GroupBy); l > 0 {
		groupby = " GROUP BY "
		for i, g := range q.GroupBy {
			if g == nil {
				continue
			}

			if i > 0 {
				groupby += ", "
			}
			groupby += fmt.Sprintf("%s.%s", q.Omit.Diff, g.Diff)
		}
	}
	return groupby
}
