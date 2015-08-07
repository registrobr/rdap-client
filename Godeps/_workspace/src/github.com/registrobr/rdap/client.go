package rdap

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/miekg/dns/idn"
	"github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/registrobr/rdap/protocol"
)

const (
	domain kind = "domain"
	autnum kind = "autnum"
	ip     kind = "ip"
	entity kind = "entity"
)

var (
	// ErrNotFound is used when the RDAP server doesn't contain any
	// information of the requested object
	ErrNotFound = errors.New("not found")
)

// Client queries a RDAP server to retrieve information of desired objects.
// You can also set the X-Forwarded-For HTTP header to work as a proxy
type Client interface {
	Domain(fqdn string) (*protocol.Domain, error)
	ASN(as uint64) (*protocol.AS, error)
	Entity(identifier string) (*protocol.Entity, error)
	IPNetwork(ipnet *net.IPNet) (*protocol.IPNetwork, error)
	IP(netIP net.IP) (*protocol.IPNetwork, error)
	SetURIs(uris []string)
	SetXForwardedFor(addr string)
}

// NewClient returns a client with the injected RDAP servers and HTTP client
var NewClient = func(uris []string, httpClient *http.Client) Client {
	return &client{
		uris:       uris,
		httpClient: httpClient,
	}
}

// client stores the HTTP client and the RDAP servers to query for retrieving
// the desired information. You can also set the X-Forward-For to work as a
// proxy
type client struct {
	httpClient    *http.Client
	uris          []string
	xForwardedFor string
}

// Domain will query each RDAP server to retrieve the desired information and
// will parse and store the response into a protocol Domain object. If
// something goes wrong an error will be returned, and if nothing is found
// the error ErrNotFound will be returned
func (c *client) Domain(fqdn string) (*protocol.Domain, error) {
	r := &protocol.Domain{}
	fqdn = idn.ToPunycode(strings.ToLower(fqdn))

	if err := c.query(domain, fqdn, r); err != nil {
		return nil, err
	}

	return r, nil
}

// ASN will query each RDAP server to retrieve the desired information and
// will parse and store the response into a protocol AS object. If
// something goes wrong an error will be returned, and if nothing is found
// the error ErrNotFound will be returned
func (c *client) ASN(as uint64) (*protocol.AS, error) {
	r := &protocol.AS{}

	if err := c.query(autnum, as, r); err != nil {
		return nil, err
	}

	return r, nil
}

// Entity will query each RDAP server to retrieve the desired information and
// will parse and store the response into a protocol Entity object. If
// something goes wrong an error will be returned, and if nothing is found
// the error ErrNotFound will be returned
func (c *client) Entity(identifier string) (*protocol.Entity, error) {
	r := &protocol.Entity{}

	if err := c.query(entity, identifier, r); err != nil {
		return nil, err
	}

	return r, nil
}

// IPNetwork will query each RDAP server to retrieve the desired information and
// will parse and store the response into a protocol IPNetwork object. If
// something goes wrong an error will be returned, and if nothing is found
// the error ErrNotFound will be returned
func (c *client) IPNetwork(ipnet *net.IPNet) (*protocol.IPNetwork, error) {
	r := &protocol.IPNetwork{}

	if err := c.query(ip, ipnet, r); err != nil {
		return nil, err
	}

	return r, nil
}

// IP will query each RDAP server to retrieve the desired information and
// will parse and store the response into a protocol IP object. If
// something goes wrong an error will be returned, and if nothing is found
// the error ErrNotFound will be returned
func (c *client) IP(netIP net.IP) (*protocol.IPNetwork, error) {
	r := &protocol.IPNetwork{}

	if err := c.query(ip, netIP, r); err != nil {
		return nil, err
	}

	return r, nil
}

func (c *client) query(kind kind, identifier interface{}, object interface{}) error {
	var lastErr error
	for _, uri := range c.uris {
		uri := fmt.Sprintf("%s/%s/%v", uri, kind, identifier)

		res, err := c.fetch(uri)
		if err != nil {
			lastErr = err
			continue
		}
		defer res.Body.Close()

		if err := handleHTTPStatusCode(kind, res); err != nil {
			lastErr = err
			continue
		}

		if err = json.NewDecoder(res.Body).Decode(&object); err != nil {
			lastErr = err
			continue
		}

		return nil
	}

	return lastErr
}

func (c *client) fetch(uri string) (response *http.Response, err error) {
	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	if c.xForwardedFor != "" {
		req.Header.Add("X-Forwarded-For", c.xForwardedFor)
	}

	if c.httpClient == nil {
		c.httpClient = &http.Client{}
	}

	response, err = c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *client) SetURIs(uris []string) {
	c.uris = uris
}

func (c *client) SetXForwardedFor(addr string) {
	c.xForwardedFor = addr
}

func handleHTTPStatusCode(kind kind, response *http.Response) error {
	if response.StatusCode == http.StatusOK {
		return nil
	}

	if response.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}

	if response.Header.Get("Content-Type") != "application/rdap+json" {
		return fmt.Errorf("unexpected response: %d %s",
			response.StatusCode, http.StatusText(response.StatusCode))
	}

	var responseErr protocol.Error
	if err := json.NewDecoder(response.Body).Decode(&responseErr); err != nil {
		return err
	}
	return responseErr
}
