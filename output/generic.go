package output

import (
	"strings"
	"text/template"
	"time"
)

var (
	genericFuncMap = template.FuncMap{
		"isDateDefined": func(time time.Time) bool {
			return time.IsZero()
		},
		"formatDate": func(time time.Time) string {
			return time.Format(dateFormat)
		},
		"join": func(in []string) string {
			return strings.Join(in, ", ")
		},
	}
)
