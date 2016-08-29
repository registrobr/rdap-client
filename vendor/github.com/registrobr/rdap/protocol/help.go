package protocol

// Help describes an answer to help queries as it is in RFC 7483, section 7
type Help struct {
	Notices []Notice `json:"notices,omitempty"`
	Conformance
	Port43
}
