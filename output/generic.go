package output

import (
	"strings"
	"text/template"
	"time"

	"github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/registrobr/rdap/protocol"
)

var (
	genericFuncMap = template.FuncMap{
		"formatDate": func(time time.Time) string {
			return time.Format(dateFormat)
		},
		"join": func(in []string) string {
			return strings.Join(in, ", ")
		},
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
