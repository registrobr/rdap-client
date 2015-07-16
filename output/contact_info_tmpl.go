package output

import (
	"strings"
	"text/template"
)

const contactTmpl = `{{range .ContactsInfos}}handle:   {{.Handle}}
{{if len .Ids}}ids:      {{.Ids | join}}
{{end}}roles:    {{.Roles | join}}
{{range .Persons}}person:   {{.}}
{{end}}{{range .Emails}}e-mail:   {{.}}
{{end}}{{range .Addresses}}address:  {{.}}
{{end}}{{range .Phones}}phone:    {{.}}
{{end}}created:  {{.CreatedAt}}
changed:  {{.UpdatedAt}}

{{end}}`

var (
	contactInfoFuncMap = template.FuncMap{
		"join": func(in []string) string {
			return strings.Join(in, ", ")
		},
	}
)
