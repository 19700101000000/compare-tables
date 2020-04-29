package yaml

// Env env.yml's struct
type Env struct {
	Driver string `yaml:"driver"`
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	User   string `yaml:"username"`
	Pass   string `yaml:"password"`
	DB     string `yaml:"database"`
}

// Condition use sql-condition
type Condition struct {
	And string `yaml:"and"`
	Or  string `yaml:"or"`
}

// Column table's column
type Column struct {
	Target string `yaml:"target"`
}

// Table any.yml's struct
type Table struct {
	Table   string      `yaml:"table"`
	Columns []Column    `yaml:"columns"`
	JoinOn  []Condition `yaml:"join_on"`
	Where   []Condition `yaml:"where"`
}
