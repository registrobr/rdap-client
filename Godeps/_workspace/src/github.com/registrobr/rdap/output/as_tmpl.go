package output

const asTmpl = `aut-num:     {{.AS.Handle}}
country:     {{.AS.Country}}
created:     {{.CreatedAt}}
changed:     {{.UpdatedAt}}

inetnum:     (ip networks)

` + contactTmpl
