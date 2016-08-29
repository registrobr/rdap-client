package protocol

// Notice describes Notices as it is in RFC 7483, section 4.3
type Notice struct {
	Title       string   `json:"title,omitempty"`
	Description []string `json:"description,omitempty"`
	Links       []Link   `json:"links,omitempty"`
}
