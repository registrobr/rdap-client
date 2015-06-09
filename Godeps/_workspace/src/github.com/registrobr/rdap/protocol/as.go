package protocol

type ASResponse struct {
	ObjectClassName string   `json:"objectClassName"`
	Handle          string   `json:"handle"`
	StartAutnum     uint32   `json:"startAutNum"`
	EndAutnum       uint32   `json:"endAutNum"`
	Type            string   `json:"type"`
	Country         string   `json:"country"`
	Links           []Link   `json:"links,omitempty"`
	Entities        []Entity `json:"entities,omitempty"`
	Events          []Event  `json:"events,omitempty"`
}
