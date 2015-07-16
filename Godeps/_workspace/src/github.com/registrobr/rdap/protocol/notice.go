package protocol

type Notice struct {
	Title       string   `json:"title,omitempty"`
	Description []string `json:"description,omitempty"`
	Links       []Link   `json:"links,omitempty"`
}
