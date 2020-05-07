package yaml

import (
	"compare-tables/types"
)

// Env env.yml's struct
type Env struct {
	Driver  types.Driver `yaml:"driver"`
	Host    string       `yaml:"host"`
	Port    string       `yaml:"port"`
	User    string       `yaml:"username"`
	Pass    string       `yaml:"password"`
	DB      string       `yaml:"database"`
	NotSame bool         `yaml:"not_same"`
}

// Condition use sql-condition
type Condition struct {
	And string `yaml:"and"`
	Or  string `yaml:"or"`
}

// Column table's column
type Column struct {
	Target       string `yaml:"target"`
	DisableMatch bool   `yaml:"disable_match"`
	Distinct     bool   `yaml:"distinct"`
}

// Table any.yml's struct
type Table struct {
	Table   string       `yaml:"table"`
	Columns []*Column    `yaml:"columns"`
	JoinOn  []*Condition `yaml:"join_on"`
	Where   []*Condition `yaml:"where"`
	GroupBy []string     `yaml:"group_by"`
}
