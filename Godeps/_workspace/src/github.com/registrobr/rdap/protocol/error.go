package protocol

type Error struct {
	RDAPConformance []string `json:"rdapConformance,omitempty"`
	Notices         []Notice `json:"notices,omitempty"`
	Lang            string   `json:"lang,omitempty"`
	ErrorCode       int
	Title           string
	Description     []string
}
