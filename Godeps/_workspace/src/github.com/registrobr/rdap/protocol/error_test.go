package protocol

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestErrorError(t *testing.T) {
	err := Error{
		ErrorCode:   http.StatusBadRequest,
		Title:       "Invalid FQDN",
		Description: []string{"The informed FQDN has an invalid format according to the RFCs"},
	}

	expected := fmt.Sprintf("HTTP status code: %d (%s)\n%s:\n  %s", err.ErrorCode, http.StatusText(err.ErrorCode), err.Title, strings.Join(err.Description, ", "))

	if err.Error() != expected {
		t.Errorf("Unexpected error message. Expected “%s” and got “%s”", expected, err.Error())
	}
}
