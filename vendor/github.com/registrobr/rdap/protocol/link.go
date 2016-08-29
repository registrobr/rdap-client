package protocol

// Link describes Links as it is in RFC 7483, section 4.2
type Link struct {
	Value string `json:"value,omitempty"`
	Rel   string `json:"rel,omitempty"`
	Href  string `json:"href,omitempty"`
	Type  string `json:"type,omitempty"`
}
