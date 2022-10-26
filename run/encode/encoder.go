package encode

import "io"

type Encoder interface {
	Encode(io.Writer, any) error
}
