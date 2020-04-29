package queries

import (
	my "compare-tables/yaml"
	"fmt"
)

// sql drivers
const (
	DriverPsql  = "posgres"
	DriverMysql = "mysql"
)

// GetSrcName get source-name when use sql.Open
func GetSrcName(env my.Env) string {
	if env.Driver == DriverPsql {
		return fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s",
			env.Host,
			env.Port,
			env.User,
			env.Pass,
			env.DB,
		)
	}

	if env.Driver == DriverMysql {
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s",
			env.User,
			env.Pass,
			env.Host,
			env.Port,
			env.DB,
		)
	}

	panic(
		fmt.Sprintf(
			"cannot found driver: %s\nused by %s | %s",
			env.Driver,
			DriverMysql,
			DriverPsql,
		),
	)
}
