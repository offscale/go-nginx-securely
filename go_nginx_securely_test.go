package main

import (
	"github.com/aluttik/go-crossplane"
	"log"
	"os"
	"testing"
)

func boolalpha(v bool) string {
	if v {
		return "true"
	} else {
		return "false"
	}
}

func TestHasUnsecuredServerName(t *testing.T) {
	unsecureServerExists := `server { server_name foo.com; listen 80; ssl off; }`
	unsecureServerNonexistent := `server { server_name bar.com; listen 80; ssl off; }`
	secureServerExists0 := `server { server_name foo.com; listen 443; ssl on; }`
	secureServerExists1 := `server { server_name foo.com; listen 443 ssl; }`

	for server, expectation := range map[string]bool{
		unsecureServerExists:      true,
		unsecureServerNonexistent: false,
		secureServerExists0:       false,
		secureServerExists1:       false,
	} {
		file, err := os.CreateTemp("", "temp_file")
		if err != nil {
			log.Fatal(err)
		}
		defer func(name string) {
			err := os.Remove(name)
			if err != nil {
				panic(err)
			}
		}(file.Name())
		err = os.WriteFile(file.Name(), []byte(server), 0)
		if err != nil {
			panic(err)
		}

		payload, err := crossplane.Parse(file.Name(), &crossplane.ParseOptions{SkipDirectiveContextCheck: true, SkipDirectiveArgsCheck: true})
		if err != nil {
			panic(err)
		}
		if hasUnsecuredServerName("foo.com", payload.Config) != expectation {
			t.Fatalf(`%s unexpectedly from hasUnsecuredServerName with "%s"`, boolalpha(!expectation), server)
		}
	}
}
