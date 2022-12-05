package run

import (
	"github.com/nextmv-io/sdk/run/decode"
	"github.com/nextmv-io/sdk/store"
)

// RunnerOption configures a Runner.
type RunnerOption[Input, Option, Solution any] func(
	Runner[Input, Option, Solution],
)

// InputDecode sets the decoder of a runner using f.
func InputDecode[Input, Option, Solution any](d InputDecoder[Input]) func(
	Runner[Input, Option, Solution],
) {
	return func(r Runner[Input, Option, Solution]) { r.SetInputDecoder(d) }
}

// Decode sets the decoder of a runner using f.
func Decode[Input any, Decoder decode.Decoder](
	d Decoder,
) func(
	Runner[Input, store.Options, store.Solution],
) {
	return InputDecode[Input, store.Options, store.Solution](
		NewGenericDecoder[Input](d),
	)
}
