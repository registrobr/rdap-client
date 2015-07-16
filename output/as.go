package output

import (
	"io"
	"text/template"

	"github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/registrobr/rdap/protocol"
)

type AS struct {
	AS *protocol.AS

	CreatedAt string
	UpdatedAt string

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
		date := e.Date.Format("20060102")

		switch e.Action {
		case protocol.EventActionRegistration:
			a.CreatedAt = date
		case protocol.EventActionLastChanged:
			a.UpdatedAt = date
		}
	}
}

func (a *AS) Print(wr io.Writer) error {
	a.setDates()
	addContacts(a, a.AS.Entities)
	filterContacts(a)

	t, err := template.New("as template").
		Funcs(contactInfoFuncMap).
		Parse(asTmpl)

	if err != nil {
		return err
	}

	return t.Execute(wr, a)
}
