package protocol

// https://tools.ietf.org/html/rfc7483#section-10.2.2
const (
	// StatusActive the object instance is in use.  For domain names, it
	// signifies that the domain name is published in DNS.  For network and autnum
	// registrations, it signifies that they are allocated or assigned for use in
	// operational networks.  This maps to the "OK" status of the Extensible
	// Provisioning Protocol (EPP) [RFC5730]
	StatusActive Status = "active"

	// StatusInactive the object instance is not in use. See "active"
	StatusInactive Status = "inactive"

	// StatusPendingCreate a request has been received for the creation of the
	// object instance, but this action is not yet complete
	StatusPendingCreate Status = "pending create"

	// StatusRemoved some of the information of the object instance has not
	// been made available and has been removed. This is most commonly applied
	// to entities
	StatusRemoved Status = "removed"
)

// Status stores one of the possible status as listed in RFC 7483, section
// 10.2.2
type Status string
