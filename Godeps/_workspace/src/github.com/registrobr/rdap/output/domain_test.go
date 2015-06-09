package output

import (
	"testing"
	"time"

	"github.com/registrobr/rdap/protocol"
)

var TestDomainToTextOutput = `domain:   example.br
nserver:  a.dns.br aa
dsrecord: 12345 RSASHA1 0123456789ABCDEF0123456789ABCDEF01234567
dsstatus: 20150301 ok
created:  20150301
changed:  20150310
status:   active

handle:   XXXX
ids:      
roles:    
person:   Joe User
e-mail:   joe.user@example.com
address:  Av Naçoes Unidas, 11541, 7 andar, Sao Paulo, SP, 04578-000, BR
phone:    tel:+55-11-5509-3506;ext=3506
created:  20150301
changed:  20150310

handle:   YYYY
ids:      
roles:    
person:   Joe User 2
e-mail:   joe.user2@example.com
address:  Av Naçoes Unidas, 11541, 7 andar, Sao Paulo, SP, 04578-000, BR
phone:    tel:+55-11-5509-3506;ext=3507
created:  20150301
changed:  20150310

`

func TestDomainToText(t *testing.T) {
	domainResponse := protocol.DomainResponse{
		ObjectClassName: "domain",
		LDHName:         "example.br",
		Status: []protocol.Status{
			"active",
		},
		Links: []protocol.Link{
			{
				Value: "https://rdap.registro.br/domain/example.br",
				Rel:   "self",
				Href:  "https://rdap.registro.br/domain/example.br",
				Type:  "application/rdap+json",
			},
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
			{
				ObjectClassName: "entity",
				Handle:          "YYYY",
				VCardArray: []interface{}{
					"vcard",
					[]interface{}{
						[]interface{}{"version", struct{}{}, "text", "4.0"},
						[]interface{}{"fn", struct{}{}, "text", "Joe User 2"},
						[]interface{}{"kind", struct{}{}, "text", "individual"},
						[]interface{}{"email", struct{ Type string }{Type: "work"}, "text", "joe.user2@example.com"},
						[]interface{}{"lang", struct{ Pref string }{Pref: "1"}, "language-tag", "pt"},
						[]interface{}{"adr", struct{ Type string }{Type: "work"}, "text",
							[]interface{}{
								"Av Naçoes Unidas", "11541", "7 andar", "Sao Paulo", "SP", "04578-000", "BR",
							},
						},
						[]interface{}{"tel", struct{ Type string }{Type: "work"}, "uri", "tel:+55-11-5509-3506;ext=3507"},
					},
				},
				Events: []protocol.Event{
					protocol.Event{Action: protocol.EventActionRegistration, Actor: "", Date: time.Date(2015, 03, 01, 12, 00, 00, 00, time.UTC)},
					protocol.Event{Action: protocol.EventActionLastChanged, Actor: "", Date: time.Date(2015, 03, 10, 14, 00, 00, 00, time.UTC)},
				},
			},
		},
		Events: []protocol.Event{
			protocol.Event{Action: protocol.EventActionRegistration, Actor: "", Date: time.Date(2015, 03, 01, 12, 00, 00, 00, time.UTC)},
			protocol.Event{Action: protocol.EventActionLastChanged, Actor: "", Date: time.Date(2015, 03, 10, 14, 00, 00, 00, time.UTC)},
		},
		Nameservers: []protocol.Nameserver{
			{
				ObjectClassName: "nameserver",
				LDHName:         "a.dns.br",
				HostStatus:      "aa",
			},
		},
		SecureDNS: protocol.SecureDNS{
			DSData: []protocol.DS{
				{
					KeyTag:    12345,
					Digest:    "0123456789ABCDEF0123456789ABCDEF01234567",
					Algorithm: 5,
					DSStatus:  "ok",
					Events: []protocol.Event{
						protocol.Event{Action: protocol.EventActionRegistration, Actor: "", Date: time.Date(2015, 03, 01, 12, 00, 00, 00, time.UTC)},
					},
				},
			},
		},
	}

	domainOutput := Domain{Domain: &domainResponse}

	var w WriterMock
	if err := domainOutput.ToText(&w); err != nil {
		t.Fatal(err)
	}

	if string(w.Content) != TestDomainToTextOutput {
		for _, l := range diff(TestDomainToTextOutput, string(w.Content)) {
			t.Log(l)
		}
		t.Fatal()
	}
}
