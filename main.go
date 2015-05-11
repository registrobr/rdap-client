package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	rdap "github.com/registrobr/rdap-client"
	"github.com/registrobr/rdap-client/output"
)

var (
	cache    = flag.String("cache", os.Getenv("HOME")+"/.rdap", "directory for caching bootstrap and RDAP data")
	endpoint = flag.String("endpoint", rdap.IANARDAPEndpoint, "endpoint for bootstrap data acquisition")
)

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		fmt.Println("rdap [options] object")
		os.Exit(1)
	}

	object := strings.Join(flag.Args(), " ")

	c := rdap.NewClient(*cache)

	if len(*endpoint) > 0 {
		c.SetRDAPEndpoint(*endpoint)
	}

	if asn, err := strconv.ParseUint(object, 10, 32); err == nil {
		r, err := c.QueryASN(asn)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		spew.Dump(r)
		os.Exit(0)
	}

	if _, cidr, err := net.ParseCIDR(object); err == nil {
		r, err := c.QueryIPNetwork(cidr)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		spew.Dump(r)
		os.Exit(0)
	}

	r, err := c.QueryDomain(object)

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
