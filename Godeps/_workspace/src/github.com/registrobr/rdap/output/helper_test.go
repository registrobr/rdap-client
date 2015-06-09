package output

import (
	"fmt"
	"strings"

	"github.com/registrobr/rdap/Godeps/_workspace/src/github.com/aryann/difflib"
)

func diff(a, b interface{}) []difflib.DiffRecord {
	return difflib.Diff(strings.Split(fmt.Sprintf("%v", a), "\n"),
		strings.Split(fmt.Sprintf("%v", b), "\n"))
}

type WriterMock struct {
	Content []byte
	Err     error

	MockWrite func(p []byte) (n int, err error)
}

func (w *WriterMock) Write(p []byte) (n int, err error) {
	if w.MockWrite != nil {
		return w.MockWrite(p)
	}

	w.Content = append(w.Content, p...)
	return len(p), w.Err
}
