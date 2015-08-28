package output

var ipnetTmpl = `inetnum:       {{.Inetnum}}
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
created:       {{.CreatedAt}}
changed:       {{.UpdatedAt}}

` + contactTmpl
