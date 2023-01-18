package decode

import (
	"encoding/json"
	"io"
)

// JSON creates a JSON decoder.
func JSON() Decoder {
	return JSONDecoder{}
}

// JSONDecoder is a Decoder that decodes a json into a struct.
type JSONDecoder struct{}

// Decode decodes JSON to the data structure v.
func (j JSONDecoder) Decode(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(&v)
}
