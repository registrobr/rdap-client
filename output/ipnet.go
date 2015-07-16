package output

import (
	"io"
	"net"
	"text/template"

	"github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/registrobr/rdap/protocol"
)

type IPNetwork struct {
	IPNetwork *protocol.IPNetwork
	Inetnum   string

	CreatedAt string
	UpdatedAt string

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
		date := e.Date.Format("20060102")

		switch e.Action {
		case protocol.EventActionRegistration:
			i.CreatedAt = date
		case protocol.EventActionLastChanged:
			i.UpdatedAt = date
		}
	}
}

func (i *IPNetwork) setInetnum() {
	start := net.ParseIP(i.IPNetwork.StartAddress)
	end := net.ParseIP(i.IPNetwork.EndAddress)
	mask := make(net.IPMask, len(start))

	for j := 0; j < len(start); j++ {
		mask[j] = start[j] | ^end[j]
	}

	cidr := net.IPNet{IP: start, Mask: mask}

	i.Inetnum = cidr.String()
}

func (i *IPNetwork) Print(wr io.Writer) error {
	i.setDates()
	i.setInetnum()
	addContacts(i, i.IPNetwork.Entities)
	filterContacts(i)

	t, err := template.New("ipnetwork template").
		Funcs(contactInfoFuncMap).
		Parse(ipnetTmpl)

	if err != nil {
		return err
	}

	return t.Execute(wr, i)
}
