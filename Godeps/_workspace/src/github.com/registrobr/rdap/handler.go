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
	isFQDN          = regexp.MustCompile(`^(([[:alnum:]](([[:alnum:]]|\-){0,61}[[:alnum:]])?\.)*[[:alnum:]](([[:alnum:]]|\-){0,61}[[:alnum:]])?)?(\.)?$`)
	ErrInvalidQuery = errors.New("invalid query")
)

type Handler struct {
	URIs          []string
	HTTPClient    *http.Client
	Bootstrap     *Bootstrap
	XForwardedFor string
}

func (h *Handler) Query(object string) (interface{}, error) {
	generic := genericQuerier{handler: h}
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

func (h *Handler) ASN(object string) (*protocol.AS, error) {
	asn, err := strconv.ParseUint(object, 10, 32)

	if err != nil {
		return nil, ErrInvalidQuery
	}

	uris := h.URIs

	if h.Bootstrap != nil {
		var err error
		uris, err = h.Bootstrap.ASN(asn)

		if err != nil {
			return nil, err
		}
	}

	cli := NewClient(uris, h.HTTPClient)
	if h.XForwardedFor != "" {
		cli.XForwardedFor = h.XForwardedFor
	}

	return cli.ASN(asn)
}

func (h *Handler) Entity(object string) (*protocol.Entity, error) {
	// Note that there is no bootstrap for entity, see [1]
	// [1] - https://tools.ietf.org/html/rfc7484#section-6

	cli := NewClient(h.URIs, h.HTTPClient)
	if h.XForwardedFor != "" {
		cli.XForwardedFor = h.XForwardedFor
	}

	return cli.Entity(object)
}

func (h *Handler) IPNetwork(object string) (*protocol.IPNetwork, error) {
	_, cidr, err := net.ParseCIDR(object)

	if err != nil {
		return nil, ErrInvalidQuery
	}

	uris := h.URIs

	if h.Bootstrap != nil {
		var err error
		uris, err = h.Bootstrap.IPNetwork(cidr)

		if err != nil {
			return nil, err
		}
	}

	cli := NewClient(uris, h.HTTPClient)
	if h.XForwardedFor != "" {
		cli.XForwardedFor = h.XForwardedFor
	}

	return cli.IPNetwork(cidr)
}

func (h *Handler) IP(object string) (*protocol.IPNetwork, error) {
	ip := net.ParseIP(object)

	if ip == nil {
		return nil, ErrInvalidQuery
	}

	uris := h.URIs

	if h.Bootstrap != nil {
		var err error
		uris, err = h.Bootstrap.IP(ip)

		if err != nil {
			return nil, err
		}
	}

	cli := NewClient(uris, h.HTTPClient)
	if h.XForwardedFor != "" {
		cli.XForwardedFor = h.XForwardedFor
	}

	return cli.IP(ip)
}

func (h *Handler) Domain(object string) (*protocol.Domain, error) {
	if !isFQDN.MatchString(object) || !strings.Contains(object, ".") {
		return nil, ErrInvalidQuery
	}

	uris := h.URIs

	if h.Bootstrap != nil {
		var err error
		uris, err = h.Bootstrap.Domain(object)

		if err != nil {
			return nil, err
		}
	}

	cli := NewClient(uris, h.HTTPClient)
	if h.XForwardedFor != "" {
		cli.XForwardedFor = h.XForwardedFor
	}

	return cli.Domain(object)
}

type genericQuerier struct {
	handler *Handler
}

func (h *genericQuerier) ASN(object string) (interface{}, error) {
	return h.handler.ASN(object)
}

func (h *genericQuerier) Entity(object string) (interface{}, error) {
	return h.handler.Entity(object)
}

func (h *genericQuerier) IPNetwork(object string) (interface{}, error) {
	return h.handler.IPNetwork(object)
}

func (h *genericQuerier) IP(object string) (interface{}, error) {
	return h.handler.IP(object)
}

func (h *genericQuerier) Domain(object string) (interface{}, error) {
	return h.handler.Domain(object)
}
