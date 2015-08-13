package protocol

// IPNetwork describes the IP Network Object Class as it is in RFC 7483,
// section 5.4
type IPNetwork struct {
	ObjectClassName    string              `json:"objectClassName"`
	Handle             string              `json:"handle"`
	StartAddress       string              `json:"startAddress"`
	EndAddress         string              `json:"endAddress"`
	IPVersion          string              `json:"ipVersion"`
	Name               string              `json:"name,omitempty"`
	Type               string              `json:"type"`
	Country            string              `json:"country"`
	ParentHandle       string              `json:"parentHandle,omitempty"`
	Status             []string            `json:"status"`
	Autnum             uint32              `json:"nicbr_autnum,omitempty"`
	Links              []Link              `json:"links"`
	Events             []Event             `json:"events"`
	Entities           []Entity            `json:"entities"`
	Notices            []Notice            `json:"notices,omitempty"`
	Remarks            []Remark            `json:"remarks,omitempty"`
	ReverseDelegations []ReverseDelegation `json:"nicbr_reverseDelegations,omitempty"`
	Conformance
}

// ReverseDelegation is a NIC.br extension to list all the IP network
// delegations and the corresponding nameservers
type ReverseDelegation struct {
	StartAddress string       `json:"startAddress"`
	EndAddress   string       `json:"endAddress"`
	Nameservers  []Nameserver `json:"nameservers,omitempty"`
	Events       []Event      `json:"events"`
}
