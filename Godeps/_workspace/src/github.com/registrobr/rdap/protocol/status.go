package protocol

const (
	StatusActive        Status = "active"
	StatusInactive      Status = "inactive"
	StatusPendingCreate Status = "pending create"
)

type Status string
