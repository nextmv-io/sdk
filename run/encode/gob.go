package encode

import (
	"encoding/gob"
	"io"
)

// Gob returns a new binary value encoder.
func Gob() Encoder {
	return gobEncoder{}
}

type gobEncoder struct{}

// Encode writes the binary encoding of v to the w stream,
// followed by a newline character.
func (g gobEncoder) Encode(w io.Writer, v any) error {
	return gob.NewEncoder(w).Encode(v)
}
