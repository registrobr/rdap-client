package output

import (
	"fmt"
	"strings"

	"github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/registrobr/rdap/protocol"
)

type contactInfo struct {
	Handle    string
	Ids       []string
	Persons   []string
	Emails    []string
	Addresses []string
	Phones    []string
	Roles     []string
	CreatedAt string
	UpdatedAt string
}

func (c *contactInfo) setContact(entity protocol.Entity) {
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
				var address []string

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

	c.Roles = entity.Roles

	for _, id := range entity.PublicIds {
		c.Ids = append(c.Ids, fmt.Sprintf("%s (%s)", id.Identifier, id.Type))
	}
}

type contactList interface {
	addContact(contactInfo)
	getContacts() []contactInfo
	setContacts(c []contactInfo)
}

func addContacts(c contactList, entities []protocol.Entity) {
	for _, entity := range entities {
		var contactInfo contactInfo
		contactInfo.setContact(entity)
		c.addContact(contactInfo)

		addContacts(c, entity.Entities)
	}
}

func filterContacts(c contactList) {
	contacts := make(map[string]*contactInfo)

	for _, contactInfo := range c.getContacts() {
		contactInfo := contactInfo

		if _, ok := contacts[contactInfo.Handle]; !ok {
			contacts[contactInfo.Handle] = &contactInfo
			continue
		}

		contacts[contactInfo.Handle].Roles = append(contacts[contactInfo.Handle].Roles,
			contactInfo.Roles...)
	}

	for _, contactInfo := range contacts {
		found := make(map[string]bool)
		roles := make([]string, 0)

		for _, role := range contactInfo.Roles {
			if _, ok := found[role]; ok {
				continue
			}

			roles = append(roles, role)
			found[role] = true
		}

		contactInfo.Roles = roles
	}

	filteredContacts := make([]contactInfo, 0)

	for _, contact := range contacts {
		filteredContacts = append(filteredContacts, *contact)
	}

	c.setContacts(filteredContacts)
}
