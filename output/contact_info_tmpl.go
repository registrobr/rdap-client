package output

const contactTmpl = `{{range .ContactsInfos}}handle:   {{.Handle}}
{{if len .Ids}}ids:      {{.Ids | join}}
{{end}}{{if len .Roles}}roles:    {{.Roles | join}}
{{end}}{{range .Persons}}person:   {{.}}
{{end}}{{range .Emails}}e-mail:   {{.}}
{{end}}{{range .Addresses}}address:  {{.}}
{{end}}{{range .Phones}}phone:    {{.}}
{{end}}created:  {{.CreatedAt}}
changed:  {{.UpdatedAt}}

{{end}}`
