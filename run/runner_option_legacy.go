package run

import (
	"github.com/nextmv-io/sdk/run/decode"
	"github.com/nextmv-io/sdk/store"
)

// Decode sets the decoder of a runner. This is a legacy function. Alternatively
// use InputDecode.
func Decode[RunnerConfig, Input any, Decoder decode.Decoder](
	d Decoder,
) func(
	Runner[RunnerConfig, Input, store.Options, store.Solution],
) {
	return InputDecode[RunnerConfig, Input, store.Options, store.Solution](
		GenericDecoder[Input](d),
	)
}
