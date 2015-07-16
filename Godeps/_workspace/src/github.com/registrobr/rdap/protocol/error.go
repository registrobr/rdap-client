package protocol

type Error struct {
	Notices     []Notice `json:"notices,omitempty"`
	Lang        string   `json:"lang,omitempty"`
	ErrorCode   int      `json:"errorCode,omitempty"`
	Title       string   `json:"title,omitempty"`
	Description []string `json:"description,omitempty"`
	Conformance
}
