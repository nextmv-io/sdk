// Package encode provides encoders for input of a runner.
package encode

import "io"

// Encoder defines an encoder.
type Encoder interface {
	// Encode encodes the data to the
	Encode(io.Writer, any) error
}
