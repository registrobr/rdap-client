package protocol

const (
	StatusActive        Status = "active"
	StatusInactive      Status = "inactive"
	StatusPendingCreate Status = "pending create"
	StatusRemoved       Status = "removed"
)

type Status string
