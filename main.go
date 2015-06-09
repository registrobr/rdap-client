package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/gregjones/httpcache"
	"github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/gregjones/httpcache/diskcache"
	"github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/registrobr/rdap"
)

func main() {
	cli.AppHelpTemplate = `
NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.Name}} {{if .Flags}}[global options]{{end}} OBJECT

VERSION:
   {{.Version}}{{if len .Authors}}

AUTHOR(S): 
   {{range .Authors}}{{ . }}{{end}}{{end}}

GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}
`

	app := cli.NewApp()
	app.Name = "rdap"
	app.Usage = "RDAP client"
	app.Author = "NIC.br"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "cache",
			Value: path.Join(os.Getenv("HOME"), ".rdap"),
			Usage: "directory for caching bootstrap and RDAP data",
		},
		cli.StringFlag{
			Name:  "bootstrap",
			Value: client.RDAPBootstrap,
			Usage: "RDAP bootstrap service URL",
		},
		cli.BoolFlag{
			Name:  "no-cache",
			Usage: "don't cache anything",
		},
		cli.BoolFlag{
			Name:  "skip-tls-verification,S",
			Usage: "skip TLS verification",
		},
		cli.BoolFlag{
			Name:  "domain",
			Usage: "force query for a domain object",
		},
		cli.BoolFlag{
			Name:  "asn",
			Usage: "force query for an ASN object",
		},
		cli.BoolFlag{
			Name:  "ip",
			Usage: "force query for an IP object",
		},
		cli.BoolFlag{
			Name:  "ipnetwork",
			Usage: "force query for an IP Network object",
		},
		cli.BoolFlag{
			Name:  "entity",
			Usage: "force query for an Entity object",
		},
		cli.StringFlag{
			Name:  "host,H",
			Value: "",
			Usage: "host where to send the query (bypass bootstrap)",
		},
	}

	app.Commands = []cli.Command{}
	app.Action = action

	app.Run(os.Args)
}

func action(ctx *cli.Context) {
	var (
		cache               = ctx.String("cache")
		bootstrapURI        = ctx.String("bootstrap")
		host                = ctx.String("host")
		skipTLSVerification = ctx.Bool("skip-tls-verification")
		forceASN            = ctx.Bool("asn")
		forceDomain         = ctx.Bool("domain")
		forceEntity         = ctx.Bool("entity")
		forceIP             = ctx.Bool("ip")
		forceIPNetwork      = ctx.Bool("ipnetwork")
		force               = forceASN || forceDomain || forceEntity || forceIP || forceIPNetwork
		httpClient          = &http.Client{}
		bs                  *client.Bootstrap
		uris                []string
	)

	forceCount := 0
	forceObjects := []bool{
		forceDomain,
		forceIP,
		forceIPNetwork,
		forceEntity,
		forceASN,
	}

	for _, force := range forceObjects {
		if force {
			if forceCount++; forceCount > 1 {
				fmt.Fprintln(os.Stderr, "you can't use -asn, -domain, -entity, -ip or -ipnetwork at the same time")
				os.Exit(1)
			}
		}
	}

	if !ctx.Bool("no-cache") {
		transport := httpcache.NewTransport(
			diskcache.New(cache),
		)

		transport.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: skipTLSVerification},
		}

		httpClient.Transport = transport
	}

	if len(host) == 0 {
		bs = client.NewBootstrap(httpClient)

		if len(bootstrapURI) > 0 {
			bs.Bootstrap = bootstrapURI
		}
	} else {
		if _, err := url.Parse(host); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		uris = []string{host}
	}

	object := strings.Join(ctx.Args(), " ")
	if object == "" {
		cli.ShowAppHelp(ctx)
		os.Exit(1)
	}

	var (
		ok  bool
		err error
		h   = &client.Handler{
			URIs:       uris,
			HTTPClient: httpClient,
			Bootstrap:  bs,
			Writer:     os.Stdout,
		}
	)

	switch {
	case forceASN:
		ok, err = h.ASN(object)
	case forceDomain:
		ok, err = h.Domain(object)
	case forceEntity:
		ok, err = h.Entity(object)
	case forceIP:
		ok, err = h.IP(object)
	case forceIPNetwork:
		ok, err = h.IPNetwork(object)
	default:
		ok = true
		err = h.Query(object)
	}

	if err == nil && !ok {
		if force {
			err = fmt.Errorf("the requested object doesn't match the requested object type")
		} else {
			err = fmt.Errorf("the requested object doesn't match an ASN, Domain, Entity, IP or IPNetwork")
		}
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)
}
