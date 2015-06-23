package rdap

import (
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"sort"

	"testing"
)

var jsonExample = []byte(`{
       "version": "1.0",
       "publication": "2015-04-17T16:00:00Z",
       "description": "Some text",
       "services": [
         [
           ["entry1", "entry2", "entry3"],
           [
             "https://registry.example.com/myrdap/",
             "http://registry.example.com/myrdap/"
           ]
         ],
         [
           ["entry4"],
           [
             "http://example.org/"
           ]
         ]
       ]
   }`)

func TestServiceRegistryConformity(t *testing.T) {
	if err := json.Unmarshal(jsonExample, &serviceRegistry{}); err != nil {
		t.Fatal(err)
	}
}

func TestServiceRegistryMatchAS(t *testing.T) {
	tests := []struct {
		description   string
		registry      serviceRegistry
		as            uint64
		expected      []string
		expectedError error
	}{
		{
			description: "it should match an as number",
			as:          65411,
			registry: serviceRegistry{
				Services: []service{
					{
						{"2045-2045"},
						{"https://rir3.example.com/myrdap/"},
					},
					{
						{"10000-12000", "300000-400000"},
						{"http://example.org/"},
					},
					{
						{"64512-65534"},
						{"http://example.net/rdaprir2/", "https://example.net/rdaprir2/"},
					},
				},
			},
			expected: []string{
				"http://example.net/rdaprir2/",
				"https://example.net/rdaprir2/",
			},
		},
		{
			description: "it should not match an as number due to invalid beginning of as range",
			as:          1,
			registry: serviceRegistry{
				Services: []service{
					{
						{"invalid-123"},
						{},
					},
				},
			},
			expectedError: fmt.Errorf(`strconv.ParseUint: parsing "invalid": invalid syntax`),
		},
		{
			description: "it should not match an as number due to invalid end of as range",
			as:          1,
			registry: serviceRegistry{
				Services: []service{
					{
						{"123-invalid"},
						{},
					},
				},
			},
			expectedError: fmt.Errorf(`strconv.ParseUint: parsing "invalid": invalid syntax`),
		},
	}

	for i, test := range tests {
		urls, err := test.registry.matchAS(test.as)

		if fmt.Sprintf("%v", test.expectedError) != fmt.Sprintf("%v", err) {
			t.Fatalf("[%d] “%s“: expected error “%s“, got “%s“", i, test.description, test.expectedError, err)
		}

		if !reflect.DeepEqual(test.expected, urls) {
			t.Fatalf("[%d] “%s“: expected “%v“, got “%v“", i, test.description, test.expected, urls)
		}
	}
}

func TestServiceRegistryMatchIPNetwork(t *testing.T) {
	tests := []struct {
		description   string
		registry      serviceRegistry
		ipnet         string
		expected      []string
		expectedError error
	}{
		{
			description: "it should match an ipv6 network",
			ipnet:       "2001:0200:1000::/48",
			registry: serviceRegistry{
				Services: []service{
					{
						{"2001:0200::/23", "2001:db8::/32"},
						{"https://rir2.example.com/myrdap/"},
					},
					{
						{"2600::/16", "2100:ffff::/32"},
						{"http://example.org/"},
					},
					{
						{"2001:0200:1000::/36"},
						{"https://example.net/rdaprir2/", "http://example.net/rdaprir2/"},
					},
				},
			},
			expected: []string{
				"https://example.net/rdaprir2/",
				"http://example.net/rdaprir2/",
			},
		},
		{
			description: "it should match an ipv4 network",
			ipnet:       "192.0.2.1/25",
			registry: serviceRegistry{
				Services: []service{
					{
						{"1.0.0.0/8", "192.0.0.0/8"},
						{"https://rir1.example.com/myrdap/"},
					},
					{
						{"28.2.0.0/16", "192.0.2.0/24"},
						{"http://example.org/"},
					},
					{
						{"28.3.0.0/16"},
						{"https://example.net/rdaprir2/", "http://example.net/rdaprir2/"},
					},
				},
			},
			expected: []string{
				"http://example.org/",
			},
		},
		{
			description: "it should not match an ip network due to invalid cidr",
			ipnet:       "127.0.0.1/32",
			registry: serviceRegistry{
				Services: []service{
					{
						{"invalid"},
						{},
					},
				},
			},
			expectedError: fmt.Errorf(`invalid CIDR address: invalid`),
		},
	}

	for i, test := range tests {
		_, ipnet, _ := net.ParseCIDR(test.ipnet)
		urls, err := test.registry.matchIPNetwork(ipnet)

		if fmt.Sprintf("%v", test.expectedError) != fmt.Sprintf("%v", err) {
			t.Fatalf("[%d] “%s“: expected error “%s“, got “%s“", i, test.description, test.expectedError, err)
		}

		if !reflect.DeepEqual(test.expected, urls) {
			t.Fatalf("[%d] “%s“: expected “%v“, got “%v“", i, test.description, test.expected, urls)
		}
	}
}

func TestServiceRegistryMatchDomain(t *testing.T) {
	tests := []struct {
		description   string
		registry      serviceRegistry
		fqdn          string
		expected      []string
		expectedError error
	}{
		{
			description: "it should match a fqdn",
			fqdn:        "a.b.example.com",
			registry: serviceRegistry{
				Services: []service{
					{
						{"net", "com"},
						{"https://registry.example.com/myrdap/"},
					},
					{
						{"org", "mytld"},
						{"http://example.org/"},
					},
					{
						{"xn--zckzah"},
						{"https://example.net/rdapxn--zckzah/", "http://example.net/rdapxn--zckzah/"},
					},
				},
			},
			expected: []string{
				"https://registry.example.com/myrdap/",
			},
		},
		{
			description: "it should match an idn",
			fqdn:        "feijão.jabá.com",
			registry: serviceRegistry{
				Services: []service{
					{
						{"xn--jab-gla.com"},
						{"https://example.com/myrdap/"},
					},
				},
			},
			expected: []string{"https://example.com/myrdap/"},
		},
		{
			description: "it should match no fqdn",
			fqdn:        "a.example.com",
			registry: serviceRegistry{
				Services: []service{
					{
						{"a.b.example.com"},
						{"https://registry.example.com/myrdap/"},
					},
				},
			},
			expected: nil,
		},
	}

	for i, test := range tests {
		urls, err := test.registry.matchDomain(test.fqdn)

		if fmt.Sprintf("%v", test.expectedError) != fmt.Sprintf("%v", err) {
			t.Fatalf("[%d] “%s“: expected error “%s“, got “%s“", i, test.description, test.expectedError, err)
		}

		if !reflect.DeepEqual(test.expected, urls) {
			t.Fatalf("[%d] “%s“: expected “%v“, got “%v“", i, test.description, test.expected, urls)
		}
	}
}

func TestServiceRegistryMatchIP(t *testing.T) {
	tests := []struct {
		description   string
		registry      serviceRegistry
		ip            string
		expected      []string
		expectedError error
	}{
		{
			description: "it should match an ipv4",
			ip:          "192.0.2.1",
			registry: serviceRegistry{
				Services: []service{
					{
						{"1.0.0.0/8", "192.0.0.0/8"},
						{"https://rir1.example.com/myrdap/"},
					},
					{
						{"28.2.0.0/16", "192.0.2.0/24"},
						{"http://example.org/"},
					},
					{
						{"28.3.0.0/16"},
						{"https://example.net/rdaprir2/", "http://example.net/rdaprir2/"},
					},
				},
			},
			expected: []string{
				"http://example.org/",
			},
		},
		{
			description: "it should not match an ipv4 due to invalid cidr",
			ip:          "127.0.0.1/32",
			registry: serviceRegistry{
				Services: []service{
					{
						{"invalid"},
						{},
					},
				},
			},
			expectedError: fmt.Errorf(`invalid CIDR address: invalid`),
		},
	}

	for i, test := range tests {
		ip := net.ParseIP(test.ip)
		urls, err := test.registry.matchIP(ip)

		if fmt.Sprintf("%v", test.expectedError) != fmt.Sprintf("%v", err) {
			t.Fatalf("[%d] “%s“: expected error “%s“, got “%s“", i, test.description, test.expectedError, err)
		}

		if !reflect.DeepEqual(test.expected, urls) {
			t.Fatalf("[%d] “%s“: expected “%v“, got “%v“", i, test.description, test.expected, urls)
		}
	}
}

func TestPrioritizeHTTPS(t *testing.T) {
	var (
		v  = prioritizeHTTPS{"http:", "https:"}
		v0 = make(prioritizeHTTPS, len(v))
	)

	copy(v0, v)
	sort.Sort(v0)

	if reflect.DeepEqual(v, v0) {
		t.Fatal("not sorting prioritizeHTTPS accordingly")
	}
}
