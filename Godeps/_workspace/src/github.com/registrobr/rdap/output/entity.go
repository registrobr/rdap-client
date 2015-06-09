package output

import (
	"io"
	"text/template"

	"github.com/registrobr/rdap/protocol"
)

type Entity struct {
	Entity *protocol.Entity

	CreatedAt string
	UpdatedAt string

	ContactsInfos []ContactInfo
}

func (e *Entity) AddContact(c ContactInfo) {
	e.ContactsInfos = append(e.ContactsInfos, c)
}

func (e *Entity) setDates() {
	for _, event := range e.Entity.Events {
		date := event.Date.Format("20060102")

		switch event.Action {
		case protocol.EventActionRegistration:
			e.CreatedAt = date
		case protocol.EventActionLastChanged:
			e.UpdatedAt = date
		}
	}
}

func (e *Entity) ToText(wr io.Writer) error {
	e.setDates()
	var contactInfo ContactInfo
	contactInfo.setContact(*e.Entity)
	e.ContactsInfos = append(e.ContactsInfos, contactInfo)

	t, err := template.New("entity template").Parse(contactTmpl)
	if err != nil {
		return err
	}

	return t.Execute(wr, e)
}
