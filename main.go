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

	data := readFile()

	ins := queries.GetInstance(env)
	defer ins.Close()

	ins.Init(data)
	ins.RunCompare()
	log.Println("done")
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

func readFile() []*my.Table {
	if len(os.Args) < 2 {
		panic(
			fmt.Sprintf("please set filename"),
		)
	}
	f := os.Args[1]
	buf, err := ioutil.ReadFile(f)
	if err != nil {
		panic(
			fmt.Sprintf("cannot readfile %s: %v", filename.Env, err),
		)
	}

	var t []*my.Table
	err = yaml.Unmarshal(buf, &t)
	if err != nil {
		panic(
			fmt.Sprintf("cannot unmarshal %s: %v", f, err),
		)
	}
	log.Printf("open file: %s\n", f)
	return t
}
