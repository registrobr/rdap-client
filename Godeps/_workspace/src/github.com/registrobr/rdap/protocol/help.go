package protocol

type Help struct {
	Notices []Notice `json:"notices,omitempty"`
	Conformance
}
