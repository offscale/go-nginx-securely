package main

import "github.com/aluttik/go-crossplane"

func hasUnsecuredServerName(serverName string, config []crossplane.Config) bool {
	hasServerName := false
	isSecure := false

	for _, conf := range config {
		for _, directive0 := range conf.Parsed {
			if directive0.Directive == "server" {
				for _, directive1 := range *directive0.Block {
					if directive1.Directive == "server_name" {
						if len(directive1.Args) > 0 {
							for _, arg := range directive1.Args {
								if arg == serverName {
									hasServerName = true
									break
								}
							}
						}
					} else if directive1.Directive == "listen" {
						if len(directive1.Args) > 0 {
							for _, arg := range directive1.Args {
								if arg == "443" || arg == "ssl" {
									isSecure = true
									break
								}
							}
						}
					} else if directive1.Directive == "ssl" {
						if len(directive1.Args) > 0 {
							for _, arg := range directive1.Args {
								if arg == "on" {
									isSecure = true
									break
								} else if arg == "off" {
									/* In case port 443 but ssl is off */
									isSecure = false
								}
							}
						}
					}
				}
			}
		}
	}

	return hasServerName && !isSecure
}

func getSecureVars(serverName string, sslCertificate string, sslCertificateKey string, sslDhparam string) []crossplane.Directive {
	/*TODO: generate the 'redirect block'*/
	/*TODO: move existing `listen` block to 'redirect block'*/
	return /*secureVars*/ []crossplane.Directive{
		{
			Directive: "server",
			Line:      1,
			Args:      []string{},
			Block: &[]crossplane.Directive{
				{
					Directive: "listen",
					Line:      2,
					Args: []string{
						"443",
						"ssl",
					},
				},
				{
					Directive: "server_name",
					Line:      5,
					Args: []string{
						serverName,
					},
				},
				{
					Directive: "ssl_certificate",
					Line:      7,
					Args: []string{
						sslCertificate,
					},
				},
				{
					Directive: "ssl_certificate_key",
					Line:      8,
					Args: []string{
						sslCertificateKey,
					},
				},
				{
					Directive: "ssl_protocols",
					Line:      10,
					Args: []string{
						"TLSv1",
						"TLSv1.1",
						"TLSv1.2",
					},
				},
				{
					Directive: "ssl_prefer_server_ciphers",
					Line:      11,
					Args: []string{
						"on",
					},
				},
				{
					Directive: "ssl_dhparam",
					Line:      12,
					Args: []string{
						sslDhparam,
					},
				},
				{
					Directive: "ssl_ciphers",
					Line:      13,
					Args: []string{
						"ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-AES256-GCM-SHA384:DHE-RSA-AES128-GCM-SHA256:DHE-DSS-AES128-GCM-SHA256:kEDH+AESGCM:ECDHE-RSA-AES128-SHA256:ECDHE-ECDSA-AES128-SHA256:ECDHE-RSA-AES128-SHA:ECDHE-ECDSA-AES128-SHA:ECDHE-RSA-AES256-SHA384:ECDHE-ECDSA-AES256-SHA384:ECDHE-RSA-AES256-SHA:ECDHE-ECDSA-AES256-SHA:DHE-RSA-AES128-SHA256:DHE-RSA-AES128-SHA:DHE-DSS-AES128-SHA256:DHE-RSA-AES256-SHA256:DHE-DSS-AES256-SHA:DHE-RSA-AES256-SHA:AES128-GCM-SHA256:AES256-GCM-SHA384:AES128-SHA256:AES256-SHA256:AES128-SHA:AES256-SHA:AES:CAMELLIA:DES-CBC3-SHA:!aNULL:!eNULL:!EXPORT:!DES:!RC4:!MD5:!PSK:!aECDH:!EDH-DSS-DES-CBC3-SHA:!EDH-RSA-DES-CBC3-SHA:!KRB5-DES-CBC3-SHA",
					},
				},
				{
					Directive: "ssl_session_timeout",
					Line:      14,
					Args: []string{
						"1d",
					},
				},
				{
					Directive: "ssl_session_cache",
					Line:      15,
					Args: []string{
						"shared:SSL:50m",
					},
				},
				{
					Directive: "ssl_stapling",
					Line:      16,
					Args: []string{
						"on",
					},
				},
				{
					Directive: "ssl_stapling_verify",
					Line:      17,
					Args: []string{
						"on",
					},
				},
				{
					Directive: "add_header",
					Line:      18,
					Args: []string{
						"Strict-Transport-Security",
						"max-age=15768000",
					},
				},
				{
					Directive: "location",
					Line:      20,
					Args: []string{
						"~",
						"/.well-known",
					},
					Block: &[]crossplane.Directive{
						{
							Directive: "allow",
							Line:      21,
							Args: []string{
								"all",
							},
						},
					},
				},
			},
		},
	}
}
