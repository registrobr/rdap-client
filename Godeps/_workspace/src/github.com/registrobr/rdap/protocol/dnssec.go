package protocol

// DS describes the dsData as it is in RFC 7483, section 5.3
type DS struct {
	KeyTag     int     `json:"keyTag"`
	Algorithm  int     `json:"algorithm"`
	Digest     string  `json:"digest"`
	DigestType int     `json:"digestType"`
	Events     []Event `json:"events,omitempty"`
}

// SecureDNS describes the secureDNS as it is in RFC 7483, section 5.3
type SecureDNS struct {
	// ZoneSigned does not make too much sense for us to use
	// it, so we need to use a pointer to hide it with omitempty. Maybe the
	// real use for it is for TLDs that publish the DS records
	// without signing it, but its not clear in RFC 7483, section 5.3
	ZoneSigned       *bool `json:"zoneSigned,omitempty"`
	DelegationSigned bool  `json:"delegationSigned"`
	DSData           []DS  `json:"dsData,omitempty"`
}
