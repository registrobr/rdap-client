package protocol

import "time"

// EventAction can store all different types of event actions described in
// RFC 7483, section 10.2.3
type EventAction string

// https://tools.ietf.org/html/rfc7483#section-10.2.3
const (
	// EventActionRegistration the object instance was initially registered
	EventActionRegistration EventAction = "registration"

	// EventActionReRegistration the object instance was registered
	// subsequently to initial registration
	EventActionReRegistration EventAction = "reregistration"

	// EventActionLastChanged action noting when the information in the
	// object instance was last changed
	EventActionLastChanged EventAction = "last changed"

	// EventActionExpiration the object instance has been removed or will be
	// removed at a predetermined date and time from the registry
	EventActionExpiration EventAction = "expiration"

	// EventActionDeletion the object instance was removed from the registry
	// at a point in time that was not predetermined
	EventActionDeletion EventAction = "deletion"

	// EventActionTransfer the object instance was transferred from one
	// registrant to another
	EventActionTransfer EventAction = "transfer"

	// EventActionLocked the object instance was locked (see the "locked"
	// status)
	EventActionLocked EventAction = "locked"

	// EventActionUnlocked the object instance was unlocked (see the "locked"
	// status)
	EventActionUnlocked EventAction = "unlocked"

	// EventDelegationCheck was proposed by NIC.br to store information about
	// DNS checks performed by the registry
	EventDelegationCheck EventAction = "delegation check"

	// EventDelegationSignCheck was proposed by NIC.br to store information about
	// DNSSEC checks performed by the registry
	EventDelegationSignCheck EventAction = "delegation sign check"

	// EventLastCorrectDelegationCheck was proposed by NIC.br to store the date
	// of the last time that the nameserver was well configured
	EventLastCorrectDelegationCheck EventAction = "last correct delegation check"

	// EventLastCorrectDelegationSignCheck was proposed by NIC.br to store the date
	// of the last time that the nameservers were well configured with DNSSEC
	// for the related DS record
	EventLastCorrectDelegationSignCheck EventAction = "last correct delegation sign check"
)

// Event describes Events as it is in RFC 7483, section 4.5
type Event struct {
	Action EventAction `json:"eventAction"`
	Actor  string      `json:"eventActor,omitempty"`
	Date   time.Time   `json:"eventDate"`

	// Status was proposed by NIC.br to store the status of a current event.
	// For NIC.br specific use was useful to store the status of a delegation
	// check event
	Status []Status `json:"status,omitempty"`
}
