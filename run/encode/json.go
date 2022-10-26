package encode

import (
	"encoding/json"
	"io"
)

func JSON() Encoder {
	return jsonEncoder{}
}

type jsonEncoder struct{}

func (j jsonEncoder) Encode(w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(v)
}
