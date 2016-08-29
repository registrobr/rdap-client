package output

import (
	"strings"
	"text/template"

	"github.com/registrobr/rdap/protocol"
)

var (
	genericFuncMap = template.FuncMap{
		"isDateDefined": func(time protocol.EventDate) bool {
			return !time.IsZero()
		},
		"formatDate": func(time protocol.EventDate) string {
			return time.Format(dateFormat)
		},
		"join": func(in []string) string {
			return strings.Join(in, ", ")
		},
	}
)
