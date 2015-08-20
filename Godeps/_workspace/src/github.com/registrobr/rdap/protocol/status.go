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

// Proposed by NIC.br for DNS and DNSSEC checks of delegations
const (
	// StatusNSAA nameserver has authority for the domain (well configured)
	StatusNSAA Status = "ns aa"

	// StatusNSTimeout did not receive any answer of the nameserver when
	// performing a DNS query
	StatusNSTimeout Status = "ns timeout"

	// StatusNSNoAA nameserver doesn't have authority for the domain
	StatusNSNoAA Status = "ns noaa"

	// StatusNSUDN nameserver answers with unknown domain name
	StatusNSUDN Status = "ns udn"

	// StatusNSUH nameserver name could not be resolved
	StatusNSUH Status = "ns uh"

	// StatusNSFail nameserver answers with an internal server error
	StatusNSFail Status = "ns fail"

	// StatusNSQueryRefused nameserver refused to give an answer
	StatusNSQueryRefused Status = "ns query refused"

	// StatusNSConnectionRefused connection was refused (firewall)
	StatusNSConnectionRefused Status = "ns connection refused"

	// StatusNSError some generic error occurred while checking the nameserver
	StatusNSError Status = "ns error"

	// StatusNSCNAME found CNAME record in zone APEX (RFC 2181, section 10.1)
	StatusNSCNAME Status = "ns cname"

	// StatusNSSOAVersion found different SOA versions between the nameservers
	StatusNSSOAVersion Status = "ns soaVersion"

	// StatusDSOK all nameservers are well configured for this DS record
	StatusDSOK Status = "ds ok"

	// StatusDSTimeout did not receive any answer of the nameserver when
	// performing a DNSSEC query
	StatusDSTimeout Status = "ds timeout"

	// StatusDSNoSig no signature (RRSIG) was found in the answer
	StatusDSNoSig Status = "ds nosig"

	// StatusDSExpiredSig signature (RRSIG) is expired
	StatusDSExpiredSig Status = "ds expiredsig"

	// StatusDSInvalidSig signature (RRSIG) is invalid when checked with the
	// public key (DNSKEY)
	StatusDSInvalidSig Status = "ds invalidsig"

	// StatusDSNotFound the corresponding public key (DNSKEY) was not found in
	// the keyset
	StatusDSNotFound Status = "ds notfound"

	// StatusDSNoSEP the corresponding public key (DNSKEY) isn't a secure entry
	// point
	StatusDSNoSEP Status = "ds nosep"
)

// Proposed by NIC.br for domain status
const (
	// StatusWaitingActivation waiting for the next publication cycle to
	// publish the domain name in the DNS
	StatusWaitingActivation Status = "nicbr waiting activation"

	// StatusWaitingInactivation waiting for the next publication cycle to
	// remove the domain name from the DNS
	StatusWaitingInactivation Status = "nicbr waiting inactivation"

	// StatusInactiveCourtOrder legal decision was executed and the domain name
	// cannot be published in the DNS
	StatusInactiveCourtOrder Status = "nicbr inactive court order"

	// StatusInactiveCG by a decision of the "Comitê Gestor da Internet no
	// Brasil", this domain name cannot be published in the DNS
	StatusInactiveCG Status = "nicbr inactive CG"
)

// Status stores one of the possible status as listed in RFC 7483, section
// 10.2.2
type Status string
