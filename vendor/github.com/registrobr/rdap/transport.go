package rdap

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/registrobr/rdap/protocol"
)

// List of resource type path segments for exact match lookup as described in
// RFC 7482, section 3.1
const (
	// QueryTypeDomain used to identify reverse DNS (RIR) or domain name (DNR)
	// information and associated data referenced using a fully qualified domain
	// name
	QueryTypeDomain QueryType = "domain"

	// QueryTypeTicket used to query a domain request. This query type was
	// created by NIC.br to allow retrieving information about the domain
	// requests
	QueryTypeTicket QueryType = "ticket"

	// QueryTypeAutnum used to identify Autonomous System number registrations
	// and associated data referenced using an asplain Autonomous System number
	QueryTypeAutnum QueryType = "autnum"

	// QueryTypeIP used to identify IP networks and associated data referenced
	// using either an IPv4 or IPv6 address
	QueryTypeIP QueryType = "ip"

	// QueryTypeEntity used to identify an entity information query using a
	// string identifier
	QueryTypeEntity QueryType = "entity"
)

// QueryType stores the query type when sending a query to an RDAP server
type QueryType string

const (
	bootstrapQueryTypeNone bootstrapQueryType = ""
	bootstrapQueryTypeDNS  bootstrapQueryType = "dns"
	bootstrapQueryTypeASN  bootstrapQueryType = "asn"
	bootstrapQueryTypeIPv4 bootstrapQueryType = "ipv4"
	bootstrapQueryTypeIPv6 bootstrapQueryType = "ipv6"
)

type bootstrapQueryType string

func newBootstrapQueryType(queryType QueryType, queryValue string) (bootstrapQueryType, bool) {
	switch queryType {
	case QueryTypeDomain:
		return bootstrapQueryTypeDNS, true

	case QueryTypeAutnum:
		return bootstrapQueryTypeASN, true

	case QueryTypeIP:
		ip := net.ParseIP(queryValue)
		if ip != nil {
			if ip.To4() != nil {
				return bootstrapQueryTypeIPv4, true
			}

			return bootstrapQueryTypeIPv6, true
		}

		var err error
		ip, _, err = net.ParseCIDR(queryValue)
		if err != nil {
			return bootstrapQueryTypeNone, false
		}

		if ip.To4() != nil {
			return bootstrapQueryTypeIPv4, true
		}

		return bootstrapQueryTypeIPv6, true
	}

	return bootstrapQueryTypeNone, false
}

const (
	// IANABootstrap stores the default URL to query to retrieve the RDAP
	// servers that contain the desired information
	IANABootstrap = "https://data.iana.org/rdap/%s.json"
)

var (
	// ErrNotFound is used when the RDAP server doesn't contain any
	// information of the requested object
	ErrNotFound  = errors.New("not found")
	ErrForbidden = errors.New("forbidden")
)

// Fetcher represents the network layer responsible for retrieving the
// resource information from a RDAP server
type Fetcher interface {
	Fetch(uris []string, queryType QueryType, queryValue string, header http.Header, queryString url.Values) (*http.Response, error)
}

// fetcherFunc is a function type that implements the Fetcher interface
type fetcherFunc func([]string, QueryType, string, http.Header, url.Values) (*http.Response, error)

// Fetch will try to use the addresses from the uris parameter to send
// requests using the queryType and queryValue parameters. You can optionally
// set HTTP headers (like X-Forwarded-For) for the RDAP server request. On
// success will return a HTTP response, otherwise an error will be returned.
// The caller is responsible for closing the response body
func (f fetcherFunc) Fetch(uris []string, queryType QueryType, queryValue string, header http.Header, queryString url.Values) (*http.Response, error) {
	return f(uris, queryType, queryValue, header, queryString)
}

type decorator func(Fetcher) Fetcher

func decorate(f Fetcher, ds ...decorator) Fetcher {
	for _, decorate := range ds {
		f = decorate(f)
	}

	return f
}

// CacheDetector is used to define how do you detect if a HTTP response is
// from cache when performing bootstrap. This depends on the proxy that you
// are using between the client and the bootstrap server
type CacheDetector func(*http.Response) bool

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type defaultFetcher struct {
	httpClient httpClient
}

// NewDefaultFetcher returns a transport layer that send requests directly to
// the RDAP servers
func NewDefaultFetcher(httpClient httpClient) Fetcher {
	return &defaultFetcher{
		httpClient: httpClient,
	}
}

func (d *defaultFetcher) Fetch(uris []string, queryType QueryType, queryValue string, header http.Header, queryString url.Values) (resp *http.Response, err error) {
	if len(uris) == 0 {
		return nil, fmt.Errorf("no URIs defined to query")
	}

	for _, uri := range uris {
		resp, err = d.fetchURI(uri, queryType, queryValue, header, queryString)
		if err != nil {
			continue
		}
		return
	}

	return
}

func (d *defaultFetcher) fetchURI(uri string, queryType QueryType, queryValue string, header http.Header, queryString url.Values) (*http.Response, error) {
	if !strings.HasPrefix(uri, "http://") && !strings.HasPrefix(uri, "https://") {
		uri = "http://" + uri
	}

	if pos := strings.Index(uri, "?"); pos != -1 {
		uri = uri[:pos]
	}

	uri = strings.TrimRight(uri, "/")
	uri = fmt.Sprintf("%s/%s/%s", uri, queryType, queryValue)

	if q := queryString.Encode(); len(q) > 0 {
		uri += "?" + q
	}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	if header != nil {
		req.Header = header
	}

	req.Header.Set("Accept", "application/rdap+json")
	req.Header.Set("User-Agent", "registrobr-rdap")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		// we will return the response here so the client can analyze the body or
		// some special HTTP headers to identify the reason why it does not
		// exists
		return resp, ErrNotFound
	} else if resp.StatusCode == http.StatusForbidden {
		return resp, ErrForbidden
	}

	contentType := resp.Header.Get("Content-Type")
	contentTypeParts := strings.Split(contentType, ";")

	if len(contentTypeParts) == 0 || contentTypeParts[0] != "application/rdap+json" {
		return nil, fmt.Errorf("unexpected response: %d %s",
			resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	if resp.StatusCode != http.StatusOK {
		var responseErr protocol.Error
		if err := json.NewDecoder(resp.Body).Decode(&responseErr); err != nil {
			return nil, err
		}

		return nil, responseErr
	}

	return resp, nil
}

// NewBootstrapFetcher returns a transport layer that tries to find the
// resource in a bootstrap strategy to detect the RDAP servers that can contain
// the information. After finding the RDAP servers, it will send the requests to
// retrieve the desired information
func NewBootstrapFetcher(httpClient httpClient, bootstrapURI string, cacheDetector CacheDetector) Fetcher {
	return decorate(
		NewDefaultFetcher(httpClient),
		bootstrap(bootstrapURI, httpClient, cacheDetector),
	)
}

func bootstrap(bootstrapURI string, httpClient httpClient, cacheDetector CacheDetector) decorator {
	return func(f Fetcher) Fetcher {
		return fetcherFunc(func(uris []string, queryType QueryType, queryValue string, header http.Header, queryString url.Values) (*http.Response, error) {
			bootstrapQueryType, ok := newBootstrapQueryType(queryType, queryValue)
			if !ok {
				// if we can't convert the queryType the resource is probably not
				// supported by the bootstrap
				return f.Fetch(uris, queryType, queryValue, header, queryString)
			}
			bootstrapURI := fmt.Sprintf(bootstrapURI, bootstrapQueryType)

			serviceRegistry, cached, err := bootstrapFetch(httpClient, bootstrapURI, false, cacheDetector)
			if err != nil {
				return nil, err
			}

			switch queryType {
			case QueryTypeDomain:
				uris, err = serviceRegistry.matchDomain(queryValue)
				if err == nil && len(uris) == 0 && cached {
					var nsSet []*net.NS
					if nsSet, err = lookupNS(queryValue); err == nil && len(nsSet) > 0 {
						serviceRegistry, cached, err = bootstrapFetch(httpClient, bootstrapURI, true, cacheDetector)
						if err == nil {
							uris, err = serviceRegistry.matchDomain(queryValue)
						}
					}
				}

			case QueryTypeAutnum:
				var asn uint64
				if asn, err = strconv.ParseUint(queryValue, 10, 32); err == nil {
					uris, err = serviceRegistry.matchAS(uint32(asn))
				}

			case QueryTypeIP:
				ip := net.ParseIP(queryValue)
				if ip != nil {
					uris, err = serviceRegistry.matchIP(ip)

				} else {
					var cidr *net.IPNet
					if _, cidr, err = net.ParseCIDR(queryValue); err == nil {
						uris, err = serviceRegistry.matchIPNetwork(cidr)
					}
				}
			}

			if err != nil {
				return nil, err
			}

			if len(uris) == 0 {
				return nil, fmt.Errorf("no matches for %v", queryValue)
			}

			sort.Sort(prioritizeHTTPS(uris))
			return f.Fetch(uris, queryType, queryValue, header, queryString)
		})
	}
}

func bootstrapFetch(httpClient httpClient, uri string, reloadCache bool, cacheDetector CacheDetector) (*serviceRegistry, bool, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, false, err
	}
	req.Header.Add("Accept", "application/json")

	if reloadCache {
		req.Header.Add("Cache-Control", "max-age=0")
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, false, err
	}
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	cached := false
	if cacheDetector != nil {
		cached = cacheDetector(resp)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotModified {
		return nil, cached, fmt.Errorf("unexpected status code %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	var serviceRegistry serviceRegistry
	if err := json.NewDecoder(resp.Body).Decode(&serviceRegistry); err != nil {
		return nil, cached, err
	}

	if serviceRegistry.Version != version {
		return nil, false, fmt.Errorf("incompatible bootstrap specification version: %s (expecting %s)", serviceRegistry.Version, version)
	}

	return &serviceRegistry, cached, nil
}

var lookupNS = func(name string) (nss []*net.NS, err error) {
	return net.LookupNS(name)
}
