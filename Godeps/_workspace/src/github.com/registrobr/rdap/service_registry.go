package client

import (
	"math"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/registrobr/rdap/Godeps/_workspace/src/github.com/miekg/dns/idn"
)

const version = "1.0"

// ServiceRegistry reflects the structure of a RDAP Bootstrap Service
// Registry.
//
// See http://tools.ietf.org/html/rfc7484#section-10.2
type serviceRegistry struct {
	Version     string    `json:"version"`
	Publication time.Time `json:"publication"`
	Description string    `json:"description,omitempty"`
	Services    []service `json:"services"`
}

// service is an array composed by two items. The first one is a list of
// entries and the second one is a list of URIs.
type service [2][]string

// entries is a helper that returns the list of entries of a service
func (s service) entries() []string {
	return s[0]
}

// uris is a helper that returns the list of URIs of a service
func (s service) uris() []string {
	return s[1]
}

// MatchAS iterates through a list of services looking for the more
// specific range to which an AS number "asn" belongs.
//
// See http://tools.ietf.org/html/rfc7484#section-5.3
func (s serviceRegistry) matchAS(asn uint64) (uris []string, err error) {
	var size uint64 = math.MaxUint32

	for _, service := range s.Services {
		for _, entry := range service.entries() {
			asRange := strings.Split(entry, "-")
			begin, err := strconv.ParseUint(asRange[0], 10, 32)

			if err != nil {
				return nil, err
			}

			end, err := strconv.ParseUint(asRange[1], 10, 32)

			if err != nil {
				return nil, err
			}

			if diff := end - begin; asn >= begin && asn <= end && diff < size {
				size = diff
				uris = service.uris()
			}
		}
	}

	return uris, nil
}

// MatchIPNetwork iterates through a list of services looking for the more
// specific IP network to which the IP network "network" belongs.
//
// See http://tools.ietf.org/html/rfc7484#section-5.1
//     http://tools.ietf.org/html/rfc7484#section-5.2
func (s serviceRegistry) matchIPNetwork(network *net.IPNet) (uris []string, err error) {
	size := 0

	for _, service := range s.Services {
		for _, entry := range service.entries() {
			_, ipnet, err := net.ParseCIDR(entry)

			if err != nil {
				return nil, err
			}

			if ipnet.Contains(network.IP) {
				lastIP := make(net.IP, len(network.IP))

				for i := 0; i < len(network.IP); i++ {
					lastIP[i] = network.IP[i] | ^network.Mask[i]
				}

				if mask, _ := ipnet.Mask.Size(); ipnet.Contains(lastIP) && mask > size {
					uris = service.uris()
					size = mask
				}
			}
		}
	}

	return uris, nil
}

// MatchIP iterates through a list of services looking for the more
// specific IP network to which the IP belongs.
//
// See http://tools.ietf.org/html/rfc7484#section-5.1
//     http://tools.ietf.org/html/rfc7484#section-5.2
func (s serviceRegistry) matchIP(ip net.IP) (uris []string, err error) {
	size := 0

	for _, service := range s.Services {
		for _, entry := range service.entries() {
			_, ipnet, err := net.ParseCIDR(entry)

			if err != nil {
				return nil, err
			}

			if mask, _ := ipnet.Mask.Size(); ipnet.Contains(ip) && mask > size {
				uris = service.uris()
				size = mask
			}
		}
	}

	return uris, nil
}

// MatchDomain iterates through a list of services looking for the label-wise
// longest match of the target domain name "fqdn".
//
// See http://tools.ietf.org/html/rfc7484#section-4
func (s serviceRegistry) matchDomain(fqdn string) (uris []string, err error) {
	var (
		size      int
		fqdnParts = strings.Split(idn.ToPunycode(fqdn), ".")
	)

	for _, service := range s.Services {
	Entries:
		for _, entry := range service.entries() {
			entryParts := strings.Split(entry, ".")

			if len(fqdnParts) < len(entryParts) {
				continue
			}

			fqdnExcerpt := fqdnParts[len(fqdnParts)-len(entryParts):]

			for i := len(entryParts) - 1; i >= 0; i-- {
				if fqdnExcerpt[i] != entryParts[i] {
					continue Entries
				}
			}

			if longest := len(entryParts); longest > size {
				uris = service.uris()
				size = longest
			}
		}
	}

	return uris, nil
}

type prioritizeHTTPS []string

func (v prioritizeHTTPS) Len() int {
	return len(v)
}

func (v prioritizeHTTPS) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v prioritizeHTTPS) Less(i, j int) bool {
	return strings.Split(v[i], ":")[0] == "https"
}
