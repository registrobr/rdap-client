package output

import (
	"io"
	"strings"
	"text/template"
	"time"

	"github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/registrobr/rdap/protocol"
)

type Domain struct {
	Domain *protocol.Domain

	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time

	Handles       map[string]string
	DS            []ds
	ContactsInfos []contactInfo
}

type ds struct {
	protocol.DS
	CreatedAt time.Time
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
		switch e.Action {
		case protocol.EventActionRegistration:
			d.CreatedAt = e.Date
		case protocol.EventActionLastChanged:
			d.UpdatedAt = e.Date
		case protocol.EventActionExpiration:
			d.ExpiresAt = e.Date
		}
	}
}

func (d *Domain) setDS() {
	if d.Domain.SecureDNS == nil {
		return
	}

	d.DS = make([]ds, len(d.Domain.SecureDNS.DSData))

	for i, dsdatum := range d.Domain.SecureDNS.DSData {
		myds := ds{DS: dsdatum}

		for _, e := range dsdatum.Events {
			if e.Action == protocol.EventActionRegistration {
				myds.CreatedAt = e.Date
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
		Funcs(genericFuncMap).
		Funcs(domainFuncMap).
		Parse(strings.Replace(domainTmpl, "\\\n", "", -1))

	if err != nil {
		return err
	}

	return t.ExecuteTemplate(wr, "domain template", d)
}
