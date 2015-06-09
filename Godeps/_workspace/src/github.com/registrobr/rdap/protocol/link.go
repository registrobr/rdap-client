package protocol

type Link struct {
	Value string `json:"value,omitempty"`
	Rel   string `json:"rel,omitempty"`
	Href  string `json:"href,omitempty"`
	Type  string `json:"type,omitempty"`
}
