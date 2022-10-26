package encode

import (
	"encoding/gob"
	"io"
)

func Gob() Encoder {
	return gobEncoder{}
}

type gobEncoder struct{}

func (g gobEncoder) Encode(w io.Writer, v any) error {
	return gob.NewEncoder(w).Encode(v)
}
