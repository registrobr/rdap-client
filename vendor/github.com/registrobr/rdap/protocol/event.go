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

	// EventActionLastUpdate last date and time the database used by the RDAP
	// service was updated from the Registry or Registrar database
	EventActionLastUpdate EventAction = "last update"

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
	Date   EventDate   `json:"eventDate"`

	// Status was proposed by NIC.br to store the status of a current event.
	// For NIC.br specific use was useful to store the status of a delegation
	// check event
	Status []Status `json:"status,omitempty"`
}

// EventDate stores a Go time type and uses a more flexible algorithm for
// parsing the date, to allow a partial RFC 3339 format
type EventDate struct {
	time.Time
}

// Date returns the EventDate corresponding to
//	yyyy-mm-dd hh:mm:ss + nsec nanoseconds
// in the appropriate zone for that time in the given location.
//
// The month, day, hour, min, sec, and nsec values may be outside
// their usual ranges and will be normalized during the conversion.
// For example, October 32 converts to November 1.
//
// A daylight savings time transition skips or repeats times.
// For example, in the United States, March 13, 2011 2:15am never occurred,
// while November 6, 2011 1:15am occurred twice.  In such cases, the
// choice of time zone, and therefore the time, is not well-defined.
// Date returns a time that is correct in one of the two zones involved
// in the transition, but it does not guarantee which.
//
// Date panics if loc is nil.
func Date(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) EventDate {
	return EventDate{
		Time: time.Date(year, month, day, hour, min, sec, nsec, loc),
	}
}

// NewEventDate creates the object with the informed time
func NewEventDate(t time.Time) EventDate {
	return EventDate{
		Time: t,
	}
}

// UnmarshalJSON implements the json.Unmarshaler interface. The time is
// expected to be a quoted string in RFC 3339 format with or without the
// time/timezone
func (e *EventDate) UnmarshalJSON(data []byte) (err error) {
	if err = e.Time.UnmarshalJSON(data); err == nil {
		return
	}

	// allow date without time
	if e.Time, err = time.Parse(`"2006-01-02"`, string(data)); err == nil {
		return
	}

	// allow date without timezone
	e.Time, err = time.Parse(`"2006-01-02T15:04:05"`, string(data))
	return
}

// UnmarshalText implements the encoding.TextUnmarshaler interface. The time
// is expected to be in RFC 3339 format with or without the time/timezone
func (e *EventDate) UnmarshalText(data []byte) (err error) {
	if err = e.Time.UnmarshalText(data); err == nil {
		return
	}

	// allow date without time
	if e.Time, err = time.Parse(`2006-01-02`, string(data)); err == nil {
		return
	}

	// allow date without timezone
	e.Time, err = time.Parse(`2006-01-02T15:04:05`, string(data))
	return
}
