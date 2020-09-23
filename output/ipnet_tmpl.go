package output

import (
	"net"
	"text/template"

	rdap "github.com/registrobr/rdap/protocol"
)

var ipnetTmpl = `
inetnum:       {{inetnum .IPNetwork.StartAddress .IPNetwork.EndAddress}}
handle:        {{.IPNetwork.Handle}}
{{if ne .IPNetwork.ParentHandle ""}}\
parent-handle: {{.IPNetwork.ParentHandle}}
{{end}}\
{{if gt .IPNetwork.Autnum 0}}\
aut-num:       {{.IPNetwork.Autnum}}
{{end}}\
start-address: {{.IPNetwork.StartAddress}}
end-address:   {{.IPNetwork.EndAddress}}
ip-version:    {{.IPNetwork.IPVersion}}
name:          {{.IPNetwork.Name}}
{{if ne .IPNetwork.Type ""}}\
type:          {{.IPNetwork.Type}}
{{end}}\
{{if ne .IPNetwork.Country ""}}\
country:       {{.IPNetwork.Country}}
{{end}}\
{{range .IPNetwork.Status}}\
status:        {{.}}
{{end}}\
{{range .IPNetwork.ReverseDelegations}}\
{{ $startAddress := .StartAddress}}
{{ $endAddress := .EndAddress }}
inetrev:       {{inetnum $startAddress $endAddress}}
{{range .Nameservers}}\
nserver:       {{.LDHName}}
{{end}}\
{{ if hasSecureDns .SecureDNS}}
{{ range .SecureDNS.DSSet }}
dsinetrev:     {{inetnum $startAddress $endAddress}}
dsrecord:      {{.Keytag}}{{.Digest}}
{{ range .Events }}
{{ if and (eq .Action "delegation sign check") (gt (lenStatus .Status) 0)}}
dsstatus:      {{ .Date.Time | formatDate }}{{dsStatusTranslate (index .Status 0)}}
{{ else if eq .Action "last correct delegation sign check" }}
dslastok: {{ .Date.Time | formatDate }}
{{ end }}
{{ end }}
{{ end }}
{{ end }}\
{{end}}\
{{if (isDateDefined .CreatedAt)}}\
created:       {{.CreatedAt | formatDate}}
{{end}}\
{{if (isDateDefined .UpdatedAt)}}\
changed:       {{.UpdatedAt | formatDate}}
{{end}}\

` + contactTmpl

var (
	ipnetFuncMap = template.FuncMap{
		"inetnum": func(startAddress, endAddress string) string {
			start := net.ParseIP(startAddress)
			end := net.ParseIP(endAddress)
			mask := make(net.IPMask, len(start))

			for j := 0; j < len(start); j++ {
				mask[j] = start[j] | ^end[j]
			}

			cidr := net.IPNet{IP: start, Mask: mask}
			return cidr.String()
		},
		"lenStatus": func(s []rdap.Status) int {
			return len(s)
		},
		"dsStatusTranslate": func(rs rdap.Status) string {
			switch rs {
			case rdap.StatusDSOK:
				return "OK"
			case rdap.StatusDSTimeout:
				return "TIMEOUT"
			case rdap.StatusDSNoSig:
				return "NOSIG"
			case rdap.StatusDSExpiredSig:
				return "EXPSIG"
			case rdap.StatusDSInvalidSig:
				return "SIGERROR"
			case rdap.StatusDSNotFound:
				return "NOKEY"
			case rdap.StatusDSNoSEP:
				return "NOSEP"
			}

			return "PLAIN DNS ERROR"
		},
		"hasSecureDns": func(secdns *rdap.ReverseDelegationSecureDNS) bool {
			return secdns != nil
		},
	}
)
