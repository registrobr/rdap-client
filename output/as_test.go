package output

import (
	"errors"
	"testing"
	"time"

	"github.com/registrobr/rdap/protocol"
)

func TestASPrint(t *testing.T) {
	as := AS{
		AS: &protocol.AS{
			ObjectClassName: "autnum",
			Handle:          "123456",
			StartAutnum:     123456,
			EndAutnum:       123456,
			Type:            "DIRECT ALLOCATION",
			Country:         "BR",
			Links: []protocol.Link{
				{
					Value: "https://rdap.registro.br/autnum/123456",
					Rel:   "self",
					Href:  "https://rdap.registro.br/autnum/123456",
					Type:  "application/rdap+json",
				},
				{
					Value: "https://rdap.registro.br/autnum/123456",
					Rel:   "related",
					Href:  "https://rdap.registro.br/ip/200.160.0.0/20",
					Type:  "application/rdap+json",
				},
			},
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
				{
					ObjectClassName: "entity",
					Handle:          "YYYY",
					VCardArray: []any{
						"vcard",
						[]any{
							[]any{"version", struct{}{}, "text", "4.0"},
							[]any{"fn", struct{}{}, "text", "Joe User 2"},
							[]any{"kind", struct{}{}, "text", "individual"},
							[]any{"email", struct{ Type string }{Type: "work"}, "text", "joe.user2@example.com"},
							[]any{"lang", struct{ Pref string }{Pref: "1"}, "language-tag", "pt"},
							[]any{"adr", struct{ Type string }{Type: "work"}, "text",
								[]any{
									"Av Naçoes Unidas", "11541", "7 andar", "Sao Paulo", "SP", "04578-000", "BR",
								},
							},
							[]any{"tel", struct{ Type string }{Type: "work"}, "uri", "tel:+55-11-5509-3506;ext=3507"},
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
		},
	}

	expected := `
aut-num:     123456
type:        DIRECT ALLOCATION
country:     BR
created:     20150301
changed:     20150310
inetnum:     200.160.0.0/20

handle:   XXXX
person:   Joe User
e-mail:   joe.user@example.com
address:  Av Naçoes Unidas, 11541, 7 andar, Sao Paulo, SP, 04578-000, BR
phone:    tel:+55-11-5509-3506;ext=3506
created:  20150301
changed:  20150310

handle:   YYYY
person:   Joe User 2
e-mail:   joe.user2@example.com
address:  Av Naçoes Unidas, 11541, 7 andar, Sao Paulo, SP, 04578-000, BR
phone:    tel:+55-11-5509-3506;ext=3507
created:  20150301
changed:  20150310

`

	var w WriterMock
	if err := as.Print(&w); err != nil {
		t.Fatal(err)
	}

	if string(w.Content) != expected {
		for _, l := range diff(expected, string(w.Content)) {
			t.Log(l)
		}
		t.Fatal("error")
	}
}

func TestAsToTextWithErrorOnWriter(t *testing.T) {
	dummyErr := errors.New("Dummy Error!")
	w := &WriterMock{
		Err: dummyErr,
	}

	as := AS{
		AS: new(protocol.AS),
	}

	if err := as.Print(w); err == nil {
		t.Fatal("Expecting an error")
	}
}
