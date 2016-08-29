package output

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/registrobr/rdap/protocol"
)

type AS struct {
	AS            *protocol.AS
	CreatedAt     protocol.EventDate
	UpdatedAt     protocol.EventDate
	IPNetworks    []string
	ContactsInfos []contactInfo
}

func (a *AS) addContact(c contactInfo) {
	a.ContactsInfos = append(a.ContactsInfos, c)
}

func (a *AS) getContacts() []contactInfo {
	return a.ContactsInfos
}

func (a *AS) setContacts(c []contactInfo) {
	a.ContactsInfos = c
}

func (a *AS) setDates() {
	for _, e := range a.AS.Events {
		switch e.Action {
		case protocol.EventActionRegistration:
			a.CreatedAt = e.Date
		case protocol.EventActionLastChanged:
			a.UpdatedAt = e.Date
		}
	}
}

func (a *AS) setIPNetworks() {
	for _, l := range a.AS.Links {
		if l.Rel != "related" {
			continue
		}

		linkParts := strings.Split(l.Href, "/")
		if len(linkParts) >= 3 && linkParts[len(linkParts)-3] == "ip" {
			cidr := fmt.Sprintf("%s/%s", linkParts[len(linkParts)-2], linkParts[len(linkParts)-1])
			a.IPNetworks = append(a.IPNetworks, cidr)
		}
	}
}

func (a *AS) Print(wr io.Writer) error {
	a.setDates()
	a.setIPNetworks()
	addContacts(a, a.AS.Entities)
	filterContacts(a)

	t, err := template.New("as template").
		Funcs(genericFuncMap).
		Parse(strings.Replace(asTmpl, "\\\n", "", -1))

	if err != nil {
		return err
	}

	return t.Execute(wr, a)
}
