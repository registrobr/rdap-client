package output

import (
	"testing"
	"time"

	"github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/registrobr/rdap/protocol"
)

func TestIPNetPrint(t *testing.T) {
	ipNetwork := IPNetwork{
		IPNetwork: &protocol.IPNetwork{
			ObjectClassName: "ip network",
			Handle:          "200.160.3.0/24",
			ParentHandle:    "200.160.0.0/16",
			StartAddress:    "200.160.3.0",
			EndAddress:      "200.160.3.255",
			IPVersion:       "v4",
			Name:            "Crazy Organization",
			Type:            "DIRECT ALLOCATION",
			Country:         "BR",
			Autnum:          1234,
			Status:          []string{"active"},
			Events: []protocol.Event{
				{Action: protocol.EventActionRegistration, Actor: "", Date: time.Date(2015, 03, 01, 12, 00, 00, 00, time.UTC)},
				{Action: protocol.EventActionLastChanged, Actor: "", Date: time.Date(2015, 03, 10, 14, 00, 00, 00, time.UTC)},
			},
			Entities: []protocol.Entity{
				{
					ObjectClassName: "entity",
					Handle:          "XXXX",
					VCardArray: []interface{}{
						"vcard",
						[]interface{}{
							[]interface{}{"version", struct{}{}, "text", "4.0"},
							[]interface{}{"fn", struct{}{}, "text", "Joe User"},
							[]interface{}{"kind", struct{}{}, "text", "individual"},
							[]interface{}{"email", struct{ Type string }{Type: "work"}, "text", "joe.user@example.com"},
							[]interface{}{"lang", struct{ Pref string }{Pref: "1"}, "language-tag", "pt"},
							[]interface{}{"adr", struct{ Type string }{Type: "work"}, "text",
								[]interface{}{
									"Av Naçoes Unidas", "11541", "7 andar", "Sao Paulo", "SP", "04578-000", "BR",
								},
							},
							[]interface{}{"tel", struct{ Type string }{Type: "work"}, "uri", "tel:+55-11-5509-3506;ext=3506"},
						},
					},
					Events: []protocol.Event{
						protocol.Event{Action: protocol.EventActionRegistration, Actor: "", Date: time.Date(2015, 03, 01, 12, 00, 00, 00, time.UTC)},
						protocol.Event{Action: protocol.EventActionLastChanged, Actor: "", Date: time.Date(2015, 03, 10, 14, 00, 00, 00, time.UTC)},
					},
				},
			},
			ReverseDelegations: []protocol.ReverseDelegation{
				{
					StartAddress: "200.160.3.0",
					EndAddress:   "200.160.3.255",
					Nameservers: []protocol.Nameserver{
						{LDHName: "a.dns.br"},
						{LDHName: "b.dns.br"},
					},
				},
			},
		},
	}

	expected := `
inetnum:       200.160.3.0/24
handle:        200.160.3.0/24
parent-handle: 200.160.0.0/16
aut-num:       1234
start-address: 200.160.3.0
end-address:   200.160.3.255
ip-version:    v4
name:          Crazy Organization
type:          DIRECT ALLOCATION
country:       BR
status:        active
inetrev:       200.160.3.0/24
nserver:       a.dns.br
nserver:       b.dns.br
created:       20150301
changed:       20150310

handle:   XXXX
person:   Joe User
e-mail:   joe.user@example.com
address:  Av Naçoes Unidas, 11541, 7 andar, Sao Paulo, SP, 04578-000, BR
phone:    tel:+55-11-5509-3506;ext=3506
created:  20150301
changed:  20150310

`

	var w WriterMock
	if err := ipNetwork.Print(&w); err != nil {
		t.Fatal(err)
	}

	if string(w.Content) != expected {
		for _, l := range diff(expected, string(w.Content)) {
			t.Log(l)
		}
		t.Fatal("error")
	}
}
