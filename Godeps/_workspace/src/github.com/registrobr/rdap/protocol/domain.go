package protocol

type Domain struct {
	ObjectClassName string                `json:"objectClassName"`
	Handle          string                `json:"handle,omitempty"`
	LDHName         string                `json:"ldhName,omitempty"`
	UnicodeName     string                `json:"unicodeName,omitempty"`
	Nameservers     []Nameserver          `json:"nameservers,omitempty"`
	SecureDNS       *SecureDNS            `json:"secureDNS,omitempty"`
	Arbitration     bool                  `json:"nicbr_arbitration,omitempty"`
	Links           []Link                `json:"links,omitempty"`
	Entities        []Entity              `json:"entities,omitempty"`
	Events          []Event               `json:"events,omitempty"`
	Status          []Status              `json:"status,omitempty"`
	PublicIDs       []PublicID            `json:"publicIds,omitempty"`
	Remarks         []Remark              `json:"remarks,omitempty"`
	Notices         []Notice              `json:"notices,omitempty"`
	Network         *IPNetwork            `json:"network,omitempty"`
	Unavailability  *DomainUnavailability `json:"-"`
	Conformance
}

type DomainUnavailability struct {
	Reason     string `json:"nicbr_reason"`
	INPINumber int    `json:"nicbr_inpiNumber,omitempty"`
}
