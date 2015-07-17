package output

import (
	"strings"
	"text/template"
	"time"
)

var (
	genericFuncMap = template.FuncMap{
		"formatDate": func(time time.Time) string {
			return time.Format(dateFormat)
		},
		"join": func(in []string) string {
			return strings.Join(in, ", ")
		},
	}
)
