package main

import (
	"io"
	"net"
	"net/http"
	"regexp"
	"strconv"

	"github.com/registrobr/rdap-client/bootstrap"
	"github.com/registrobr/rdap-client/client"
	"github.com/registrobr/rdap-client/output"
)

var isFQDN = regexp.MustCompile(`^(([[:alnum:]](([[:alnum:]]|\-){0,61}[[:alnum:]])?\.)*[[:alnum:]](([[:alnum:]]|\-){0,61}[[:alnum:]])?)?(\.)?$`)

type cli struct {
	uris       []string
	httpClient *http.Client
	bootstrap  *bootstrap.Client
	wr         io.Writer
}

type handler func(object string) (bool, error)

func (c *cli) guess(object string) (bool, error) {
	handlers := []handler{
		c.asn(),
		c.ip(),
		c.ipnetwork(),
		c.domain(),
		c.entity(),
	}

	ok := false

	for _, handler := range handlers {
		var err error
		ok, err = handler(object)

		if err != nil {
			return ok, err
		}

		if ok {
			break
		}
	}

	return ok, nil
}

func (c *cli) asn() handler {
	return func(object string) (bool, error) {
		asn, err := strconv.ParseUint(object, 10, 32)

		if err != nil {
			return false, nil
		}

		uris := c.uris

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
}

func (c *cli) entity() handler {
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

func (c *cli) ipnetwork() handler {
	return func(object string) (bool, error) {
		_, cidr, err := net.ParseCIDR(object)

		if err != nil {
			return false, nil
		}

		uris := c.uris

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
}

func (c *cli) ip() handler {
	return func(object string) (bool, error) {
		ip := net.ParseIP(object)

		if ip == nil {
			return false, nil
		}

		uris := c.uris

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
}

func (c *cli) domain() handler {
	return func(object string) (bool, error) {
		if !isFQDN.MatchString(object) {
			return false, nil
		}

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
			return true, nil
		}

		domain := output.Domain{Domain: r}
		if err := domain.ToText(c.wr); err != nil {
			return true, err
		}

		return true, nil
	}
}
