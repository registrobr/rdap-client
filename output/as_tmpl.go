package output

const asTmpl = `aut-num:     {{.AS.Handle}}
country:     {{.AS.Country}}
created:     {{formatDate .CreatedAt}}
changed:     {{formatDate .UpdatedAt}}

inetnum:     (ip networks)

` + contactTmpl
