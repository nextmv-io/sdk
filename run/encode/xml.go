package encode

import (
	"encoding/xml"
	"io"
)

// XML returns a new binary value encoder.
func XML() Encoder {
	return XMLEncoder{}
}

// XMLEncoder is a Encoder that encodes a struct into a gob.
type XMLEncoder struct{}

// Encode writes the binary encoding of v to the w stream,
// followed by a newline character.
func (g XMLEncoder) Encode(w io.Writer, v any) error {
	return xml.NewEncoder(w).Encode(v)
}

// ContentType returns the content type of the encoder.
func (g XMLEncoder) ContentType() string {
	return "application/xml"
}
