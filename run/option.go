package run

import (
	"github.com/nextmv-io/sdk/run/decode"
	"github.com/nextmv-io/sdk/run/encode"
)

// An Option configures a Runner.
type Option func(Runner)

// Decode sets the decoder of a runner using f.
func Decode(f func() decode.Decoder) func(Runner) {
	return func(r Runner) { r.SetDecoder(f()) }
}

// Encode sets the encoder of a runner using f.
func Encode(f func() encode.Encoder) func(Runner) {
	return func(r Runner) { r.SetEncoder(f()) }
}
