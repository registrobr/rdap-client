package protocol

import (
	"fmt"
	"net/http"
	"strings"
)

// Error describes an Error Response Body as it is in RFC 7483, section 6
type Error struct {
	Notices     []Notice `json:"notices,omitempty"`
	Lang        string   `json:"lang,omitempty"`
	ErrorCode   int      `json:"errorCode,omitempty"`
	Title       string   `json:"title,omitempty"`
	Description []string `json:"description,omitempty"`
	Conformance
	Port43
}

// Error make it easy to transport the protocol error via Go error interface
// between different API levels
func (e Error) Error() string {
	return fmt.Sprintf("HTTP status code: %d (%s)\n%s:\n  %s",
		e.ErrorCode,
		http.StatusText(e.ErrorCode),
		e.Title,
		strings.Join(e.Description, ", "))
}
