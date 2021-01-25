package main

import (
	"flag"
	"fmt"
	"log"

	jsonenv "github.com/palchukovsky/json-env/env"
)

var (
	source = flag.String("source", "CONFIG_BASE64",
		`source, Base64 (without paddings) coded JSON structure`)
	path = flag.String("path", "",
		`path to value, ex.: "root/branch/key"`)
)

func main() {
	flag.Parse()
	if *source == "" {
		log.Fatalf(`Source is not provided.`)
	}
	if *path == "" {
		log.Fatalf(`Path is not provided.`)
	}

	env, err := jsonenv.NewEnv(*source)
	if err != nil {
		log.Fatalf(`Failed to read source: "%v".`, err)
	}

	pathVal, err := env.Read(*path)
	if err != nil {
		log.Fatalf(`Failed to read path %q: "%v".`, *path, err)
	}

	fmt.Print(pathVal)
}
