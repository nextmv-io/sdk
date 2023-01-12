package decode

import (
	"encoding/gob"
	"io"
)

// Gob returns a new binary value decoder.
func Gob() Decoder {
	return GobDecoder{}
}

// GobDecoder is a Decoder that decodes a gob into a struct.
type GobDecoder struct{}

// Decode decodes Gob to the data structure v.
func (g GobDecoder) Decode(r io.Reader, v any) error {
	return gob.NewDecoder(r).Decode(&v)
}
