package yaml

import (
	"compare-tables/types"
)

// Env type
type Env struct {
	Driver types.Driver `yaml:"driver"`
	Host   string       `yaml:"host"`
	Port   string       `yaml:"port"`
	User   string       `yaml:"username"`
	Pass   string       `yaml:"password"`
	DB     string       `yaml:"database"`
}

// Condition type
type Condition struct {
	And *string `yaml:"and"`
	Or  *string `yaml:"or"`
}

// Join type
type Join struct {
	Name string       `yaml:"name"`
	On   []*Condition `yaml:"on"`
}

// Table type
type Table struct {
	Name    string       `yaml:"name"`
	Joins   []*Join      `yaml:"joins"`
	Where   []*Condition `yaml:"where"`
	GroupBy []string     `yaml:"group_by"`
	OrderBy []string     `yaml:"order_by"`
}

// Target type
type Target struct {
	Left  *string `yaml:"left"`
	Right *string `yaml:"right"`
}

// Column type
type Column struct {
	Name       string  `yaml:"name"`
	Join       *Target `yaml:"join"`
	IsDistinct *string `yaml:"distinct"`
	IsRaw      *string `yaml:"is_raw"`
}

// Compare type
type Compare struct {
	Columns []*Column `yaml:"columns"`
	Left    Table     `yaml:"left"`
	Right   Table     `yaml:"right"`
}

// File type
type File struct {
	Env      Env        `yaml:"env"`
	Compares []*Compare `yaml:"compares"`
}
