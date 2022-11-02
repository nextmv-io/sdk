package encode

import (
	"encoding/json"
	"io"
)

// JSON returns a new encoder that writes JSON.
func JSON() Encoder {
	return jsonEncoder{}
}

type jsonEncoder struct{}

// Encode writes the JSON encoding of v to the w stream,
// followed by a newline character.
func (j jsonEncoder) Encode(w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(v)
}
