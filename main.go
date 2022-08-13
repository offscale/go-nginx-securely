package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	crossplane "github.com/aluttik/go-crossplane"
)

func main() {
	serverName, nginxConfig, inplace, sslCertificate, sslCertificateKey, sslDhparam := cliParser()

	/*directives, err := nginxparser.New(nil).ParseFile(*nginxConfig)
	if err != nil {
		panic(err)
	}

	body, err := json.MarshalIndent(directives, "", "  ")
	if err != nil {
		panic(err)
	}*/

	payload, err := crossplane.Parse(*nginxConfig, &crossplane.ParseOptions{SkipDirectiveContextCheck: true, SkipDirectiveArgsCheck: true})
	if err != nil {
		panic(err)
	}

	b, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	if hasUnsecuredServerName(*serverName, payload.Config) {
		redirectServerBlock := getRedirectServerBlock(*serverName)
		secureVarsBlock := getSecureVars(*serverName, *sslCertificate, *sslCertificateKey, *sslDhparam)
		secureConfig(
			&payload.Config,
			*serverName,
			&redirectServerBlock,
			&secureVarsBlock,
		)
		payload.Config[0].File = filepath.Join( /*os.TempDir()*/ "/tmp", path.Base(payload.Config[0].File))
		fmt.Printf("payload.Config[0].File to \"%s\"\n", payload.Config[0].File)

		err = crossplane.BuildFiles(
			*payload,
			"/tmp/a",
			&crossplane.BuildOptions{},
		)
		if err != nil {
			panic(err)
		}
		fileContents, err := os.ReadFile(payload.Config[0].File)
		if err != nil {
			panic(err)
		}
		_, err = os.Stdout.Write(fileContents)
		if err != nil {
			panic(err)
		}
	} else {
		os.Exit(33)
	}

	if *inplace {
		stat, err := os.Stat(*nginxConfig)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(*nginxConfig, b, stat.Mode().Perm())
		if err != nil {
			panic(err)
		}
	} /*else {
		_, err := fmt.Fprintln(os.Stdout, string(b))
		if err != nil {
			panic(err)
		}
	}*/
}

func cliParser() (*string, *string, *bool, *string, *string, *string) {
	serverName := flag.String("server_name", "", "server_name, only this one will be secured")
	nginxConfig := flag.String("nginx_config", "", "nginx config filepath")
	inplace := flag.Bool("inplace", false, "whether to edit inplace (defaults to stdout emission)")
	sslCertificate := flag.String("ssl_certificate", "/etc/letsencrypt/live/${server_name}/fullchain.pem",
		"SSL certificate, defaults to LetsEncrypt of `server_name`")
	sslCertificateKey := flag.String("ssl_certificate_key", "/etc/letsencrypt/live/${server_name}/privkey.pem",
		"SSL certificate key, defaults to LetsEncrypt of `server_name`")
	sslDhparam := flag.String("ssl_dh_param", "/etc/ssl/certs/dhparam.pem",
		"SSL Diffie-Helmann, defaults to nginx default")
	flag.Parse()

	cliParseErrors := 0
	cliErrorCode := syscall.Errno(1)

	if *serverName == "" {
		_, err := fmt.Fprintln(flag.CommandLine.Output(), "`server_name`\trequired but found unset")
		if err != nil {
			panic(err)
		}
		cliParseErrors++
	}
	if *nginxConfig == "" {
		_, err := fmt.Fprintln(flag.CommandLine.Output(), "`nginx_config`\trequired but found unset")
		if err != nil {
			panic(err)
		}
		cliParseErrors++
	} else {
		_, err := os.Stat(*nginxConfig)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				_, err := fmt.Fprintf(flag.CommandLine.Output(), "\"%s\" file does not exists\n", *nginxConfig)
				if err != nil {
					panic(err)
				}
				cliParseErrors++
				cliErrorCode = syscall.ENOENT
			} else {
				panic(err)
			}
		}
	}
	if cliParseErrors > 0 {
		if cliParseErrors > 1 || cliErrorCode != 2 {
			flag.Usage()
		}
		os.Exit(int(cliErrorCode))
	}
	*sslCertificate = strings.Replace(*sslCertificate, "${server_name}", *serverName, 1)
	*sslCertificateKey = strings.Replace(*sslCertificateKey, "${server_name}", *serverName, 1)
	*sslDhparam = strings.Replace(*sslDhparam, "${server_name}", *serverName, 1)
	return serverName, nginxConfig, inplace, sslCertificate, sslCertificateKey, sslDhparam
}
