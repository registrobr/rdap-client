package protocol

// Port43 described as in RFC 7483, section 4.7. The port43 is usually inserted in all responses to
// identify the WHOIS server where the containing object instance may be found
type Port43 struct {
	Port43 string `json:"port43,omitempty"`
}

// Port43Setter interface for identifying response objects that can
// contain the conformance structure
type Port43Setter interface {
	SetPort43(string)
}

// SetPort43 implements the Port43Setter and is used to fill port43 field in response objects
func (l *Port43) SetPort43(port43 string) {
	l.Port43 = port43
}
