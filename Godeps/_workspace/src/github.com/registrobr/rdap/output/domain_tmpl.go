package output

import "text/template"

const domainTmpl = `domain:   {{.Domain.LDHName}}
{{range .Domain.Nameservers}}nserver:  {{.LDHName}} {{.HostStatus}}
{{end}}{{range .DS}}dsrecord: {{.KeyTag}} {{.Algorithm | dsAlgorithm}} {{.Digest}}
dsstatus: {{.CreatedAt}} {{.DSStatus}}
{{end}}created:  {{.CreatedAt}}
changed:  {{.UpdatedAt}}
{{range .Domain.Status}}status:   {{.}}
{{end}}
` + contactTmpl

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
		"dsAlgorithm": func(id int) string {
			return dsAlgorithms[id]
		},
	}
)
