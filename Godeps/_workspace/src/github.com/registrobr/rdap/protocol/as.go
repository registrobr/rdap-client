package protocol

type AS struct {
	ObjectClassName string          `json:"objectClassName"`
	Handle          string          `json:"handle"`
	StartAutnum     uint32          `json:"startAutNum"`
	EndAutnum       uint32          `json:"endAutNum"`
	Type            string          `json:"type"`
	Country         string          `json:"country"`
	Links           []Link          `json:"links,omitempty"`
	Entities        []Entity        `json:"entities,omitempty"`
	RoutingPolicy   []RoutingPolicy `json:"nicbr_routing_policy,omitempty"`
	Events          []Event         `json:"events,omitempty"`
}

type RoutingPolicy struct {
	Autnum  uint32 `json:"autnum"`
	Traffic int    `json:"traffic"`
	Cost    int    `json:"cost"`
	Policy  string `json:"policy"`
}
