package output

import "io"

type Printer interface {
	Print(io.Writer) error
}
