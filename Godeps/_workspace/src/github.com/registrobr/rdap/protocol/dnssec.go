package protocol

type DS struct {
	KeyTag     int     `json:"keyTag"`
	Algorithm  int     `json:"algorithm"`
	Digest     string  `json:"digest"`
	DigestType int     `json:"digestType"`
	Events     []Event `json:"events,omitempty"`
	DSStatus   string  `json:"nicbr_status,omitempty"`
}

type SecureDNS struct {
	ZoneSigned        bool `json:"zoneSigned"`
	DelegationsSigned bool `json:"delegationsSigned"`
	DSData            []DS `json:"dsData"`
}
