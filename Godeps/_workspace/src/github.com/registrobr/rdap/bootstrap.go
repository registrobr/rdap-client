package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"sort"

	"github.com/registrobr/rdap/Godeps/_workspace/src/github.com/gregjones/httpcache"
)

type kind string

const (
	dns           kind = "dns"
	asn           kind = "asn"
	ipv4          kind = "ipv4"
	ipv6          kind = "ipv6"
	RDAPBootstrap      = "https://data.iana.org/rdap/%s.json"
)

type Bootstrap struct {
	httpClient  *http.Client
	cacheKey    string
	Bootstrap   string
	reloadCache bool
}

func NewBootstrap(httpClient *http.Client) *Bootstrap {
	return &Bootstrap{
		Bootstrap:  RDAPBootstrap,
		httpClient: httpClient,
	}
}

func (c *Bootstrap) Domain(fqdn string) ([]string, error) {
	return c.query(dns, fqdn)
}

func (c *Bootstrap) ASN(as uint64) ([]string, error) {
	return c.query(asn, as)
}

func (c *Bootstrap) IPNetwork(ipnet *net.IPNet) ([]string, error) {
	kind := ipv4

	if ipnet.IP.To4() == nil {
		kind = ipv6
	}

	return c.query(kind, ipnet)
}

func (c *Bootstrap) IP(ip net.IP) ([]string, error) {
	kind := ipv4

	if ip.To4() == nil {
		kind = ipv6
	}

	return c.query(kind, ip)
}

func (c *Bootstrap) checkDomain(fqdn string, cached bool, r serviceRegistry) (uris []string, err error) {
	uris, err = r.matchDomain(fqdn)
	if err != nil {
		return
	}

	if len(uris) > 0 {
		return
	}

	if !cached {
		return
	}

	nsSet, err := net.LookupNS(fqdn)
	if err != nil {
		return nil, nil
	}

	if len(nsSet) == 0 {
		return nil, nil
	}

	c.reloadCache = true
	body, cached, err := c.fetch(fmt.Sprintf(c.Bootstrap, dns))
	if err != nil {
		return nil, err
	}
	defer body.Close()

	if err := json.NewDecoder(body).Decode(&r); err != nil {
		return nil, err
	}

	return r.matchDomain(fqdn)
}

func (c *Bootstrap) query(kind kind, identifier interface{}) ([]string, error) {
	uris := []string{}
	r := serviceRegistry{}
	uri := fmt.Sprintf(c.Bootstrap, kind)
	body, cached, err := c.fetch(uri)

	if err != nil {
		return nil, err
	}

	defer body.Close()

	if err := json.NewDecoder(body).Decode(&r); err != nil {
		return nil, err
	}

	if r.Version != version {
		return nil, fmt.Errorf("incompatible bootstrap specification version: %s (expecting %s)", r.Version, version)
	}

	switch kind {
	case dns:
		uris, err = c.checkDomain(identifier.(string), cached, r)
	case asn:
		uris, err = r.matchAS(identifier.(uint64))
	case ipv4, ipv6:
		if ip, ok := identifier.(net.IP); ok {
			uris, err = r.matchIP(ip)
			break
		}

		if ipNet, ok := identifier.(*net.IPNet); ok {
			uris, err = r.matchIPNetwork(ipNet)
			break
		}
	}

	if err != nil {
		return nil, err
	}

	if len(uris) == 0 {
		return nil, fmt.Errorf("no matches for %v", identifier)
	}

	sort.Sort(prioritizeHTTPS(uris))

	return uris, nil
}

func (c *Bootstrap) fetch(uri string) (_ io.ReadCloser, cached bool, err error) {
	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		return nil, cached, err
	}

	if c.reloadCache {
		req.Header.Add("Cache-Control", "max-age=0")
	}

	c.cacheKey = req.URL.String()

	if c.httpClient == nil {
		c.httpClient = &http.Client{}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, cached, err
	}

	if resp.Header.Get(httpcache.XFromCache) == "1" {
		cached = true
	}

	return resp.Body, cached, nil
}
