package output

import (
	"io"
	"net"
	"text/template"

	"github.com/registrobr/rdap/protocol"
)

type IPNetwork struct {
	IPNetwork *protocol.IPNetwork
	Inetnum   string

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

func (i *IPNetwork) setInetnum() {
	start := net.ParseIP(i.IPNetwork.StartAddress)
	end := net.ParseIP(i.IPNetwork.EndAddress)
	mask := make(net.IPMask, len(start))

	for j := 0; j < len(start); j++ {
		mask[j] = start[j] | ^end[j]
	}

	cidr := net.IPNet{start, mask}

	i.Inetnum = cidr.String()
}

func (i *IPNetwork) Print(wr io.Writer) error {
	i.setDates()
	i.setInetnum()
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
