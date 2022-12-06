package run

import (
	"github.com/nextmv-io/sdk/run/decode"
	"github.com/nextmv-io/sdk/store"
)

// RunnerOption configures a Runner.
type RunnerOption[Input, Option, Solution any] func(
	Runner[Input, Option, Solution],
)

// InputDecode sets the decoder of a runner.
func InputDecode[Input, Option, Solution any](i InputDecoder[Input]) func(
	Runner[Input, Option, Solution],
) {
	return func(r Runner[Input, Option, Solution]) { r.SetInputDecoder(i) }
}

// Decode sets the decoder of a runner. This is a legacy function. Alternatively
// use InputDecode.
func Decode[Input any, Decoder decode.Decoder](
	d Decoder,
) func(
	Runner[Input, store.Options, store.Solution],
) {
	return InputDecode[Input, store.Options, store.Solution](
		NewGenericDecoder[Input](d),
	)
}

// OptionDecode sets the options decoder of a runner.
func OptionDecode[Input, Option, Solution any](o OptionDecoder[Option]) func(
	Runner[Input, Option, Solution],
) {
	return func(r Runner[Input, Option, Solution]) { r.SetOptionDecoder(o) }
}

// Encode sets the encoder of a runner.
func Encode[Input, Option, Solution any](e Encoder[Solution, Option]) func(
	Runner[Input, Option, Solution],
) {
	return func(r Runner[Input, Option, Solution]) { r.SetEncoder(e) }
}

// IOProduce sets the ioProducer of a runner.
func IOProduce[Input, Option, Solution any](i IOProducer) func(
	Runner[Input, Option, Solution],
) {
	return func(r Runner[Input, Option, Solution]) { r.SetIOProducer(i) }
}
