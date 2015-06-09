package output

import (
	"fmt"
	"strings"

	"github.com/registrobr/rdap/protocol"
)

const contactTmpl = `{{range .ContactsInfos}}handle:   {{.Handle}}
ids:      {{.Ids}}
roles:    {{.Roles}}
{{range .Persons}}person:   {{.}}
{{end}}{{range .Emails}}e-mail:   {{.}}
{{end}}{{range .Addresses}}address:  {{.}}
{{end}}{{range .Phones}}phone:    {{.}}
{{end}}created:  {{.CreatedAt}}
changed:  {{.UpdatedAt}}

{{end}}`

type ContactInfo struct {
	Handle    string
	Ids       string
	Persons   []string
	Emails    []string
	Addresses []string
	Phones    []string
	Roles     string
	CreatedAt string
	UpdatedAt string
}

func (c *ContactInfo) setContact(entity protocol.Entity) {
	c.Handle = entity.Handle
	for _, vCardValues := range entity.VCardArray {
		vCardValue, ok := vCardValues.([]interface{})
		if !ok {
			continue
		}

		for _, value := range vCardValue {
			v, ok := value.([]interface{})
			if !ok {
				continue
			}

			switch v[0] {
			case "fn":
				c.Persons = append(c.Persons, v[3].(string))
			case "email":
				c.Emails = append(c.Emails, v[3].(string))
			case "adr":
				address := make([]string, 0)

				addresses, ok := v[3].([]interface{})
				if !ok {
					continue
				}

				for _, v := range addresses {
					v := v.(string)

					if len(v) > 0 {
						address = append(address, v)
					}
				}

				c.Addresses = append(c.Addresses, strings.Join(address, ", "))
			case "tel":
				c.Phones = append(c.Phones, v[3].(string))
			}
		}
	}

	for _, event := range entity.Events {
		date := event.Date.Format("20060102")

		switch event.Action {
		case protocol.EventActionRegistration:
			c.CreatedAt = date
		case protocol.EventActionLastChanged:
			c.UpdatedAt = date
		}
	}

	c.Roles = strings.Join(entity.Roles, ", ")

	ids := make([]string, len(entity.PublicIds))

	for i, id := range entity.PublicIds {
		ids[i] = fmt.Sprintf("%s (%s)", id.Identifier, id.Type)
	}

	c.Ids = strings.Join(ids, ", ")
}

type ContactList interface {
	AddContact(ContactInfo)
}

func AddContacts(c ContactList, entities []protocol.Entity) {
	contacts := make(map[string]bool)
	for _, entity := range entities {
		if contacts[entity.Handle] {
			continue
		}
		contacts[entity.Handle] = true

		var contactInfo ContactInfo
		contactInfo.setContact(entity)
		c.AddContact(contactInfo)
	}
}
