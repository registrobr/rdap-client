package output

import (
	"net"
	"text/template"
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
inetrev:       {{inetnum .StartAddress .EndAddress}}
{{range .Nameservers}}\
nserver:       {{.LDHName}}
{{end}}\
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
	}
)
