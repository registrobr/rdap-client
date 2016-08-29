package protocol

// IPAddresses describes the ipAddresses field as it is in RFC 7483, section
// 5.2
type IPAddresses struct {
	V4 []string `json:"v4,omitempty"`
	V6 []string `json:"v6,omitempty"`
}

// Nameserver describes the Nameserver Object Class as it is in RFC 7483,
// section 5.2
type Nameserver struct {
	ObjectClassName string       `json:"objectClassName"`
	Handle          string       `json:"handle,omitempty"`
	LDHName         string       `json:"ldhName,omitempty"`
	UnicodeName     string       `json:"unicodeName,omitempty"`
	Entities        []Entity     `json:"entities,omitempty"`
	Status          []Status     `json:"status,omitempty"`
	IPAddresses     *IPAddresses `json:"ipAddresses,omitempty"`
	Remarks         []Remark     `json:"remarks,omitempty"`
	Links           []Link       `json:"links,omitempty"`
	Port43          string       `json:"port43,omitempty"`
	Events          []Event      `json:"events,omitempty"`
}
