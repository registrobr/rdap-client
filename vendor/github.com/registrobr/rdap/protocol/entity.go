package protocol

// PublicID describes Public IDs as it is in RFC 7483, section 4.8
type PublicID struct {
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
}

// CustomerSupportService is a NIC.br extension to store some extra contact
// information for an entity
type CustomerSupportService struct {
	Email   string `json:"nicbr_email,omitempty"`
	Website string `json:"nicbr_website,omitempty"`
	Phone   string `json:"nicbr_phone,omitempty"`
}

// Entity describes the Entity Object Class as it is in RFC 7483, section 5.1
type Entity struct {
	ObjectClassName        string                  `json:"objectClassName"`
	Handle                 string                  `json:"handle"`
	VCardArray             []interface{}           `json:"vcardArray,omitempty"`
	Roles                  []string                `json:"roles,omitempty"`
	PublicIds              []PublicID              `json:"publicIds,omitempty"`
	Networks               []IPNetwork             `json:"networks,omitempty"`
	Autnums                []AS                    `json:"autnums,omitempty"`
	CustomerSupportService *CustomerSupportService `json:"nicbr_customerSupportService,omitempty"`
	Entities               []Entity                `json:"entities,omitempty"`
	Events                 []Event                 `json:"events,omitempty"`
	Links                  []Link                  `json:"links,omitempty"`
	Remarks                []Remark                `json:"remarks,omitempty"`
	Notices                []Notice                `json:"notices,omitempty"`
	DomainCount            int                     `json:"nicbr_domainCount,omitempty"`
	InetCount              int                     `json:"nicbr_inetCount,omitempty"`
	AutnumCount            int                     `json:"nicbr_autnumCount,omitempty"`
	Lang                   string                  `json:"lang,omitempty"`
	Conformance
	Port43

	// LegalRepresentative was proposed by NIC.br to store the name of the
	// persons that is responsible for this entity
	LegalRepresentative string `json:"legalRepresentative,omitempty"`
}

// GetEntity is an easy way to find an entity with a given role. If more than
// one entity has the same role, the last one is returned
func (e *Entity) GetEntity(role string) (entity Entity, found bool) {
	for _, v := range e.Entities {
		for _, r := range v.Roles {
			if r == role {
				entity = v
				found = true
				return
			}
		}
	}

	return
}
