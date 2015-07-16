package output

import (
	"io"
	"text/template"

	"github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/registrobr/rdap/protocol"
)

type Domain struct {
	Domain *protocol.Domain

	CreatedAt string
	UpdatedAt string
	ExpiresAt string

	Handles       map[string]string
	DS            []ds
	ContactsInfos []contactInfo
}

type ds struct {
	protocol.DS
	CreatedAt string
}

func (d *Domain) addContact(c contactInfo) {
	d.ContactsInfos = append(d.ContactsInfos, c)
}

func (d *Domain) getContacts() []contactInfo {
	return d.ContactsInfos
}

func (d *Domain) setContacts(c []contactInfo) {
	d.ContactsInfos = c
}

func (d *Domain) setDates() {
	for _, e := range d.Domain.Events {
		date := e.Date.Format("20060102")

		switch e.Action {
		case protocol.EventActionRegistration:
			d.CreatedAt = date
		case protocol.EventActionLastChanged:
			d.UpdatedAt = date
		case protocol.EventActionExpiration:
			d.ExpiresAt = date
		}
	}
}

func (d *Domain) setDS() {
	d.DS = make([]ds, len(d.Domain.SecureDNS.DSData))

	for i, dsdatum := range d.Domain.SecureDNS.DSData {
		myds := ds{DS: dsdatum}

		for _, e := range dsdatum.Events {
			if e.Action == protocol.EventActionRegistration {
				myds.CreatedAt = e.Date.Format("20060102")
			}
		}

		d.DS[i] = myds
	}
}

func (d *Domain) Print(wr io.Writer) error {
	d.setDates()
	d.setDS()
	addContacts(d, d.Domain.Entities)
	filterContacts(d)

	t, err := template.New("domain template").
		Funcs(contactInfoFuncMap).
		Funcs(domainFuncMap).
		Parse(domainTmpl)

	if err != nil {
		return err
	}

	return t.ExecuteTemplate(wr, "domain template", d)
}
