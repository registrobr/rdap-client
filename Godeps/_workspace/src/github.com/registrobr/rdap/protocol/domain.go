package protocol

// Domain describes Domain Object Class as it is in RFC 7483, section 5.3
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

// DomainUnavailability is a NIC.br extension used to determinate the reason
// that a domain name cannot be registered
type DomainUnavailability struct {
	Reason     string `json:"nicbr_reason"`
	INPINumber int    `json:"nicbr_inpiNumber,omitempty"`
}
