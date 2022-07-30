package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	nginxparser "github.com/faceair/nginx-parser"
)

func main() {
	serverName := flag.String("server_name", "", "server_name, only this one will be secured")
	nginxConfig := flag.String("nginx_config", "", "nginx config filepath")
	inplace := flag.Bool("inplace", false, "whether to edit inplace (defaults to stdout emission)")
	flag.Parse()

	cliParseErrors := false
	cliErrorCode := 1

	if *serverName == "" {
		_, err := fmt.Fprintln(flag.CommandLine.Output(), "`server_name`\trequired but found unset")
		if err != nil {
			panic(err)
		}
		cliParseErrors = true
	}
	if *nginxConfig == "" {
		_, err := fmt.Fprintln(flag.CommandLine.Output(), "`nginx_config`\trequired but found unset")
		if err != nil {
			panic(err)
		}
		cliParseErrors = true
	} else {
		_, err := os.Stat(*nginxConfig)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				_, err := fmt.Fprintf(flag.CommandLine.Output(), "\"%s\" file does not exists\n", *nginxConfig)
				if err != nil {
					panic(err)
				}
				cliParseErrors = true
				cliErrorCode = 2
			} else {
				panic(err)
			}
		}
	}
	if cliParseErrors {
		flag.Usage()
		os.Exit(cliErrorCode)
	}

	_, err := fmt.Fprintf(flag.CommandLine.Output(), "server_name = \"%s\"\n", *serverName)
	if err != nil {
		panic(err)
	}

	directives, err := nginxparser.New(nil).ParseFile(*nginxConfig)
	if err != nil {
		panic(err)
	}

	body, err := json.MarshalIndent(directives, "", "  ")
	if err != nil {
		panic(err)
	}

	if *inplace {
		stat, err := os.Stat(*nginxConfig)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(*nginxConfig, body, stat.Mode().Perm())
		if err != nil {
			panic(err)
		}
	} else {
		_, err := fmt.Fprintln(os.Stdout, string(body))
		if err != nil {
			panic(err)
		}
	}
}
