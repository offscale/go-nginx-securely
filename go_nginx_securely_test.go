package main

import (
	"github.com/aluttik/go-crossplane"
	"log"
	"os"
	"path"
	"testing"
)

func boolalpha(v bool) string {
	if v {
		return "true"
	} else {
		return "false"
	}
}

const unsecureServerExists string = `server { server_name foo.com; listen 80; ssl off; }`

func TestHasUnsecuredServerName(t *testing.T) {
	unsecureServerNonexistent := `server { server_name bar.com; listen 80; ssl off; }`
	secureServerExists0 := `server { server_name foo.com; listen 443; ssl on; }`
	secureServerExists1 := `server { server_name foo.com; listen 443 ssl; }`

	file, err := os.CreateTemp("", "temp_file")
	if err != nil {
		log.Fatal(err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Fatal(err)
		}
	}(file.Name())
	for server, expectation := range map[string]bool{
		unsecureServerExists:      true,
		unsecureServerNonexistent: false,
		secureServerExists0:       false,
		secureServerExists1:       false,
	} {
		err = os.WriteFile(file.Name(), []byte(server), 0)
		if err != nil {
			log.Fatal(err)
		}

		payload, err := crossplane.Parse(file.Name(), &crossplane.ParseOptions{SkipDirectiveContextCheck: true, SkipDirectiveArgsCheck: true})
		if err != nil {
			log.Fatal(err)
		}
		if hasUnsecuredServerName("foo.com", payload.Config) != expectation {
			t.Fatalf("%s unexpectedly from `hasUnsecuredServerName` with:\t%s", boolalpha(!expectation), server)
		}
	}
}

const secureVars string = `server {
    listen 443 ssl;
    server_name foo.com;
    ssl_certificate ssl.cert;
    ssl_certificate_key ssl.key;
    ssl_protocols TLSv1 TLSv1.1 TLSv1.2;
    ssl_prefer_server_ciphers on;
    ssl_dhparam dhparam.ssl;
    ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-AES256-GCM-SHA384:DHE-RSA-AES128-GCM-SHA256:DHE-DSS-AES128-GCM-SHA256:kEDH+AESGCM:ECDHE-RSA-AES128-SHA256:ECDHE-ECDSA-AES128-SHA256:ECDHE-RSA-AES128-SHA:ECDHE-ECDSA-AES128-SHA:ECDHE-RSA-AES256-SHA384:ECDHE-ECDSA-AES256-SHA384:ECDHE-RSA-AES256-SHA:ECDHE-ECDSA-AES256-SHA:DHE-RSA-AES128-SHA256:DHE-RSA-AES128-SHA:DHE-DSS-AES128-SHA256:DHE-RSA-AES256-SHA256:DHE-DSS-AES256-SHA:DHE-RSA-AES256-SHA:AES128-GCM-SHA256:AES256-GCM-SHA384:AES128-SHA256:AES256-SHA256:AES128-SHA:AES256-SHA:AES:CAMELLIA:DES-CBC3-SHA:!aNULL:!eNULL:!EXPORT:!DES:!RC4:!MD5:!PSK:!aECDH:!EDH-DSS-DES-CBC3-SHA:!EDH-RSA-DES-CBC3-SHA:!KRB5-DES-CBC3-SHA;
    ssl_session_timeout 1d;
    ssl_session_cache shared:SSL:50m;
    ssl_stapling on;
    ssl_stapling_verify on;
    add_header Strict-Transport-Security max-age=15768000;
    location ~ /.well-known {
        allow all;
    }
}
`

const redirectVars string = `server {
    listen 80;
    server_name foo.com;
    return 301 https://$server_name$request_uri;
}
` /* `listen [::]:80;` if exposing with IPv6 */

func TestGetSecureVars(t *testing.T) {
	DirectivesTestFixture(
		t,
		secureVars,
		getSecureVars("foo.com", "ssl.cert", "ssl.key", "dhparam.ssl"),
	)
}

func TestGetRedirectServerBlock(t *testing.T) {
	DirectivesTestFixture(
		t,
		redirectVars,
		[]crossplane.Directive{getRedirectServerBlock("foo.com")},
	)
}

func DirectivesTestFixture(t *testing.T, expectation string, directives []crossplane.Directive) {
	dir, err := os.MkdirTemp("", "DirectivesTestFixture")
	if err != nil {
		log.Fatal(err)
	}
	defer func(name string) {
		dirs, e := os.ReadDir(name)
		if e != nil {
			log.Fatal(e)
		}
		if len(dirs) == 0 {
			log.Fatal("Dir is empty")
		}
		for _, file := range dirs {
			absoluteFilename := path.Join(name, file.Name())
			buf, er := os.ReadFile(absoluteFilename)
			if er != nil {
				log.Fatal(er)
			}
			if string(buf) != expectation {
				t.Fatalf("Expected `%s` got `%s`", expectation, buf)
			}
			er = os.Remove(absoluteFilename)
			if er != nil {
				log.Fatal(er)
			}
		}
		err := os.Remove(name)
		if err != nil {
			log.Fatal(err)
		}
	}(dir)

	err = crossplane.BuildFiles(
		crossplane.Payload{
			Status: "ok",
			Errors: []crossplane.PayloadError{},
			Config: []crossplane.Config{
				{
					File:   "filename0.conf",
					Status: "ok",
					Errors: []crossplane.ConfigError{},
					Parsed: directives,
				},
			},
		},
		dir,
		&crossplane.BuildOptions{},
	)
	if err != nil {
		log.Fatal(err)
	}
}
