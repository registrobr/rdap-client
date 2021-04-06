package protocol

// Mar/2015 - https://tools.ietf.org/html/rfc7483#section-10.2.2
const (
	// Signifies that the data of the object instance has
	// been found to be accurate.  This type of status is usually
	// found on entity object instances to note the validity of
	// identifying contact information.
	StatusValidated Status = "validated"

	// The registration of the object instance has been
	// performed by a third party.  This is most commonly applied to
	// entities.
	StatusProxy Status = "proxy"

	// The information of the object instance is not
	// designated for public consumption.  This is most commonly
	// applied to entities.
	StatusPrivate Status = "private"

	// Some of the information of the object instance has
	// been altered for the purposes of not readily revealing the
	// actual information of the object instance.  This is most
	// commonly applied to entities.
	StatusObscured Status = "obscured"

	// The object instance is associated with other object
	// instances in the registry.  This is most commonly used to
	// signify that a nameserver is associated with a domain or that
	// an entity is associated with a network resource or domain.
	StatusAssociated Status = "associated"

	// Changes to the object instance cannot be made,
	// including the association of other object instances.
	StatusLocked Status = "locked"

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

	// StatusPendingRenew A request has been received for the renewal of the
	// object instance but this action is not yet complete.
	StatusPendingRenew Status = "pending renew"

	// StatusPendingtransfer A request has been received for the transfer of the
	// object instance but this action is not yet complete.
	StatusPendingTransfer Status = "pending transfer"

	// StatusPendingUpdate A request has been received for the update or
	// modification of the object instance but this action is not yet complete.
	StatusPendingUpdate Status = "pending update"

	// StatusPendingDelete A request has been received for the deletion or removal
	// of the object instance but this action is not yet complete. For domains,
	// this might mean that the name is no longer published in DNS but has not yet
	// been purged from the registry database.
	StatusPendingDelete Status = "pending delete"

	// StatusRenewProhibited Renewal or reregistration of the object instance is
	// forbidden.
	StatusRenewProhibited = "renew prohibited"

	// StatusTransferProhibited Transfers of the registration from one registrar
	// to another are forbidden. This type of status normally applies to DNR
	// domain names.
	StatusTransferProhibited = "transfer prohibited"

	// StatusUpdateProhibited Updates to the object instance are forbidden.
	StatusUpdateProhibited = "update prohibited"

	// StatusDeleteProhibited Deletion of the registration of the object instance
	// is forbidden. This type of status normally applies to DNR domain names.
	StatusDeleteProhibited = "delete prohibited"

	// StatusRemoved some of the information of the object instance has not
	// been made available and has been removed. This is most commonly applied
	// to entities
	StatusRemoved Status = "removed"
)

// The following values appeared at RFC8056 - Jan/2017 - https://tools.ietf.org/html/rfc8056#3.1
// EPP related constants
const (
	// This grace period is provided after the initial
	// registration of the object.  If the object is deleted by the
	// client during this period, the server provides a credit to the
	// client for the cost of the registration.  This maps to the Domain
	// Registry Grace Period Mapping for the Extensible Provisioning
	// Protocol (EPP) [RFC3915] 'addPeriod' status.
	StatusAddPeriod Status = "add period"

	// This grace period is provided after an object
	// registration period expires and is extended (renewed)
	// automatically by the server.  If the object is deleted by the
	// client during this period, the server provides a credit to the
	// client for the cost of the auto renewal.  This maps to the Domain
	// Registry Grace Period Mapping for the Extensible Provisioning
	// Protocol (EPP) [RFC3915] 'autoRenewPeriod' status.
	StatusAutoRenewPeriod Status = "auto renew period"

	// The client requested that requests to delete the
	// object MUST be rejected.  This maps to the Extensible Provisioning
	// Protocol (EPP) Domain Name Mapping [RFC5731], Extensible
	// Provisioning Protocol (EPP) Host Mapping [RFC5732], and Extensible
	// Provisioning Protocol (EPP) Contact Mapping [RFC5733]
	// 'clientDeleteProhibited' status.
	StatusClientDeleteProhibited Status = "client delete prohibited"

	// The client requested that the DNS delegation
	// information MUST NOT be published for the object.  This maps to
	// the Extensible Provisioning Protocol (EPP) Domain Name Mapping
	// [RFC5731] 'clientHold' status.
	StatusClientHold Status = "client hold"

	// The client requested that requests to renew the
	// object MUST be rejected.  This maps to the Extensible Provisioning
	// Protocol (EPP) Domain Name Mapping [RFC5731]
	// 'clientRenewProhibited' status.
	StatusClientRenewProhibited Status = "client renew prohibited"

	// The client requested that requests to transfer the
	// object MUST be rejected.  This maps to the Extensible Provisioning
	// Protocol (EPP) Domain Name Mapping [RFC5731] and Extensible
	// Provisioning Protocol (EPP) Contact Mapping [RFC5733]
	// 'clientTransferProhibited' status.
	StatusClientTransferProhibited Status = "client transfer prohibited"

	// The client requested that requests to update the
	// object (other than to remove this status) MUST be rejected.  This
	// maps to the Extensible Provisioning Protocol (EPP) Domain Name
	// Mapping [RFC5731], Extensible Provisioning Protocol (EPP) Host
	// Mapping [RFC5732], and Extensible Provisioning Protocol (EPP)
	// Contact Mapping [RFC5733] 'clientUpdateProhibited' status.
	StatusClientUpdateProhibited Status = "client update prohibited"

	// An object is in the process of being restored after
	// being in the redemption period state.  This maps to the Domain
	// Registry Grace Period Mapping for the Extensible Provisioning
	// Protocol (EPP) [RFC3915] 'pendingRestore' status.
	StatusPendingRestore Status = "pending restore"

	// A delete has been received, but the object has not
	// yet been purged because an opportunity exists to restore the
	// object and abort the deletion process.  This maps to the Domain
	// Registry Grace Period Mapping for the Extensible Provisioning
	// Protocol (EPP) [RFC3915] 'redemptionPeriod' status.
	StatusRedemptionPeriod Status = "redemption period"

	// This grace period is provided after an object
	// registration period is explicitly extended (renewed) by the
	// client.  If the object is deleted by the client during this
	// period, the server provides a credit to the client for the cost of
	// the renewal.  This maps to the Domain Registry Grace Period
	// Mapping for the Extensible Provisioning Protocol (EPP) [RFC3915]
	// 'renewPeriod' status.
	StatusRenewPeriod Status = "renew period"

	// The server set the status so that requests to delete
	// the object MUST be rejected.  This maps to the Extensible
	// Provisioning Protocol (EPP) Domain Name Mapping [RFC5731],
	// Extensible Provisioning Protocol (EPP) Host Mapping [RFC5732], and
	// Extensible Provisioning Protocol (EPP) Contact Mapping [RFC5733]
	// 'serverDeleteProhibited' status.
	StatusServerDeleteProhibited Status = "server delete prohibited"

	// The server set the status so that requests to renew
	// the object MUST be rejected.  This maps to the Extensible
	// Provisioning Protocol (EPP) Domain Name Mapping [RFC5731]
	// 'serverRenewProhibited' status.
	StatusServerRenewProhibited Status = "server renew prohibited"

	// The server set the status so that requests to
	// transfer the object MUST be rejected.  This maps to the Extensible
	// Provisioning Protocol (EPP) Domain Name Mapping [RFC5731] and
	// Extensible Provisioning Protocol (EPP) Contact Mapping [RFC5733]
	// 'serverTransferProhibited' status.
	StatusServerTransferProhibited Status = "server transfer prohibited"

	// The server set the status so that requests to update
	// the object (other than to remove this status) MUST be rejected.
	// This maps to the Extensible Provisioning Protocol (EPP) Domain
	// Name Mapping [RFC5731], Extensible Provisioning Protocol (EPP)
	// Host Mapping [RFC5732], and Extensible Provisioning Protocol (EPP)
	// Contact Mapping [RFC5733] 'serverUpdateProhibited' status.
	StatusServerUpdateProhibited Status = "server update prohibited"

	// The server set the status so that DNS delegation
	// information MUST NOT be published for the object.  This maps to
	// the Extensible Provisioning Protocol (EPP) Domain Name Mapping
	// [RFC5731] 'serverHold' status.
	StatusServerHold Status = "server hold"

	// This grace period is provided after the successful
	// transfer of object registration sponsorship from one client to
	// another client.  If the object is deleted by the client during
	// this period, the server provides a credit to the client for the
	// cost of the transfer.  This maps to the Domain Registry Grace
	// Period Mapping for the Extensible Provisioning Protocol (EPP)
	// [RFC3915] 'transferPeriod' status.
	StatusTransferPeriod Status = "transfer period"
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

	// StatusNONE dns error has occurred
	StatusNone Status = "plain dns error"
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
