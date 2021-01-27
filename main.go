package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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
	encodeFile = flag.String("encode_file", "",
		`encode JSON-file to Base64, ex: "/config.json`)
)

func main() {
	flag.Parse()
	if *source == "" && *encodeFile == "" {
		log.Fatalf(`Source or export-file is not provided.`)
	}
	if *read != "" && *write != "" {
		log.Fatalf(`Reading and writing cannot be executed by one start.`)
	}

	var env jsonenv.Env
	if *source != "" {
		var err error
		env, err = jsonenv.NewEnv(*source)
		if err != nil {
			log.Fatalf(`Failed to read source: "%v".`, err)
		}
	}

	if *read != "" {
		if *export {
			log.Fatalf(`Cannot export with reading.`)
		}
		if *source == "" {
			log.Fatalf(`Source is not set for reading.`)
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
		if *source == "" {
			log.Fatalf(`Source is not set for writing.`)
		}
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
		if *source == "" {
			log.Fatalf(`Source is not send for exporting.`)
		}
		result, err := env.Export()
		if err != nil {
			log.Fatalf(`Failed to export source: "%v".`, err)
		}
		fmt.Print(result)
		return
	}

	if *encodeFile != "" {
		fmt.Print(encode(*encodeFile))
		return
	}

	log.Fatalf(`Unknown command.`)
}

func encode(jsonSourceFile string) string {
	file, err := ioutil.ReadFile(jsonSourceFile)
	if err != nil {
		log.Fatalf(`Failed to read JSON-file %q: "%v".`, jsonSourceFile, err)
	}
	var obj json.RawMessage
	if err := json.Unmarshal([]byte(file), &obj); err != nil {
		log.Fatalf(`Failed to parse JSON-file %q: "%v".`, jsonSourceFile, err)
	}
	cleanSource, err := json.Marshal(obj)
	if err != nil {
		log.Fatalf(`Failed to export JSON: "%v".`, err)
	}
	return base64.RawStdEncoding.EncodeToString([]byte(cleanSource))
}
