package yaml

import (
	"compare-tables/types"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

const (
	driverPsql  = "postgres"
	driverMysql = "mysql"
)

// ReadFile open read file
func ReadFile(fname string) *File {
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(
			fmt.Sprintf("cannot readfile %s\n%v", fname, err),
		)
	}

	f := new(File)

	if err = yaml.Unmarshal(buf, f); err != nil {
		panic(
			fmt.Sprintf("cannot unmarshal %s\n%v", fname, err),
		)
	}

	return f
}

// Output output file content to standart output
func (f *File) Output() {
	fmt.Println("file output start ---->>")
	defer fmt.Println("<<---- file output end")

	s := "\t%s\t%v\n"
	fmt.Printf("%s\n", "env:")
	fmt.Printf(s, "driver:", f.Env.Driver)
	fmt.Printf(s, "host:", f.Env.Host)
	fmt.Printf(s, "port:", f.Env.Port)
	fmt.Printf(s, "username:", f.Env.User)
	fmt.Printf(s, "passrowd:", f.Env.Pass)
	fmt.Printf(s, "database:", f.Env.DB)

	fmt.Println("compares:")
	for _, v := range f.Compares {
		fmt.Println("\tcolumns:")
		for i, v := range v.Columns {
			s := "\t" + s
			fmt.Printf("\t- %d\n", i)
			fmt.Printf(s, "name:", v.Name)
			fmt.Println("\t\tjoin:")
			fmt.Print(outputTarget("\t\t\t", v.Join))
			fmt.Printf(s, "is_distinct:", v.IsDistinct)
			fmt.Printf(s, "is_raw:", v.IsRaw)
		}
		fmt.Print(outputTable("\t", "left", &v.Left))
		fmt.Print(outputTable("\t", "right", &v.Right))
	}
}

// GetDriverNames getting db driver-names.
func (env *Env) GetDriverNames() (left, right string) {
	if env == nil {
		return
	}

	left, right, _ = splitWord(string(env.Driver), ":")
	return
}

// GetSrcNames getting db source-names.
func (env *Env) GetSrcNames() (left, right string, isSame bool) {
	if env == nil {
		return
	}

	l, r := new(Env), new(Env)
	sep, isSame := ":", true

	var is bool
	if ld, rd, is := splitWord(string(env.Driver), sep); !is && isSame {
		isSame = false
		l.Driver, r.Driver = types.Driver(ld), types.Driver(rd)
	} else {
		l.Driver, r.Driver = types.Driver(ld), types.Driver(rd)
	}
	if l.Host, r.Host, is = splitWord(env.Host, sep); !is && isSame {
		isSame = false
	}
	if l.Port, r.Port, is = splitWord(env.Port, sep); !is && isSame {
		isSame = false
	}
	if l.User, r.User, is = splitWord(env.User, sep); !is && isSame {
		isSame = false
	}
	if l.Pass, r.Pass, is = splitWord(env.Pass, sep); !is && isSame {
		isSame = false
	}
	if l.DB, r.DB, is = splitWord(env.DB, sep); !is && isSame {
		isSame = false
	}

	left = l.getSrcName()
	right = r.getSrcName()
	return
}

// GetLeftTableNames getting left table-names
func (cmp *Compare) GetLeftTableNames() (name, fullname string) {
	if cmp == nil {
		return
	}
	name, fullname = getTableNames(cmp.Left.Name)
	return
}

// GetRightTableNames getting right table-names
func (cmp *Compare) GetRightTableNames() (name, fullname string) {
	if cmp == nil {
		return
	}
	name, fullname = getTableNames(cmp.Right.Name)
	return
}

// GetWheres
func (cmp *Compare) GetWheres() (left, right *string) {
	left = cmp.GetLeftWhere()
	right = cmp.GetRightWhere()
	return
}

// GetLeftWhere getting left where
func (cmp *Compare) GetLeftWhere() *string {
	if cmp == nil {
		return nil
	}
	return getCondition(cmp.Left.Where)
}

// GetRightWhere getting right where
func (cmp *Compare) GetRightWhere() *string {
	if cmp == nil {
		return nil
	}
	return getCondition(cmp.Right.Where)
}

// GetGroupBys getting group-bys
func (cmp *Compare) GetGroupBys() (left, right *string) {
	if cmp == nil {
		return
	}
	left = joinStrings(cmp.Left.GroupBy)
	right = joinStrings(cmp.Right.GroupBy)
	return
}

// GetOrderBys getting order-bys
func (cmp *Compare) GetOrderBys() (left, right *string) {
	if cmp == nil {
		return
	}
	left = joinStrings(cmp.Left.OrderBy)
	right = joinStrings(cmp.Right.OrderBy)
	return
}

// GetColumnName getting column-name
func (col *Column) GetColumnName(tablename string, isLeft bool) string {
	if col == nil {
		return ""
	}

	var name string
	if l, r, _ := splitWord(col.Name, ":"); isLeft {
		name = l
	} else {
		name = r
	}

	if GetBool(col.IsRaw, isLeft) {
		return name
	}

	if col.Join != nil {
		if col.Join.Left != nil && isLeft {
			tablename = *col.Join.Left
		} else if col.Join.Right != nil {
			tablename = *col.Join.Right
		}
	}

	name = fmt.Sprintf("%s.%s", tablename, name)
	if GetBool(col.IsDistinct, isLeft) {
		name = "DISTINCT " + name
	}
	return name
}

// GetBool is getting string to bool
func GetBool(p *string, isLeft bool) bool {
	if p == nil {
		return false
	}

	l, r, _ := splitWord(*p, ":")

	var s string
	if isLeft {
		s = l
	} else {
		s = r
	}

	if s == "true" {
		return true
	} else if s == "false" {
		return false
	}
	panic(fmt.Sprintf("cannot cast to bool: %s", s))
}

// GetJoinNames getting join-names
func (j *Join) GetJoinNames() (name, fullname string) {
	if j == nil {
		return
	}
	fullname, name, _ = splitWord(j.Name, " ")
	return
}

// GetJoinOn getting join-on
func (j *Join) GetJoinOn() *string {
	if j == nil {
		return nil
	}
	return getCondition(j.On)
}

func joinStrings(s []string) *string {
	if len(s) == 0 {
		return nil
	}
	groupby := strings.Join(s, ", ")
	return &groupby
}

func (env *Env) getSrcName() string {
	if env == nil {
		return ""
	}

	if env.Driver == driverPsql {
		return fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			env.Host,
			env.Port,
			env.User,
			env.Pass,
			env.DB,
		)
	}

	if env.Driver == driverMysql {
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
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
			driverMysql,
			driverPsql,
		),
	)
}

func getCondition(cs []*Condition) *string {
	if len(cs) == 0 {
		return nil
	}

	var s string
	for i, v := range cs {
		if v == nil || (v.And == nil && v.Or == nil) {
			continue
		}
		if i > 0 {
			if v.And != nil {
				s += " AND "
			} else {
				s += " OR "
			}
		}
		if v.And != nil {
			s += *v.And
		} else {
			s += *v.Or
		}
	}
	return &s
}

func getTableNames(str string) (name, fullname string) {
	sep := " "
	fullname, name, _ = splitWord(str, sep)
	return
}

func splitWord(w, sep string) (left, right string, isSame bool) {
	if s := strings.Split(w, sep); 1 < len(s) {
		left, right, isSame = s[0], s[1], false
	} else {
		left, right, isSame = w, w, true
	}
	return
}

func outputTable(indent, title string, table *Table) string {
	if table == nil {
		return ""
	}

	var s string
	s += fmt.Sprintf(indent+"%s:\n", title)
	indent += "\t"

	s += fmt.Sprintf(indent+"name:\t%v\n", table.Name)
	if len(table.Joins) > 0 {
		s += fmt.Sprintln(indent + "joins:")
		for i, v := range table.Joins {
			if v == nil {
				continue
			}
			s += fmt.Sprintf(indent+"- %d\n", i)
			indent := indent + "\t"

			s += fmt.Sprintf(indent+"name:\t%v\n", v.Name)
			s += fmt.Sprintln(indent + "on:")
			for i, v := range v.On {
				if v == nil {
					continue
				}
				s += fmt.Sprintf(indent+"- %d\n", i)
				indent := indent + "\t"

				s += outputCondition(indent, v)
			}
		}
	}

	if len(table.Where) > 0 {
		s += fmt.Sprintln(indent + "where:")
		for i, v := range table.Where {
			if v == nil {
				continue
			}
			s += fmt.Sprintf(indent+"- %d\n", i)
			indent := indent + "\t"

			s += outputCondition(indent, v)
		}
	}

	if len(table.GroupBy) > 0 {
		s += fmt.Sprintln(indent + "group_by:")
		for i, v := range table.GroupBy {
			s += fmt.Sprintf(indent+"- %d\t%s\n", i, v)
		}
	}

	return s
}

func outputCondition(indent string, c *Condition) string {
	if c == nil {
		return ""
	}

	var s string
	if c.And != nil {
		s += fmt.Sprintf(indent+"and:\t%v\n", *c.And)
	} else if c.Or != nil {
		s += fmt.Sprintf(indent+"or:\t%v\n", *c.Or)
	}
	return s
}

func outputTarget(indent string, t *Target) string {
	if t == nil {
		return ""
	}

	var s string
	if t.Left != nil {
		s += fmt.Sprintf(indent+"left:\t%v\n", *t.Left)
	}

	if t.Right != nil {
		s += fmt.Sprintf(indent+"right:\t%v\n", *t.Right)
	}
	return s
}
