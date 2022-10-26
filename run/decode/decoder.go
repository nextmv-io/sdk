package decode

import "io"

type Decoder interface {
	Decode(io.Reader, any) error
}
