package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"

	"github.com/registrobr/rdap-client/bootstrap"
	"github.com/registrobr/rdap-client/client"
	"github.com/registrobr/rdap-client/output"
)

type CLI struct {
	uris       []string
	httpClient *http.Client
	bootstrap  *bootstrap.Client
	wr         io.Writer
}

type handler func(object string) (bool, error)

func (c *CLI) asn() handler {
	return func(object string) (bool, error) {
		uris := c.uris

		if asn, err := strconv.ParseUint(object, 10, 32); err == nil {
			if c.bootstrap != nil {
				var err error
				uris, err = c.bootstrap.ASN(asn)

				if err != nil {
					return true, err
				}
			}

			r, err := client.NewClient(uris, c.httpClient).ASN(asn)

			if err != nil {
				return true, err
			}

			as := output.AS{AS: r}
			if err := as.ToText(c.wr); err != nil {
				return true, err
			}

			return true, nil
		}

		return false, nil
	}
}

func (c *CLI) entity() handler {
	// Note that there is no bootstrap for entity, see [1]
	// [1] - https://tools.ietf.org/html/rfc7484#section-6
	return func(object string) (bool, error) {
		r, err := client.NewClient(c.uris, c.httpClient).Entity(object)
		if err != nil {
			return true, err
		}

		entity := output.Entity{Entity: r}
		if err := entity.ToText(c.wr); err != nil {
			return true, err
		}
		return true, nil

	}
}

func (c *CLI) ipnetwork() handler {
	return func(object string) (bool, error) {
		uris := c.uris

		if ip := net.ParseIP(object); ip != nil {
			if c.bootstrap != nil {
				var err error
				uris, err = c.bootstrap.IP(ip)

				if err != nil {
					return true, err
				}
			}

			r, err := client.NewClient(uris, c.httpClient).IP(ip)
			if err != nil {
				return true, err
			}

			ipNetwork := output.IPNetwork{IPNetwork: r}
			if err := ipNetwork.ToText(c.wr); err != nil {
				return true, err
			}

			return true, nil
		}

		if _, cidr, err := net.ParseCIDR(object); err == nil {
			if c.bootstrap != nil {
				var err error
				uris, err = c.bootstrap.IPNetwork(cidr)

				if err != nil {
					return true, err
				}
			}

			r, err := client.NewClient(uris, c.httpClient).IPNetwork(cidr)

			if err != nil {
				return true, err
			}

			ipNetwork := output.IPNetwork{IPNetwork: r}
			if err := ipNetwork.ToText(c.wr); err != nil {
				return true, err
			}

			return true, nil
		}

		return false, nil
	}
}

func (c *CLI) domain() handler {
	return func(object string) (bool, error) {
		uris := c.uris

		if c.bootstrap != nil {
			var err error
			uris, err = c.bootstrap.Domain(object)

			if err != nil {
				return true, err
			}
		}

		r, err := client.NewClient(uris, c.httpClient).Domain(object)

		if err != nil {
			return true, err
		}

		if r == nil {
			return false, nil
		}

		if err := output.PrintDomain(r, c.wr); err != nil {
			return true, err
		}

		return false, nil
	}
}
