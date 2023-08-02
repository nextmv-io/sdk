// Package cli provides additional features for CLIRunner.
package cli

import (
	"github.com/nextmv-io/sdk/run"
	"github.com/nextmv-io/sdk/run/decode"
	"github.com/nextmv-io/sdk/run/encode"
)

// Decode sets the decoder of a CLIRunner.
func Decode[Input, Option, Solution any, Decoder decode.Decoder](
	d Decoder,
) func(
	run.Runner[run.CLIRunnerConfig, Input, Option, Solution],
) {
	return run.InputDecode[run.CLIRunnerConfig, Input, Option, Solution](
		run.GenericDecoder[Input](d),
	)
}

// Encode sets the encoder of a CLIRunner.
func Encode[Input, Option, Solution any, Encoder encode.Encoder](
	e Encoder,
) func(
	run.Runner[run.CLIRunnerConfig, Input, Option, Solution],
) {
	return run.Encode[run.CLIRunnerConfig, Input](
		run.GenericEncoder[Solution, Option](e),
	)
}

// Validate sets the validator of a CLIRunner.
func Validate[Input, Option, Solution any](
	v run.Validator[Input],
) func(
	run.Runner[run.CLIRunnerConfig, Input, Option, Solution],
) {
	return run.InputValidate[run.CLIRunnerConfig, Input, Option, Solution](v)
}
