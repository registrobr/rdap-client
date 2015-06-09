package output

import (
	"io"
	"text/template"

	"github.com/registrobr/rdap/protocol"
)

type IPNetwork struct {
	IPNetwork *protocol.IPNetwork

	CreatedAt string
	UpdatedAt string

	ContactsInfos []ContactInfo
}

func (i *IPNetwork) AddContact(c ContactInfo) {
	i.ContactsInfos = append(i.ContactsInfos, c)
}

func (i *IPNetwork) setDates() {
	for _, e := range i.IPNetwork.Events {
		date := e.Date.Format("20060102")

		switch e.Action {
		case protocol.EventActionRegistration:
			i.CreatedAt = date
		case protocol.EventActionLastChanged:
			i.UpdatedAt = date
		}
	}
}

func (i *IPNetwork) ToText(wr io.Writer) error {
	i.setDates()
	AddContacts(i, i.IPNetwork.Entities)

	for _, entity := range i.IPNetwork.Entities {
		AddContacts(i, entity.Entities)
	}

	t, err := template.New("ipnetwork template").Parse(ipnetTmpl)
	if err != nil {
		return err
	}

	return t.Execute(wr, i)
}
