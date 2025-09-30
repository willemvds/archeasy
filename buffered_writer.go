package archeasy

import (
	"io"
)

type BufferedWriter interface {
	io.Writer
	Flush() error
}
