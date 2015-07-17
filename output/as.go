package output

import (
	"io"
	"text/template"
	"time"

	"github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/registrobr/rdap/protocol"
)

type AS struct {
	AS *protocol.AS

	CreatedAt time.Time
	UpdatedAt time.Time

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

func (a *AS) Print(wr io.Writer) error {
	a.setDates()
	addContacts(a, a.AS.Entities)
	filterContacts(a)

	t, err := template.New("as template").
		Funcs(genericFuncMap).
		Parse(asTmpl)

	if err != nil {
		return err
	}

	return t.Execute(wr, a)
}
