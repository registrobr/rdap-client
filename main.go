package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/gregjones/httpcache"
	"github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/gregjones/httpcache/diskcache"
	"github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/registrobr/rdap"
	"github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/registrobr/rdap/protocol"
	"github.com/registrobr/rdap-client/output"
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
			Value: rdap.IANABootstrap,
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
			Usage: "force query for an IP or IPNetwork object",
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
		httpClient          = &http.Client{}
		uris                []string
	)

	forceCount := 0
	forceObjects := []bool{
		forceDomain,
		forceIP,
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
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: skipTLSVerification,
			},
		}

		httpClient.Transport = transport
	}

	if len(host) > 0 {
		if _, err := url.Parse(host); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		uris = []string{host}
	}

	identifier := strings.Join(ctx.Args(), " ")
	if identifier == "" {
		cli.ShowAppHelp(ctx)
		os.Exit(1)
	}

	var client rdap.Client
	client.URIs = uris

	if len(client.URIs) == 0 {
		cacheDetector := rdap.CacheDetector(func(resp *http.Response) bool {
			return resp.Header.Get(httpcache.XFromCache) == "1"
		})

		client.Transport = rdap.NewBootstrapFetcher(httpClient, bootstrapURI, cacheDetector)

	} else {
		client.Transport = rdap.NewDefaultFetcher(httpClient)
	}

	var err error
	var object interface{}

	switch {
	case forceASN:
		var asn uint64
		if asn, err = strconv.ParseUint(identifier, 10, 32); err == nil {
			object, err = client.ASN(uint32(asn), nil)
		}

	case forceDomain:
		object, err = client.Domain(identifier, nil)

	case forceEntity:
		object, err = client.Entity(identifier, nil)

	case forceIP:
		if ip := net.ParseIP(identifier); ip != nil {
			object, err = client.IP(ip, nil)
		} else {
			var ipnetwork *net.IPNet

			if _, ipnetwork, err = net.ParseCIDR(identifier); err != nil {
				err = fmt.Errorf("invalid ip or ip network “%s”", identifier)
			} else {
				object, err = client.IPNetwork(ipnetwork, nil)
			}
		}

	default:
		object, err = client.Query(identifier, nil)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var printer output.Printer

	switch object.(type) {
	case *protocol.AS:
		printer = &output.AS{
			AS: object.(*protocol.AS),
		}
	case *protocol.Domain:
		printer = &output.Domain{
			Domain: object.(*protocol.Domain),
		}
	case *protocol.Entity:
		printer = &output.Entity{
			Entity: object.(*protocol.Entity),
		}
	case *protocol.IPNetwork:
		printer = &output.IPNetwork{
			IPNetwork: object.(*protocol.IPNetwork),
		}
	}

	if err := printer.Print(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)
}
