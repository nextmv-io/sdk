package run

import (
	"github.com/nextmv-io/sdk/run/decode"
	"github.com/nextmv-io/sdk/store"
)

// RunnerOption configures a Runner.
type RunnerOption[RunnerConfig, Input, Option, Solution any] func(
	Runner[RunnerConfig, Input, Option, Solution],
)

// InputDecode sets the input decoder of a runner.
func InputDecode[
	RunnerConfig, Input, Option, Solution any,
](i Decoder[Input]) func(
	Runner[RunnerConfig, Input, Option, Solution],
) {
	return func(r Runner[RunnerConfig, Input, Option, Solution]) {
		r.SetInputDecoder(i)
	}
}

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

// OptionDecode sets the options decoder of a runner.
func OptionDecode[
	RunnerConfig, Input, Option, Solution any,
](o Decoder[Option]) func(
	Runner[RunnerConfig, Input, Option, Solution],
) {
	return func(r Runner[RunnerConfig, Input, Option, Solution]) {
		r.SetOptionDecoder(o)
	}
}

// Encode sets the encoder of a runner.
func Encode[
	RunnerConfig, Input, Option, Solution any,
](e Encoder[Solution, Option]) func(
	Runner[RunnerConfig, Input, Option, Solution],
) {
	return func(r Runner[RunnerConfig, Input, Option, Solution]) {
		r.SetEncoder(e)
	}
}

// IOProduce sets the IOProducer of a runner.
func IOProduce[
	RunnerConfig, Input, Option, Solution any,
](i IOProducer[RunnerConfig]) func(
	Runner[RunnerConfig, Input, Option, Solution],
) {
	return func(r Runner[RunnerConfig, Input, Option, Solution]) {
		r.SetIOProducer(i)
	}
}
