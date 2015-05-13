package main

import (
	"fmt"
	"net"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/davecgh/go-spew/spew"
	rdap "github.com/registrobr/rdap-client/client"
	"github.com/registrobr/rdap-client/output"
)

var (
	cache     string
	bootstrap string
	host      string
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
			Value: rdap.RDAPBootstrap,
			Usage: "RDAP bootstrap service URL",
		},
		cli.StringFlag{
			Name:  "host,H",
			Value: "",
			Usage: "Host where to send the query (bypass bootstrap)",
		},
	}
	app.Commands = []cli.Command{}

	app.Action = func(ctx *cli.Context) {
		cache = ctx.String("cache")
		bootstrap = ctx.String("bootstrap")
		host = ctx.String("host")

		client := rdap.NewClient(cache)

		if host != "" {
			client.Host = host
		}

		if len(bootstrap) > 0 {
			client.Bootstrap = bootstrap
		}

		object := strings.Join(ctx.Args(), " ")
		if object == "" {
			cli.ShowAppHelp(ctx)
			os.Exit(1)
		}

		if asn, err := strconv.ParseUint(object, 10, 32); err == nil {
			r, err := client.QueryASN(asn)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			spew.Dump(r)
			os.Exit(0)
		}

		if _, cidr, err := net.ParseCIDR(object); err == nil {
			r, err := client.QueryIPNetwork(cidr)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			spew.Dump(r)
			os.Exit(0)
		}

		r, err := client.QueryDomain(object)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := output.PrintDomain(r, os.Stdout); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		os.Exit(0)
	}

	app.Run(os.Args)
}
