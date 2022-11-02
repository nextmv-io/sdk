// Package decode provides decoders for output of a runner.
package decode

import "io"

// Decoder defines a decoder.
type Decoder interface {
	// Decode decodes the data in reader to structure in the second
	// argument of the method.
	Decode(io.Reader, any) error
}
