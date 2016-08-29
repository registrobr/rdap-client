package protocol

// AS describes the Autonomous System Number Entity Object Class as it is in
// RFC 7483, section 5.5
type AS struct {
	ObjectClassName string          `json:"objectClassName"`
	Handle          string          `json:"handle"`
	StartAutnum     uint32          `json:"startAutnum"`
	EndAutnum       uint32          `json:"endAutnum"`
	Name            string          `json:"name,omitempty"`
	Type            string          `json:"type"`
	Country         string          `json:"country"`
	Links           []Link          `json:"links,omitempty"`
	Entities        []Entity        `json:"entities,omitempty"`
	RoutingPolicy   []RoutingPolicy `json:"nicbr_routingPolicy,omitempty"`
	Events          []Event         `json:"events,omitempty"`
	Notices         []Notice        `json:"notices,omitempty"`
	Remarks         []Remark        `json:"remarks,omitempty"`
	Conformance
	Port43
}

// RoutingPolicy is a NIC.br extension that stores the information of network
// announces
type RoutingPolicy struct {
	Autnum  uint32 `json:"nicbr_autnum"`
	Traffic int    `json:"nicbr_traffic"`
	Cost    int    `json:"nicbr_cost"`
	Policy  string `json:"nicbr_policy"`
}
