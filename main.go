package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/go-yaml/yaml"
)

const (
	envYml = "env.yml"
)

type Info struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	User string `yaml:"username"`
	Pass string `yaml:"password"`
	DB   string `yaml:"database"`
}

type Env struct {
	Psql Info `yaml:"psql"`
}

func main() {
	buf, err := ioutil.ReadFile(envYml)
	if err != nil {
		log.Fatalf("cannot readfile %s: %v", envYml, err)
	}

	var env Env
	err = yaml.Unmarshal(buf, &env)
	if err != nil {
		log.Fatalf("cannot unmarshal %s: %v", envYml, err)
	}

	fmt.Printf("%#v\n", env)

	fmt.Println("hello, world")
}
