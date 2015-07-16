package protocol

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
}

type RoutingPolicy struct {
	Autnum  uint32 `json:"nicbr_autnum"`
	Traffic int    `json:"nicbr_traffic"`
	Cost    int    `json:"nicbr_cost"`
	Policy  string `json:"nicbr_policy"`
}
