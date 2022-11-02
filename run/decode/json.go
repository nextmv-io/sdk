package decode

import (
	"encoding/json"
	"io"
)

// JSON creates a JSON decoder.
func JSON() Decoder {
	return jsonDecoder{}
}

type jsonDecoder struct{}

// Decode decodes JSON to the data structure v.
func (j jsonDecoder) Decode(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(&v)
}
