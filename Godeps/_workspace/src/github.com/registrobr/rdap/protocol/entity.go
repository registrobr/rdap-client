package protocol

type PublicID struct {
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
}

type CustomerSupportService struct {
	Email   string `json:"nicbr_email,omitempty"`
	Website string `json:"nicbr_website,omitempty"`
	Phone   string `json:"nicbr_phone,omitempty"`
}

type Entity struct {
	ObjectClassName        string                  `json:"objectClassName"`
	Handle                 string                  `json:"handle"`
	VCardArray             []interface{}           `json:"vcardArray,omitempty"`
	Roles                  []string                `json:"roles,omitempty"`
	PublicIds              []PublicID              `json:"publicIds,omitempty"`
	Responsible            string                  `json:"nicbr_responsible,omitempty"`
	CustomerSupportService *CustomerSupportService `json:"nicbr_customerSupportService,omitempty"`
	Entities               []Entity                `json:"entities,omitempty"`
	Events                 []Event                 `json:"events,omitempty"`
	Links                  []Link                  `json:"links,omitempty"`
	Remarks                []Remark                `json:"remarks,omitempty"`
	Notices                []Notice                `json:"notices,omitempty"`
	DomainCount            int                     `json:"nicbr_domainCount,omitempty"`
	InetCount              int                     `json:"nicbr_inetCount,omitempty"`
	AutnumCount            int                     `json:"nicbr_autnumCount,omitempty"`
	Conformance
}
