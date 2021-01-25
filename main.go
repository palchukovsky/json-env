package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	jsonenv "github.com/palchukovsky/json-env/env"
)

var (
	source = flag.String("source", "",
		`source, Base64 (without paddings) coded JSON structure`)
	read = flag.String("read", "",
		`path to value to read, ex.: root/branch/key`)
	defaultValue = flag.String("default", "",
		`default result for value reading, ex.: "default value"`)
	write = flag.String("write", "",
		`values to write, ex.: root/key1=val1 root/key2=val2`)
	export = flag.Bool("export", false,
		`export source, could be used with writing`)
)

func main() {
	flag.Parse()
	if *source == "" {
		log.Fatalf(`Source is not provided.`)
	}
	if *read != "" && *write != "" {
		log.Fatalf(`Reading and writing cannot be executed by one start.`)
	}

	env, err := jsonenv.NewEnv(*source)
	if err != nil {
		log.Fatalf(`Failed to read source: "%v".`, err)
	}

	if *read != "" {
		if *export {
			log.Fatalf(`Cannot export with reading.`)
		}
		result, err := env.Read(*read)
		if err != nil {
			log.Fatalf(`Failed to read path %q: "%v".`, *read, err)
		}
		if result == nil {
			if *defaultValue == "" {
				log.Fatalf(`Path %q is not existent.`, *read)
			}
			result = defaultValue
		}
		fmt.Print(*result)
		return
	}

	if *write != "" {
		args := flag.Args()
		args = append(args, *write)
		for _, a := range args {
			v := strings.Split(strings.TrimSpace(a), "=")
			if len(v) != 2 {
				log.Fatalf(`Invalid format of key-value pair %q.`, a)
			}
			err := env.Set(strings.TrimSpace(v[0]), strings.TrimSpace(v[1]))
			if err != nil {
				log.Fatalf(`Failed to set %q = %q: "%v".`, v[0], v[1], err)
			}
		}
		if !*export {
			result, err := env.Dump()
			if err != nil {
				log.Fatalf(`Failed to dump source: "%v".`, err)
			}
			fmt.Print(result)
			return
		}
	}

	if *export {
		result, err := env.Export()
		if err != nil {
			log.Fatalf(`Failed to export source: "%v".`, err)
		}
		fmt.Print(result)
		return
	}

	log.Fatalf(`Unknown command.`)
}
