// Package cli provides additional features for CLI runners.
package cli

import (
	"github.com/nextmv-io/sdk/run"
	"github.com/nextmv-io/sdk/run/decode"
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
