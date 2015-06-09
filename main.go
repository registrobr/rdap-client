package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	cgcli "github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/gregjones/httpcache"
	"github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/gregjones/httpcache/diskcache"
	"github.com/registrobr/rdap-client/bootstrap"
	"github.com/registrobr/rdap-client/handler"
)

func main() {
	cgcli.AppHelpTemplate = `
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

	app := cgcli.NewApp()
	app.Name = "rdap"
	app.Usage = "RDAP cgclient"
	app.Author = "NIC.br"
	app.Version = "0.0.1"

	app.Flags = []cgcli.Flag{
		cgcli.StringFlag{
			Name:  "cache",
			Value: path.Join(os.Getenv("HOME"), ".rdap"),
			Usage: "directory for caching bootstrap and RDAP data",
		},
		cgcli.StringFlag{
			Name:  "bootstrap",
			Value: bootstrap.RDAPBootstrap,
			Usage: "RDAP bootstrap service URL",
		},
		cgcli.BoolFlag{
			Name:  "no-cache",
			Usage: "don't cache anything",
		},
		cgcli.BoolFlag{
			Name:  "skip-tls-verification,S",
			Usage: "skip TLS verification",
		},
		cgcli.BoolFlag{
			Name:  "domain",
			Usage: "force query for a domain object",
		},
		cgcli.BoolFlag{
			Name:  "asn",
			Usage: "force query for an ASN object",
		},
		cgcli.BoolFlag{
			Name:  "ip",
			Usage: "force query for an IP object",
		},
		cgcli.BoolFlag{
			Name:  "ipnetwork",
			Usage: "force query for an IP Network object",
		},
		cgcli.BoolFlag{
			Name:  "entity",
			Usage: "force query for an Entity object",
		},
		cgcli.StringFlag{
			Name:  "host,H",
			Value: "",
			Usage: "host where to send the query (bypass bootstrap)",
		},
	}

	app.Commands = []cgcli.Command{}
	app.Action = action

	app.Run(os.Args)
}

func action(ctx *cgcli.Context) {
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
		bs                  *bootstrap.Client
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
		bs = bootstrap.NewClient(httpClient)

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
		cgcli.ShowAppHelp(ctx)
		os.Exit(1)
	}

	var (
		ok  bool
		err error
		h   = &handler.Handler{
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
