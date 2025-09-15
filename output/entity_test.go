package output

import (
	"testing"
	"time"

	"github.com/registrobr/rdap/protocol"
)

var expectedEntityOutput = `handle:   XXXX
person:   Joe User
e-mail:   joe.user@example.com
address:  Av Naçoes Unidas, 11541, 7 andar, Sao Paulo, SP, 04578-000, BR
phone:    tel:+55-11-5509-3506;ext=3506
created:  20150301
changed:  20150310

`

func TestEntityPrint(t *testing.T) {
	entityResponse := protocol.Entity{
		ObjectClassName: "entity",
		Handle:          "XXXX",
		VCardArray: []any{
			"vcard",
			[]any{
				[]any{"version", struct{}{}, "text", "4.0"},
				[]any{"fn", struct{}{}, "text", "Joe User"},
				[]any{"kind", struct{}{}, "text", "individual"},
				[]any{"email", struct{ Type string }{Type: "work"}, "text", "joe.user@example.com"},
				[]any{"lang", struct{ Pref string }{Pref: "1"}, "language-tag", "pt"},
				[]any{"adr", struct{ Type string }{Type: "work"}, "text",
					[]any{
						"Av Naçoes Unidas", "11541", "7 andar", "Sao Paulo", "SP", "04578-000", "BR",
					},
				},
				[]any{"tel", struct{ Type string }{Type: "work"}, "uri", "tel:+55-11-5509-3506;ext=3506"},
			},
		},
		LegalRepresentative: "Joe User",
		Entities: []protocol.Entity{
			{
				ObjectClassName: "entity",
				Handle:          "XXXX",
				VCardArray: []any{
					"vcard",
					[]any{
						[]any{"version", struct{}{}, "text", "4.0"},
						[]any{"fn", struct{}{}, "text", "Joe User"},
						[]any{"kind", struct{}{}, "text", "individual"},
						[]any{"email", struct{ Type string }{Type: "work"}, "text", "joe.user@example.com"},
						[]any{"lang", struct{ Pref string }{Pref: "1"}, "language-tag", "pt"},
						[]any{"adr", struct{ Type string }{Type: "work"}, "text",
							[]any{
								"Av Naçoes Unidas", "11541", "7 andar", "Sao Paulo", "SP", "04578-000", "BR",
							},
						},
						[]any{"tel", struct{ Type string }{Type: "work"}, "uri", "tel:+55-11-5509-3506;ext=3506"},
					},
				},
				Events: []protocol.Event{
					{
						Action: protocol.EventActionRegistration,
						Date:   protocol.Date(2015, 03, 01, 12, 00, 00, 00, time.UTC),
					},
					{
						Action: protocol.EventActionLastChanged,
						Date:   protocol.Date(2015, 03, 10, 14, 00, 00, 00, time.UTC),
					},
				},
			},
		},
		Events: []protocol.Event{
			{
				Action: protocol.EventActionRegistration,
				Date:   protocol.Date(2015, 03, 01, 12, 00, 00, 00, time.UTC),
			},
			{
				Action: protocol.EventActionLastChanged,
				Date:   protocol.Date(2015, 03, 10, 14, 00, 00, 00, time.UTC),
			},
		},
	}

	entityOutput := Entity{Entity: &entityResponse}

	var w WriterMock
	if err := entityOutput.Print(&w); err != nil {
		t.Fatal(err)
	}

	if string(w.Content) != expectedEntityOutput {
		for _, l := range diff(expectedEntityOutput, string(w.Content)) {
			t.Log(l)
		}
		t.Fatal("error")
	}
}
