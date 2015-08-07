package rdap

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"sort"
)

type kind string

const (
	dns  kind = "dns"
	asn  kind = "asn"
	ipv4 kind = "ipv4"
	ipv6 kind = "ipv6"

	// RDAPBootstrap stores the default URL to query to retrieve the RDAP
	// servers that contain the desired information
	RDAPBootstrap = "https://data.iana.org/rdap/%s.json"
)

// Bootstrap stores all necessary information to query IANA and retrieve the
// RDAP server that contains the desired information
type Bootstrap struct {
	Bootstrap   string
	IsFromCache func(*http.Response) bool
	httpClient  *http.Client
	cacheKey    string
	reloadCache bool
}

// NewBootstrap returns build a Bootstrap instance that will use the injected
// http client. By default the bootstrap will query the URI stored in the
// RDAPBootstrap constant
func NewBootstrap(httpClient *http.Client) *Bootstrap {
	return &Bootstrap{
		Bootstrap:  RDAPBootstrap,
		httpClient: httpClient,
	}
}

// Domain retrieve the RDAP servers that can contain information about the
// given domain. If something goes wrong an error is returned
func (c *Bootstrap) Domain(fqdn string) ([]string, error) {
	return c.query(dns, fqdn)
}

// ASN retrieve the RDAP servers that can contain information about the
// given ASN. If something goes wrong an error is returned
func (c *Bootstrap) ASN(as uint64) ([]string, error) {
	return c.query(asn, as)
}

// IPNetwork retrieve the RDAP servers that can contain information about the
// given IP network. If something goes wrong an error is returned
func (c *Bootstrap) IPNetwork(ipnet *net.IPNet) ([]string, error) {
	kind := ipv4

	if ipnet.IP.To4() == nil {
		kind = ipv6
	}

	return c.query(kind, ipnet)
}

// IP retrieve the RDAP servers that can contain information about the
// given IP. If something goes wrong an error is returned
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

	if len(uris) > 0 || !cached {
		return
	}

	if nsSet, err := net.LookupNS(fqdn); err != nil || len(nsSet) == 0 {
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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotModified {
		return nil, false, fmt.Errorf("unexpected status code %d %s",
			resp.StatusCode,
			http.StatusText(resp.StatusCode),
		)
	}

	if c.IsFromCache != nil {
		cached = c.IsFromCache(resp)
	}

	return resp.Body, cached, nil
}
