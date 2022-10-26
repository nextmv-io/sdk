package run

import (
	"github.com/nextmv-io/sdk/run/decode"
	"github.com/nextmv-io/sdk/run/encode"
)

type Option func(Runner)

func Decode(f func() decode.Decoder) func(Runner) {
	return func(r Runner) { r.SetDecoder(f()) }
}

func Encode(f func() encode.Encoder) func(Runner) {
	return func(r Runner) { r.SetEncoder(f()) }
}
