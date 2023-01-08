package decode

import (
	"encoding/xml"
	"io"
)

// XML creates a JSON decoder.
func XML() Decoder {
	return XMLDecoder{}
}

// XMLDecoder is a Decoder that decodes a json into a struct.
type XMLDecoder struct{}

// Decode decodes JSON to the data structure v.
func (j XMLDecoder) Decode(r io.Reader, v any) error {
	return xml.NewDecoder(r).Decode(&v)
}
