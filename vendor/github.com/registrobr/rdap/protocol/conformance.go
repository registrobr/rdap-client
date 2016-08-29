package protocol

// Conformance describes the RDAP conformance as it is in RFC 7483, section
// 4.1. The conformance is usually inserted in all responses to identify the
// extensions that the response includes
type Conformance struct {
	Levels []string `json:"rdapConformance,omitempty"`
}

// ConformanceSetter interface for identifying response objects that can
// contain the conformance structure
type ConformanceSetter interface {
	SetConformance([]string)
}

// SetConformance implements the ConformanceSetter and is used to enable
// conformance fields in response objects
func (l *Conformance) SetConformance(levels []string) {
	l.Levels = levels
}
