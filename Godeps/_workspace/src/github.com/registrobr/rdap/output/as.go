package output

import (
	"io"
	"text/template"

	"github.com/registrobr/rdap/protocol"
)

type AS struct {
	AS *protocol.ASResponse

	CreatedAt string
	UpdatedAt string

	ContactsInfos []ContactInfo
}

func (a *AS) AddContact(c ContactInfo) {
	a.ContactsInfos = append(a.ContactsInfos, c)
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

func (a *AS) ToText(wr io.Writer) error {
	a.setDates()
	AddContacts(a, a.AS.Entities)

	for _, entity := range a.AS.Entities {
		AddContacts(a, entity.Entities)
	}

	t, err := template.New("as template").Parse(asTmpl)
	if err != nil {
		return err
	}

	return t.Execute(wr, a)
}
