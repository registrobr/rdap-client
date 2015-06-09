package client

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/registrobr/rdap/Godeps/_workspace/src/github.com/miekg/dns/idn"
	"github.com/registrobr/rdap/protocol"
)

const (
	domain kind = "domain"
	autnum kind = "autnum"
	ip     kind = "ip"
	entity kind = "entity"
)

type Client struct {
	httpClient *http.Client
	uris       []string
}

func NewClient(uris []string, httpClient *http.Client) *Client {
	return &Client{
		uris:       uris,
		httpClient: httpClient,
	}
}

func (c *Client) Domain(fqdn string) (*protocol.DomainResponse, error) {
	r := &protocol.DomainResponse{}
	fqdn = idn.ToPunycode(strings.ToLower(fqdn))

	if err := c.query(domain, fqdn, r); err != nil {
		return nil, err
	}

	return r, nil
}

func (c *Client) ASN(as uint64) (*protocol.ASResponse, error) {
	r := &protocol.ASResponse{}

	if err := c.query(autnum, as, r); err != nil {
		return nil, err
	}

	return r, nil
}

func (c *Client) Entity(identifier string) (*protocol.Entity, error) {
	r := &protocol.Entity{}

	if err := c.query(entity, identifier, r); err != nil {
		return nil, err
	}

	return r, nil
}

func (c *Client) IPNetwork(ipnet *net.IPNet) (*protocol.IPNetwork, error) {
	r := &protocol.IPNetwork{}

	if err := c.query(ip, ipnet, r); err != nil {
		return nil, err
	}

	return r, nil
}

func (c *Client) IP(netIP net.IP) (*protocol.IPNetwork, error) {
	r := &protocol.IPNetwork{}

	if err := c.query(ip, netIP, r); err != nil {
		return nil, err
	}

	return r, nil
}

func (c *Client) handleHTTPStatusCode(kind kind, response *http.Response) error {
	if response.StatusCode == http.StatusOK {
		return nil
	}

	if response.StatusCode == http.StatusNotFound {
		return fmt.Errorf("%s not found.", kind)
	}

	if response.Header.Get("Content-Type") != "application/json" {
		return fmt.Errorf("unexpected response: %d %s",
			response.StatusCode, http.StatusText(response.StatusCode))
	}

	var responseErr protocol.Error
	if err := json.NewDecoder(response.Body).Decode(&responseErr); err != nil {
		return err
	}

	return fmt.Errorf("HTTP status code: %d (%s)\n%s:\n  %s",
		responseErr.ErrorCode,
		http.StatusText(responseErr.ErrorCode),
		responseErr.Title,
		strings.Join(responseErr.Description, "\n  "))
}

func (c *Client) query(kind kind, identifier interface{}, object interface{}) (err error) {
	errors := make([]string, 0)
	for _, uri := range c.uris {
		uri := fmt.Sprintf("%s/%s/%v", uri, kind, identifier)

		res, err := c.fetch(uri)
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}
		defer res.Body.Close()

		if err := c.handleHTTPStatusCode(kind, res); err != nil {
			errors = append(errors, err.Error())
			continue
		}

		if err = json.NewDecoder(res.Body).Decode(&object); err != nil {
			errors = append(errors, err.Error())
			continue
		}

		return nil
	}

	return fmt.Errorf("error(s) fetching RDAP data from %v:\n  %s", identifier, strings.Join(errors, "\n  "))
}

func (c *Client) fetch(uri string) (response *http.Response, err error) {
	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	if c.httpClient == nil {
		c.httpClient = &http.Client{}
	}

	response, err = c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	return response, nil
}
