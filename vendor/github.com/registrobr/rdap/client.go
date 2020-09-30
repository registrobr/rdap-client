package rdap

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/registrobr/rdap/protocol"
	"golang.org/x/net/idna"
)

var (
	fqdnRX = regexp.MustCompile(`^((([a-z0-9][a-z0-9\-]*[a-z0-9])|[a-z0-9]+)\.)*([a-z]+|xn\-\-[a-z0-9]+)\.?$`)
)

// Client is responsible for building, sending the request and parsing the
// result. It can set the URIs attribute if you want to query RDAP servers
// directly without using bootstrap
type Client struct {
	// Transport is the network layer that you can fill with a direct query to
	// the RDAP servers or with an extra layer of RDAP bootstrap strategy
	Transport Fetcher

	// URIs store the addresses of the RDAP servers that you want to query
	// directly. Remember that if you use a bootstrap transport layer this
	// information might not be used
	URIs []string
}

// NewClient is an easy way to create a client with bootstrap support or not,
// depending if you inform direct RDAP addresses
func NewClient(URIs []string) *Client {
	client := Client{
		URIs: URIs,
	}

	var httpClient http.Client

	if len(URIs) == 0 {
		client.Transport = NewBootstrapFetcher(&httpClient, IANABootstrap, nil)
	} else {
		client.Transport = NewDefaultFetcher(&httpClient)
	}

	return &client
}

// Domain will query each RDAP server to retrieve the desired information and
// will parse and store the response into a protocol Domain object. You can
// optionally define the HTTP headers parameters to send to the RDAP server. If
// something goes wrong an error will be returned, and if nothing is found
// the error ErrNotFound will be returned. The HTTP header of the RDAP
// response is also returned to analyze any specific flag
func (c *Client) Domain(fqdn string, header http.Header, queryString url.Values) (*protocol.Domain, http.Header, error) {
	fqdn, err := idna.ToASCII(strings.ToLower(fqdn))
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.Transport.Fetch(c.URIs, QueryTypeDomain, fqdn, header, queryString)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if err != nil {
		if resp != nil {
			return nil, resp.Header, err
		}
		return nil, nil, err
	}

	domain := &protocol.Domain{}
	if err = json.NewDecoder(resp.Body).Decode(domain); err != nil {
		return nil, resp.Header, err
	}

	return domain, resp.Header, nil
}

// Ticket will query each RDAP server to retrieve the desired information and
// will parse and store the response into a protocol Domain object. You can
// optionally define the HTTP headers parameters to send to the RDAP server. If
// something goes wrong an error will be returned, and if nothing is found
// the error ErrNotFound will be returned. The HTTP header of the RDAP
// response is also returned to analyze any specific flag
func (c *Client) Ticket(ticketNumber int, header http.Header, queryString url.Values) (*protocol.Domain, http.Header, error) {
	resp, err := c.Transport.Fetch(c.URIs, QueryTypeTicket, strconv.Itoa(ticketNumber), header, queryString)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if err != nil {
		if resp != nil {
			return nil, resp.Header, err
		}
		return nil, nil, err
	}

	domain := &protocol.Domain{}
	if err = json.NewDecoder(resp.Body).Decode(domain); err != nil {
		return nil, resp.Header, err
	}

	return domain, resp.Header, nil
}

// ASN will query each RDAP server to retrieve the desired information and
// will parse and store the response into a protocol AS object. You can
// optionally define the HTTP headers parameters to send to the RDAP server.
// If something goes wrong an error will be returned, and if nothing is found
// the error ErrNotFound will be returned. The HTTP header of the RDAP
// response is also returned to analyze any specific flag
func (c *Client) ASN(asn uint32, header http.Header, queryString url.Values) (*protocol.AS, http.Header, error) {
	asnStr := strconv.FormatUint(uint64(asn), 10)

	resp, err := c.Transport.Fetch(c.URIs, QueryTypeAutnum, asnStr, header, queryString)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if err != nil {
		if resp != nil {
			return nil, resp.Header, err
		}
		return nil, nil, err
	}

	as := &protocol.AS{}
	if err = json.NewDecoder(resp.Body).Decode(as); err != nil {
		return nil, resp.Header, err
	}

	return as, resp.Header, nil
}

// Entity will query each RDAP server to retrieve the desired information and
// will parse and store the response into a protocol Entity object. You can
// optionally define the HTTP headers parameters to send to the RDAP server.
// If something goes wrong an error will be returned, and if nothing is found
// the error ErrNotFound will be returned. The HTTP header of the RDAP
// response is also returned to analyze any specific flag
func (c *Client) Entity(identifier string, header http.Header, queryString url.Values) (*protocol.Entity, http.Header, error) {
	resp, err := c.Transport.Fetch(c.URIs, QueryTypeEntity, identifier, header, queryString)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if err != nil {
		if resp != nil {
			return nil, resp.Header, err
		}
		return nil, nil, err
	}

	entity := &protocol.Entity{}
	if err = json.NewDecoder(resp.Body).Decode(entity); err != nil {
		return nil, resp.Header, err
	}

	return entity, resp.Header, nil
}

// IPNetwork will query each RDAP server to retrieve the desired information and
// will parse and store the response into a protocol IPNetwork object. You can
// optionally define the HTTP headers parameters to send to the RDAP server.
// If something goes wrong an error will be returned, and if nothing is found
// the error ErrNotFound will be returned. The HTTP header of the RDAP
// response is also returned to analyze any specific flag
func (c *Client) IPNetwork(ipnet *net.IPNet, header http.Header, queryString url.Values) (*protocol.IPNetwork, http.Header, error) {
	if ipnet == nil {
		return nil, nil, fmt.Errorf("undefined IP network")
	}

	resp, err := c.Transport.Fetch(c.URIs, QueryTypeIP, ipnet.String(), header, queryString)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if err != nil {
		if resp != nil {
			return nil, resp.Header, err
		}
		return nil, nil, err
	}

	ipNetwork := &protocol.IPNetwork{}
	if err = json.NewDecoder(resp.Body).Decode(ipNetwork); err != nil {
		return nil, resp.Header, err
	}

	return ipNetwork, resp.Header, nil
}

// IP will query each RDAP server to retrieve the desired information and
// will parse and store the response into a protocol IP object. You can
// optionally define the HTTP headers parameters to send to the RDAP server.
// If something goes wrong an error will be returned, and if nothing is found
// the error ErrNotFound will be returned. The HTTP header of the RDAP
// response is also returned to analyze any specific flag
func (c *Client) IP(ip net.IP, header http.Header, queryString url.Values) (*protocol.IPNetwork, http.Header, error) {
	if ip == nil {
		return nil, nil, fmt.Errorf("undefined IP")
	}

	resp, err := c.Transport.Fetch(c.URIs, QueryTypeIP, ip.String(), header, queryString)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if err != nil {
		if resp != nil {
			return nil, resp.Header, err
		}
		return nil, nil, err
	}

	ipNetwork := &protocol.IPNetwork{}
	if err = json.NewDecoder(resp.Body).Decode(ipNetwork); err != nil {
		return nil, resp.Header, err
	}

	return ipNetwork, resp.Header, nil
}

// Query will try to search the object in the following order: ASN, IP, IP
// network, domain and entity. If the format is not valid for the specific
// search, the search is ignored. The HTTP header of the RDAP
// response is also returned to analyze any specific flag
func (c *Client) Query(object string, header http.Header, queryString url.Values) (interface{}, http.Header, error) {
	if asn, err := strconv.ParseUint(object, 10, 32); err == nil {
		return c.ASN(uint32(asn), header, queryString)
	}

	if ip := net.ParseIP(object); ip != nil {
		return c.IP(ip, header, queryString)
	}

	if _, ipnetwork, err := net.ParseCIDR(object); err == nil {
		return c.IPNetwork(ipnetwork, header, queryString)
	}

	fqdn, err := idna.ToASCII(strings.ToLower(object))
	if err == nil && fqdnRX.MatchString(fqdn) {
		return c.Domain(fqdn, header, queryString)
	}

	return c.Entity(object, header, queryString)
}
