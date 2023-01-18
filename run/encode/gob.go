package encode

import (
	"encoding/gob"
	"io"
)

// Gob returns a new binary value encoder.
func Gob() Encoder {
	return GobEncoder{}
}

// GobEncoder is a Encoder that encodes a struct into a gob.
type GobEncoder struct{}

// Encode writes the binary encoding of v to the w stream,
// followed by a newline character.
func (g GobEncoder) Encode(w io.Writer, v any) error {
	return gob.NewEncoder(w).Encode(v)
}

// ContentType returns the content type of the encoder.
func (g GobEncoder) ContentType() string {
	return "application/gob"
}
