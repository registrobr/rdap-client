package protocol

import "time"

type EventAction string

const (
	EventActionRegistration   EventAction = "registration"
	EventActionReRegistration EventAction = "reregistration"
	EventActionLastChanged    EventAction = "last changed"
	EventActionExpiration     EventAction = "expiration"
	EventActionDeletion       EventAction = "deletion"
	EventActionTransfer       EventAction = "transfer"
	EventActionLocked         EventAction = "locked"
	EventActionUnlocked       EventAction = "unlocked"
)

type Event struct {
	Action EventAction `json:"eventAction"`
	Actor  string      `json:"eventActor,omitempty"`
	Date   time.Time   `json:"eventDate"`
}
