package output

import (
	"io"
	"text/template"

	"github.com/registrobr/rdap/protocol"
)

type Domain struct {
	Domain *protocol.DomainResponse

	CreatedAt string
	UpdatedAt string
	ExpiresAt string

	Handles       map[string]string
	DS            []ds
	ContactsInfos []ContactInfo
}

type ds struct {
	protocol.DS
	CreatedAt string
}

func (d *Domain) AddContact(c ContactInfo) {
	d.ContactsInfos = append(d.ContactsInfos, c)
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

func (d *Domain) ToText(wr io.Writer) error {
	d.setDates()
	d.setDS()

	AddContacts(d, d.Domain.Entities)

	for _, entity := range d.Domain.Entities {
		AddContacts(d, entity.Entities)
	}

	t, err := template.New("domain").Funcs(domainFuncMap).Parse(domainTmpl)

	if err != nil {
		return err
	}

	return t.ExecuteTemplate(wr, "domain", d)
}
