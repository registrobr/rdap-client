package rdap

import (
	"errors"
	"net"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/registrobr/rdap/protocol"
)

var (
	isFQDN = regexp.MustCompile(`^((([a-z0-9][a-z0-9\-]*[a-z0-9])|[a-z0-9]+)\.)*([a-z]+|xn\-\-[a-z0-9]+)\.?$`)

	// ErrInvalidQuery reports an query with an invalid format
	ErrInvalidQuery = errors.New("invalid query")
)

// BootstrapClient is a client with the bootstrap logic
type BootstrapClient interface {
	Query(object string) (interface{}, error)
	Domain(object string) (*protocol.Domain, error)
	ASN(object string) (*protocol.AS, error)
	Entity(object string) (*protocol.Entity, error)
	IPNetwork(object string) (*protocol.IPNetwork, error)
	IP(object string) (*protocol.IPNetwork, error)
	SetXForwardedFor(addr string)
	SetCacheDetector(isFromCache func(*http.Response) bool)
}

// NewBootstrapClient returns a client with bootstrap capabilities. You can
// inject a HTTP client and a bootstrap URL starting point, if no HTTP client is
// defined, a new one is created, and if no bootstrap URL is given, the
// RDAPBootstrap constant will be used
var NewBootstrapClient = func(uris []string, httpClient *http.Client, bootstrap string) BootstrapClient {
	b := bootstrapClient{
		uris:      uris,
		client:    NewClient(uris, httpClient),
		bootstrap: NewBootstrap(httpClient),
	}

	if bootstrap != "" {
		b.bootstrap.Bootstrap = bootstrap
	}
	return &b
}

type bootstrapClient struct {
	uris      []string
	client    Client
	bootstrap *Bootstrap
}

func (b *bootstrapClient) Query(object string) (interface{}, error) {
	generic := genericQuerier{handler: b}
	handlers := []func(object string) (interface{}, error){
		generic.ASN,
		generic.IP,
		generic.IPNetwork,
		generic.Domain,
		generic.Entity,
	}

	for _, handler := range handlers {
		resp, err := handler(object)

		if err != nil && err != ErrInvalidQuery {
			return nil, err
		}

		// interface{} holding nil value...
		if !reflect.ValueOf(resp).IsNil() {
			return resp, nil
		}
	}

	return nil, ErrInvalidQuery
}

func (b *bootstrapClient) Domain(object string) (*protocol.Domain, error) {
	if !isFQDN.MatchString(object) || !strings.Contains(object, ".") {
		return nil, ErrInvalidQuery
	}

	uris := b.uris

	if b.bootstrap != nil {
		var err error
		uris, err = b.bootstrap.Domain(object)

		if err != nil {
			return nil, err
		}
	}

	b.client.SetURIs(uris)
	return b.client.Domain(object)
}

func (b *bootstrapClient) ASN(object string) (*protocol.AS, error) {
	asn, err := strconv.ParseUint(object, 10, 32)

	if err != nil {
		return nil, ErrInvalidQuery
	}

	uris := b.uris

	if b.bootstrap != nil {
		var err error
		uris, err = b.bootstrap.ASN(asn)

		if err != nil {
			return nil, err
		}
	}

	b.client.SetURIs(uris)
	return b.client.ASN(asn)
}

func (b *bootstrapClient) Entity(object string) (*protocol.Entity, error) {
	// Note that there is no bootstrap for entity, see [1]
	// [1] - https://tools.ietf.org/html/rfc7484#section-6
	b.client.SetURIs(b.uris)
	return b.client.Entity(object)
}

func (b *bootstrapClient) IPNetwork(object string) (*protocol.IPNetwork, error) {
	_, cidr, err := net.ParseCIDR(object)

	if err != nil {
		return nil, ErrInvalidQuery
	}

	uris := b.uris

	if b.bootstrap != nil {
		var err error
		uris, err = b.bootstrap.IPNetwork(cidr)

		if err != nil {
			return nil, err
		}
	}

	b.client.SetURIs(uris)
	return b.client.IPNetwork(cidr)
}

func (b *bootstrapClient) IP(object string) (*protocol.IPNetwork, error) {
	ip := net.ParseIP(object)

	if ip == nil {
		return nil, ErrInvalidQuery
	}

	uris := b.uris

	if b.bootstrap != nil {
		var err error
		uris, err = b.bootstrap.IP(ip)

		if err != nil {
			return nil, err
		}
	}

	b.client.SetURIs(uris)
	return b.client.IP(ip)
}

func (b *bootstrapClient) SetXForwardedFor(addr string) {
	b.client.SetXForwardedFor(addr)
}

func (b *bootstrapClient) SetCacheDetector(isFromCache func(*http.Response) bool) {
	b.bootstrap.IsFromCache = isFromCache
}

type genericQuerier struct {
	handler BootstrapClient
}

func (g *genericQuerier) ASN(object string) (interface{}, error) {
	return g.handler.ASN(object)
}

func (g *genericQuerier) Entity(object string) (interface{}, error) {
	return g.handler.Entity(object)
}

func (g *genericQuerier) IPNetwork(object string) (interface{}, error) {
	return g.handler.IPNetwork(object)
}

func (g *genericQuerier) IP(object string) (interface{}, error) {
	return g.handler.IP(object)
}

func (g *genericQuerier) Domain(object string) (interface{}, error) {
	return g.handler.Domain(object)
}
