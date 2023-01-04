// Package http provides additional features for HTTPRunner.
package http

import (
	"github.com/nextmv-io/sdk/run"
	"github.com/nextmv-io/sdk/run/decode"
	"github.com/nextmv-io/sdk/run/encode"
)

// Decode sets the decoder of a HTTPRunner.
func Decode[Input, Option, Solution any, Decoder decode.Decoder](
	d Decoder,
) func(
	run.Runner[run.HTTPRunnerConfig, Input, Option, Solution],
) {
	return run.InputDecode[run.HTTPRunnerConfig, Input, Option, Solution](
		run.GenericDecoder[Input](d),
	)
}

// Encode sets the encoder of a HTTPRunner.
func Encode[Input, Option, Solution any, Encoder encode.Encoder](
	e Encoder,
) func(
	run.Runner[run.HTTPRunnerConfig, Input, Option, Solution],
) {
	return run.Encode[run.HTTPRunnerConfig, Input](
		run.GenericEncoder[Solution, Option](e),
	)
}
