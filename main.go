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
	write = flag.String("write", "",
		`values to write, ex.: root/key1=val1 root/key2=val2`)
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
		result, err := env.Read(*read)
		if err != nil {
			log.Fatalf(`Failed to read path %q: "%v".`, *read, err)
		}
		fmt.Print(result)
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
		result, err := env.Dump()
		if err != nil {
			log.Fatalf(`Failed to dump value: "%v".`, err)
		}
		fmt.Print(result)
		return
	}

	log.Fatalf(`Unknown command.`)
}
