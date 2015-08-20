package rdap

import (
	"io"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/aryann/difflib"
)

func diff(a, b interface{}) []difflib.DiffRecord {
	return difflib.Diff(
		strings.SplitAfter(spew.Sdump(a), "\n"),
		strings.SplitAfter(spew.Sdump(b), "\n"),
	)
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }
