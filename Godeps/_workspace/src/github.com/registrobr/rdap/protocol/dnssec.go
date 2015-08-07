package protocol

import "time"

// DS describes the dsData as it is in RFC 7483, section 5.3
type DS struct {
	KeyTag      int       `json:"keyTag"`
	Algorithm   int       `json:"algorithm"`
	Digest      string    `json:"digest"`
	DigestType  int       `json:"digestType"`
	Events      []Event   `json:"events,omitempty"`
	DSStatus    string    `json:"nicbr_status,omitempty"`
	LastCheckAt time.Time `json:"nicbr_lastCheck,omitempty"`
	LastOKAt    time.Time `json:"nicbr_lastOK,omitempty"`
}

// SecureDNS describes the secureDNS as it is in RFC 7483, section 5.3
type SecureDNS struct {
	ZoneSigned       bool `json:"zoneSigned"`
	DelegationSigned bool `json:"delegationSigned"`
	DSData           []DS `json:"dsData,omitempty"`
}
