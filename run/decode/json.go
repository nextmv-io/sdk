package decode

import (
	"encoding/json"
	"io"
)

func JSON() Decoder {
	return jsonDecoder{}
}

type jsonDecoder struct{}

func (j jsonDecoder) Decode(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(&v)
}
