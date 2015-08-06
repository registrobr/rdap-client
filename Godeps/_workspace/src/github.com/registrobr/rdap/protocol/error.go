package protocol

import (
	"fmt"
	"net/http"
	"strings"
)

type Error struct {
	Notices     []Notice `json:"notices,omitempty"`
	Lang        string   `json:"lang,omitempty"`
	ErrorCode   int      `json:"errorCode,omitempty"`
	Title       string   `json:"title,omitempty"`
	Description []string `json:"description,omitempty"`
	Conformance
}

func (e Error) Error() string {
	return fmt.Sprintf("HTTP status code: %d (%s)\n%s:\n  %s",
		e.ErrorCode,
		http.StatusText(e.ErrorCode),
		e.Title,
		strings.Join(e.Description, ", "))
}
