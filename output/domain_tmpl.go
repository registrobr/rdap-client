package output

import (
	"text/template"
	"time"

	"github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/registrobr/rdap/protocol"
)

const (
	domainTmpl = `domain:   {{.Domain.LDHName}}
{{range .Domain.Nameservers}}\
nserver:  {{.LDHName}}
nsstat:   {{nsLastCheck .Events | formatDate}} {{nsStatus .Events}}
nslastaa: {{nsLastOK .Events | formatDate}}
{{end}}\
{{range .DS}}\
dsrecord: {{.KeyTag}} {{.Algorithm | dsAlgorithm}} {{.Digest}}
dsstatus: {{dsLastCheck .Events | formatDate}} {{dsStatus .Events}}
dslastok: {{dsLastOK .Events | formatDate}}
{{end}}\
created:  {{.CreatedAt | formatDate}}
changed:  {{.UpdatedAt | formatDate}}
{{range .Domain.Status}}\
status:   {{.}}
{{end}}\

` + contactTmpl
	dateFormat = "20060102"
)

var (
	dsAlgorithms = map[int]string{
		1:   "RSAMD5",
		2:   "DH",
		3:   "DSASHA1",
		4:   "ECC",
		5:   "RSASHA1",
		6:   "DSASHA1NSEC3",
		7:   "RSASHA1NSEC3",
		8:   "RSASHA256",
		10:  "RSASHA512",
		12:  "ECCGOST",
		13:  "ECDSASHA256",
		14:  "ECDSASHA384",
		252: "INDIRECT",
		253: "PRIVATEDNS",
		254: "PRIVATEOID",
	}

	domainFuncMap = template.FuncMap{
		"nsStatus": func(events []protocol.Event) protocol.Status {
			for _, event := range events {
				if event.Action == protocol.EventDelegationCheck && len(event.Status) > 0 {
					return event.Status[0]
				}
			}

			return protocol.Status("")
		},
		"nsLastCheck": func(events []protocol.Event) time.Time {
			for _, event := range events {
				if event.Action == protocol.EventDelegationCheck && len(event.Status) > 0 {
					return event.Date
				}
			}

			return time.Time{}
		},
		"nsLastOK": func(events []protocol.Event) time.Time {
			for _, event := range events {
				if event.Action == protocol.EventLastCorrectDelegationCheck {
					return event.Date
				}
			}

			return time.Time{}
		},
		"dsAlgorithm": func(id int) string {
			return dsAlgorithms[id]
		},
		"dsStatus": func(events []protocol.Event) protocol.Status {
			for _, event := range events {
				if event.Action == protocol.EventDelegationSignCheck && len(event.Status) > 0 {
					return event.Status[0]
				}
			}

			return protocol.Status("")
		},
		"dsLastCheck": func(events []protocol.Event) time.Time {
			for _, event := range events {
				if event.Action == protocol.EventDelegationSignCheck && len(event.Status) > 0 {
					return event.Date
				}
			}

			return time.Time{}
		},
		"dsLastOK": func(events []protocol.Event) time.Time {
			for _, event := range events {
				if event.Action == protocol.EventLastCorrectDelegationSignCheck {
					return event.Date
				}
			}

			return time.Time{}
		},
	}
)
