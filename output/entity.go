package output

import (
	"io"
	"text/template"
	"time"

	"github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/registrobr/rdap/protocol"
)

type Entity struct {
	Entity *protocol.Entity

	CreatedAt time.Time
	UpdatedAt time.Time

	ContactsInfos []contactInfo
}

func (e *Entity) AddContact(c contactInfo) {
	e.ContactsInfos = append(e.ContactsInfos, c)
}

func (e *Entity) setDates() {
	for _, ev := range e.Entity.Events {
		switch ev.Action {
		case protocol.EventActionRegistration:
			e.CreatedAt = ev.Date
		case protocol.EventActionLastChanged:
			e.UpdatedAt = ev.Date
		}
	}
}

func (e *Entity) Print(wr io.Writer) error {
	e.setDates()
	var contactInfo contactInfo
	contactInfo.setContact(*e.Entity)
	e.ContactsInfos = append(e.ContactsInfos, contactInfo)

	t, err := template.New("entity template").
		Funcs(genericFuncMap).
		Parse(contactTmpl)

	if err != nil {
		return err
	}

	return t.Execute(wr, e)
}
