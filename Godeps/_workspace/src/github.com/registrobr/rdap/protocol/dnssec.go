package protocol

import "time"

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

type SecureDNS struct {
	ZoneSigned       bool `json:"zoneSigned"`
	DelegationSigned bool `json:"delegationSigned"`
	DSData           []DS `json:"dsData,omitempty"`
}
