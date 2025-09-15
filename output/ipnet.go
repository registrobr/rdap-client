package output

import (
	"io"
	"strings"
	"text/template"

	"github.com/registrobr/rdap/protocol"
)

type IPNetwork struct {
	IPNetwork     *protocol.IPNetwork
	CreatedAt     protocol.EventDate
	UpdatedAt     protocol.EventDate
	ContactsInfos []contactInfo
}

func (i *IPNetwork) addContact(c contactInfo) {
	i.ContactsInfos = append(i.ContactsInfos, c)
}

func (i *IPNetwork) getContacts() []contactInfo {
	return i.ContactsInfos
}

func (i *IPNetwork) setContacts(c []contactInfo) {
	i.ContactsInfos = c
}

func (i *IPNetwork) setDates() {
	for _, e := range i.IPNetwork.Events {
		switch e.Action {
		case protocol.EventActionRegistration:
			i.CreatedAt = e.Date
		case protocol.EventActionLastChanged:
			i.UpdatedAt = e.Date
		}
	}
}

func (i *IPNetwork) Print(wr io.Writer) error {
	i.setDates()
	addContacts(i, i.IPNetwork.Entities)
	filterContacts(i)

	t, err := template.New("ipnetwork template").
		Funcs(genericFuncMap).
		Funcs(ipnetFuncMap).
		Parse(strings.ReplaceAll(ipnetTmpl, "\\\n", ""))

	if err != nil {
		return err
	}

	return t.Execute(wr, i)
}
