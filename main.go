package main

import (
	"compare-tables/queries"
	"compare-tables/yaml"
	"fmt"
	"log"
	"os"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalln(r)
		}
	}()

	if len(os.Args) < 2 {
		panic(
			fmt.Sprintf("please set filename"),
		)
	}
	fname := os.Args[1]
	f := yaml.ReadFile(fname)
	log.Println("file open ok:", fname)
	// f.Output()

	ins, err := queries.NewInstance(f)
	if err != nil {
		log.Fatalln(err)
	}
	defer ins.Close()

	err = ins.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("db connection ok.")

	rs := ins.Exec()
	rs.Compare()
}
