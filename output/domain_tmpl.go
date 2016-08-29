package output

import (
	"text/template"

	"github.com/registrobr/rdap/protocol"
)

const (
	domainTmpl = `
domain:   {{.Domain.LDHName}}
{{range .Domain.Nameservers}}\
nserver:  {{.LDHName}}
{{$lastCheck := nsLastCheck .Events}}\
{{if (isDateDefined $lastCheck)}}\
nsstat:   {{$lastCheck | formatDate}} {{nsStatus .Events}}
{{end}}\
{{$lastOK := nsLastOK .Events}}\
{{if (isDateDefined $lastOK)}}\
nslastaa: {{$lastOK | formatDate}}
{{end}}\
{{end}}\
{{range .DS}}\
dsrecord: {{.KeyTag}} {{.Algorithm | dsAlgorithm}} {{.Digest}}
dsstatus: {{dsLastCheck .Events | formatDate}} {{dsStatus .Events}}
dslastok: {{dsLastOK .Events | formatDate}}
{{end}}\
{{if (isDateDefined .CreatedAt)}}\
created:  {{.CreatedAt | formatDate}}
{{end}}\
{{if (isDateDefined .UpdatedAt)}}\
changed:  {{.UpdatedAt | formatDate}}
{{end}}\
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
		"nsLastCheck": func(events []protocol.Event) protocol.EventDate {
			for _, event := range events {
				if event.Action == protocol.EventDelegationCheck && len(event.Status) > 0 {
					return event.Date
				}
			}

			return protocol.EventDate{}
		},
		"nsLastOK": func(events []protocol.Event) protocol.EventDate {
			for _, event := range events {
				if event.Action == protocol.EventLastCorrectDelegationCheck {
					return event.Date
				}
			}

			return protocol.EventDate{}
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
		"dsLastCheck": func(events []protocol.Event) protocol.EventDate {
			for _, event := range events {
				if event.Action == protocol.EventDelegationSignCheck && len(event.Status) > 0 {
					return event.Date
				}
			}

			return protocol.EventDate{}
		},
		"dsLastOK": func(events []protocol.Event) protocol.EventDate {
			for _, event := range events {
				if event.Action == protocol.EventLastCorrectDelegationSignCheck {
					return event.Date
				}
			}

			return protocol.EventDate{}
		},
	}
)
