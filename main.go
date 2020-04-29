package main

import (
	"compare-tables/filename"
	"compare-tables/queries"
	my "compare-tables/yaml"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalln(r)
		}
	}()

	env := readEnv()

	readFile()

	queries.GetSrcName(env)
}

func readEnv() my.Env {
	buf, err := ioutil.ReadFile(filename.Env)
	if err != nil {
		log.Fatalf("cannot readfile %s: %v", filename.Env, err)
	}

	var env my.Env
	err = yaml.Unmarshal(buf, &env)
	if err != nil {
		panic(
			fmt.Sprintf("cannot unmarshal %s: %v", filename.Env, err),
		)
	}
	log.Printf("open file: %s\n", filename.Env)

	return env
}

func readFile() []my.Table {
	if len(os.Args) < 2 {
		panic(
			fmt.Sprintf("please set filename"),
		)
	}
	tgt := os.Args[1]
	buf, err := ioutil.ReadFile(tgt)
	if err != nil {
		panic(
			fmt.Sprintf("cannot readfile %s: %v", filename.Env, err),
		)
	}

	var tables []my.Table
	err = yaml.Unmarshal(buf, &tables)
	if err != nil {
		panic(
			fmt.Sprintf("cannot unmarshal %s: %v", tgt, err),
		)
	}
	log.Printf("open file: %s\n", tgt)

	return tables
}
