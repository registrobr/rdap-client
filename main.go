package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/gregjones/httpcache"
	"github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/gregjones/httpcache/diskcache"
	"github.com/registrobr/rdap-client/bootstrap"
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

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "cache",
			Value: path.Join(os.Getenv("HOME"), ".rdap"),
			Usage: "directory for caching bootstrap and RDAP data",
		},
		cli.StringFlag{
			Name:  "bootstrap",
			Value: bootstrap.RDAPBootstrap,
			Usage: "RDAP bootstrap service URL",
		},
		cli.BoolFlag{
			Name:  "no-cache",
			Usage: "Don't cache anything",
		},
		cli.StringFlag{
			Name:  "host,H",
			Value: "",
			Usage: "Host where to send the query (bypass bootstrap)",
		},
	}

	app.Commands = []cli.Command{}
	app.Action = action

	app.Run(os.Args)
}

func action(ctx *cli.Context) {
	var (
		cache        = ctx.String("cache")
		bootstrapURI = ctx.String("bootstrap")
		host         = ctx.String("host")
		httpClient   = &http.Client{}
		bs           *bootstrap.Client
		uris         []string
	)

	if !ctx.Bool("no-cache") {
		httpClient.Transport = httpcache.NewTransport(
			diskcache.New(cache),
		)
	}

	if len(host) == 0 {
		bs = bootstrap.NewClient(httpClient)

		if len(bootstrapURI) > 0 {
			bs.Bootstrap = bootstrapURI
		}
	} else {
		if _, err := url.Parse(host); err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}

		uris = []string{host}
	}

	object := strings.Join(ctx.Args(), " ")
	if object == "" {
		cli.ShowAppHelp(ctx)
		os.Exit(1)
	}

	rdapCLI := &CLI{
		uris:       uris,
		httpClient: httpClient,
		bootstrap:  bs,
		wr:         os.Stdout,
	}

	handlers := []handler{
		rdapCLI.asn(),
		rdapCLI.ipnetwork(),
		rdapCLI.domain(),
	}

	var err error
	var executed bool
	for _, handler := range handlers {
		executed, err = handler(object)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}

		if executed {
			break
		}
	}

	if !executed {
		_, err := rdapCLI.entity()(object)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	}

	os.Exit(0)
}
