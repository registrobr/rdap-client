package rdap

import (
	"fmt"
	"net/http"

	"reflect"
	"testing"

	"github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/registrobr/rdap/protocol"
)

func TestHandlerDomain(t *testing.T) {
	tests := []struct {
		description     string
		identifier      string
		bootstrapEntry  string
		bootstrapStatus int
		expectedObject  interface{}
		expectedError   error
	}{
		{
			description:   "Domain handler should not be executed due to an invalid identifier in input",
			identifier:    "invalid&invalid",
			expectedError: ErrInvalidQuery,
		},
		{
			description: "Domain handler should return a valid RDAP response",
			identifier:  "example.br",
			expectedObject: protocol.Domain{
				ObjectClassName: "domain",
				LDHName:         "example.br",
			},
		},
		{
			description:    "Domain handler should return a valid RDAP response (bootstrapped)",
			identifier:     "example.br",
			bootstrapEntry: "br",
			expectedObject: protocol.Domain{
				ObjectClassName: "domain",
				LDHName:         "example.br",
			},
		},
		{
			description:     "Domain handler should return a HTTP 500 error from bootstrap server",
			identifier:      "example.br",
			bootstrapEntry:  "br",
			bootstrapStatus: http.StatusInternalServerError,
			expectedError:   fmt.Errorf("unexpected status code 500 Internal Server Error"),
		},
	}

	for i, test := range tests {
		ts, bs := createTestServers(
			test.expectedObject,
			test.bootstrapEntry,
			test.bootstrapStatus,
			0,
		)

		h := Handler{
			URIs: []string{ts.URL},
		}

		if len(test.bootstrapEntry) > 0 {
			h.Bootstrap = &Bootstrap{Bootstrap: bs.URL + "/%v"}
		}

		object, err := h.Domain(test.identifier)

		if test.expectedError != nil {
			if fmt.Sprintf("%v", test.expectedError) != fmt.Sprintf("%v", err) {
				t.Fatalf("[%d] “%s“: expected error “%s”, got “%s”", i, test.description, test.expectedError, err)
			}
		} else {
			if object != nil && !reflect.DeepEqual(test.expectedObject, *object) {
				for _, l := range diff(test.expectedObject, *object) {
					t.Log(l)
				}

				t.Fatalf("[%d] “%s”", i, test.description)
			}
		}
	}
}

func TestHandlerASN(t *testing.T) {
	tests := []struct {
		description     string
		identifier      string
		bootstrapEntry  string
		bootstrapStatus int
		expectedObject  interface{}
		expectedError   error
	}{
		{
			description:   "ASN handler should not be executed due to an invalid identifier in input",
			identifier:    "invalid&invalid",
			expectedError: ErrInvalidQuery,
		},
		{
			description: "ASN handler should return a valid RDAP response",
			identifier:  "1",
			expectedObject: protocol.AS{
				ObjectClassName: "as",
				StartAutnum:     1,
				EndAutnum:       16,
			},
		},
		{
			description:    "ASN handler should return a valid RDAP response (bootstrapped)",
			identifier:     "1",
			bootstrapEntry: "1-16",
			expectedObject: protocol.AS{
				ObjectClassName: "as",
				StartAutnum:     1,
				EndAutnum:       16,
			},
		},
		{
			description:     "ASN handler should return a HTTP 500 error from bootstrap server",
			identifier:      "1",
			bootstrapEntry:  "1-16",
			bootstrapStatus: http.StatusInternalServerError,
			expectedError:   fmt.Errorf("unexpected status code 500 Internal Server Error"),
		},
	}

	for i, test := range tests {
		ts, bs := createTestServers(
			test.expectedObject,
			test.bootstrapEntry,
			test.bootstrapStatus,
			0,
		)

		h := Handler{
			URIs: []string{ts.URL},
		}

		if len(test.bootstrapEntry) > 0 {
			h.Bootstrap = &Bootstrap{Bootstrap: bs.URL + "/%v"}
		}

		object, err := h.ASN(test.identifier)

		if test.expectedError != nil {
			if fmt.Sprintf("%v", test.expectedError) != fmt.Sprintf("%v", err) {
				t.Fatalf("[%d] “%s“: expected error “%s”, got “%s”", i, test.description, test.expectedError, err)
			}
		} else {
			if object != nil && !reflect.DeepEqual(test.expectedObject, *object) {
				for _, l := range diff(test.expectedObject, *object) {
					t.Log(l)
				}

				t.Fatalf("[%d] “%s”", i, test.description)
			}
		}
	}
}

func TestHandlerIP(t *testing.T) {
	tests := []struct {
		description     string
		identifier      string
		bootstrapEntry  string
		bootstrapStatus int
		expectedObject  interface{}
		expectedError   error
	}{
		{
			description:   "IP handler should not be executed due to an invalid identifier in input",
			identifier:    "invalid",
			expectedError: ErrInvalidQuery,
		},
		{
			description: "IP handler should return a valid RDAP response",
			identifier:  "192.168.0.1",
			expectedObject: protocol.IPNetwork{
				ObjectClassName: "ipnetwork",
				StartAddress:    "192.168.0.0",
				EndAddress:      "192.168.0.255",
			},
		},
		{
			description:    "IP handler should return a valid RDAP response (bootstrapped)",
			identifier:     "192.168.0.1",
			bootstrapEntry: "192.168.0.0/24",
			expectedObject: protocol.IPNetwork{
				ObjectClassName: "ipnetwork",
				StartAddress:    "192.168.0.0",
				EndAddress:      "192.168.0.255",
			},
		},
		{
			description:     "IP handler should return a HTTP 500 error from bootstrap server",
			identifier:      "192.168.0.1",
			bootstrapEntry:  "192.168.0.0/24",
			bootstrapStatus: http.StatusInternalServerError,
			expectedError:   fmt.Errorf("unexpected status code 500 Internal Server Error"),
		},
	}

	for i, test := range tests {
		ts, bs := createTestServers(
			test.expectedObject,
			test.bootstrapEntry,
			test.bootstrapStatus,
			0,
		)

		h := Handler{
			URIs: []string{ts.URL},
		}

		if len(test.bootstrapEntry) > 0 {
			h.Bootstrap = &Bootstrap{Bootstrap: bs.URL + "/%v"}
		}

		object, err := h.IP(test.identifier)

		if test.expectedError != nil {
			if fmt.Sprintf("%v", test.expectedError) != fmt.Sprintf("%v", err) {
				t.Fatalf("[%d] “%s“: expected error “%s”, got “%s”", i, test.description, test.expectedError, err)
			}
		} else {
			if object != nil && !reflect.DeepEqual(test.expectedObject, *object) {
				for _, l := range diff(test.expectedObject, *object) {
					t.Log(l)
				}

				t.Fatalf("[%d] “%s”", i, test.description)
			}
		}
	}
}

func TestHandlerIPNetwork(t *testing.T) {
	tests := []struct {
		description     string
		identifier      string
		bootstrapEntry  string
		bootstrapStatus int
		expectedObject  interface{}
		expectedError   error
	}{
		{
			description:   "IP handler should not be executed due to an invalid identifier in input",
			identifier:    "invalid",
			expectedError: ErrInvalidQuery,
		},
		{
			description: "IP handler should return a valid RDAP response",
			identifier:  "192.168.0.0/24",
			expectedObject: protocol.IPNetwork{
				ObjectClassName: "ipnetwork",
				StartAddress:    "192.168.0.0",
				EndAddress:      "192.168.0.255",
			},
		},
		{
			description:    "IP handler should return a valid RDAP response (bootstrapped)",
			identifier:     "192.168.0.0/24",
			bootstrapEntry: "192.168.0.0/16",
			expectedObject: protocol.IPNetwork{
				ObjectClassName: "ipnetwork",
				StartAddress:    "192.168.0.0",
				EndAddress:      "192.168.0.255",
			},
		},
		{
			identifier:      "192.168.0.0/24",
			bootstrapEntry:  "192.168.0.0/16",
			bootstrapStatus: http.StatusInternalServerError,
			expectedError:   fmt.Errorf("unexpected status code 500 Internal Server Error"),
		},
	}

	for i, test := range tests {
		ts, bs := createTestServers(
			test.expectedObject,
			test.bootstrapEntry,
			test.bootstrapStatus,
			0,
		)

		h := Handler{
			URIs: []string{ts.URL},
		}

		if len(test.bootstrapEntry) > 0 {
			h.Bootstrap = &Bootstrap{Bootstrap: bs.URL + "/%v"}
		}

		object, err := h.IPNetwork(test.identifier)

		if test.expectedError != nil {
			if fmt.Sprintf("%v", test.expectedError) != fmt.Sprintf("%v", err) {
				t.Fatalf("[%d] “%s“: expected error “%s”, got “%s”", i, test.description, test.expectedError, err)
			}
		} else {
			if object != nil && !reflect.DeepEqual(test.expectedObject, *object) {
				for _, l := range diff(test.expectedObject, *object) {
					t.Log(l)
				}

				t.Fatalf("[%d] “%s”", i, test.description)
			}
		}
	}
}

func TestHandlerEntity(t *testing.T) {
	tests := []struct {
		description     string
		identifier      string
		bootstrapEntry  string
		bootstrapStatus int
		expectedObject  interface{}
	}{
		{
			description: "Entity handler should return a valid RDAP response",
			identifier:  "someone",
			expectedObject: protocol.Entity{
				ObjectClassName: "entity",
				Handle:          "someone",
			},
		},
	}

	for i, test := range tests {
		ts, bs := createTestServers(
			test.expectedObject,
			test.bootstrapEntry,
			test.bootstrapStatus,
			0,
		)

		h := Handler{
			URIs: []string{ts.URL},
		}

		if len(test.bootstrapEntry) > 0 {
			h.Bootstrap = &Bootstrap{Bootstrap: bs.URL + "/%v"}
		}

		object, err := h.Entity(test.identifier)

		if err != nil {
			t.Fatalf("[%d] “%s”: not expecting error “%v”", i, test.description, err)
		}

		if object != nil && !reflect.DeepEqual(test.expectedObject, *object) {
			for _, l := range diff(test.expectedObject, *object) {
				t.Log(l)
			}

			t.Fatalf("[%d] “%s”", i, test.description)
		}
	}
}

func TestHandlerQuery(t *testing.T) {
	tests := []struct {
		description    string
		identifier     string
		rdapStatus     int
		expectedObject interface{}
		expectedError  error
	}{
		{
			description:    "Generic handler should return an object of type protocol.AS",
			identifier:     "1",
			expectedObject: protocol.AS{},
		},
		{
			description:    "Generic handler should return an object of type protocol.IPNetwork",
			identifier:     "192.168.0.1",
			expectedObject: protocol.IPNetwork{},
		},
		{
			description:    "Generic handler should return an object of type protocol.IPNetwork",
			identifier:     "192.168.0.0/24",
			expectedObject: protocol.IPNetwork{},
		},
		{
			description:    "Generic handler should return an object of type protocol.Domain",
			identifier:     "example.br",
			expectedObject: protocol.Domain{},
		},
		{
			description:    "Generic handler should return an object of type protocol.Entity",
			identifier:     "someone",
			expectedObject: protocol.Entity{},
		},
		{
			description:   "Generic handler should return an HTTP 500 error from RDAP server",
			identifier:    "someone",
			rdapStatus:    http.StatusInternalServerError,
			expectedError: fmt.Errorf("unexpected response: 500 Internal Server Error"),
		},
	}

	for i, test := range tests {
		ts, _ := createTestServers(test.expectedObject, "", 0, test.rdapStatus)

		h := Handler{
			URIs: []string{ts.URL},
		}

		object, err := h.Query(test.identifier)

		if test.expectedError != nil {
			if fmt.Sprintf("%v", test.expectedError) != fmt.Sprintf("%v", err) {
				t.Fatalf("[%d] “%s“: expected error “%s”, got “%s”", i, test.description, test.expectedError, err)
			}
		} else {
			expectedObjType := objType(test.expectedObject)
			objType := objType(object)

			if expectedObjType != objType {
				t.Fatalf("[%d] “%s” expected type “%s”, got “%s” ", i, test.description, expectedObjType, objType)
			}
		}
	}
}
