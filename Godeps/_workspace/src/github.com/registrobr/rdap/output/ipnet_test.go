package output

import (
	"testing"
	"time"

	"github.com/registrobr/rdap/protocol"
)

var TestIPNetworkToTextOutput = `inetnum:       (IPNetwork)
aut-num:       ip_123456-NICBR
abuse-c:       (handle)
owner:         
ownerid:       (CPF/CNPJ)
responsible:   
address:     
address:     
country:       BR
phone:       
start-address: 200.160.3.0
end-address:   200.160.3.255
ip-version:    v4
type:          DIRECT ALLOCATION
parent-handle: ip_1-NICBR
status:        [active]
owner-c:     
tech-c:      
inetrev:     
nserver:     
nsstat:      
nslastaa:    
created:     20150301
changed:     20150310

handle:   XXXX
ids:      
roles:    
person:   Joe User
e-mail:   joe.user@example.com
address:  Av Naçoes Unidas, 11541, 7 andar, Sao Paulo, SP, 04578-000, BR
phone:    tel:+55-11-5509-3506;ext=3506
created:  20150301
changed:  20150310

`

func TestIPNetworkToText(t *testing.T) {
	ipNetworkResponse := protocol.IPNetwork{
		ObjectClassName: "ip network",
		Handle:          "ip_123456-NICBR",
		StartAddress:    "200.160.3.0",
		EndAddress:      "200.160.3.255",
		IPVersion:       "v4",
		Name:            "",
		Type:            "DIRECT ALLOCATION",
		Country:         "BR",
		ParentHandle:    "ip_1-NICBR",
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
	}

	ipNetworkOutput := IPNetwork{IPNetwork: &ipNetworkResponse}

	var w WriterMock
	if err := ipNetworkOutput.ToText(&w); err != nil {
		t.Fatal(err)
	}

	if string(w.Content) != TestIPNetworkToTextOutput {
		for _, l := range diff(TestIPNetworkToTextOutput, string(w.Content)) {
			t.Log(l)
		}
		t.Fatal()
	}
}
