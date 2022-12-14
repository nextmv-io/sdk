package run

import (
	"github.com/nextmv-io/sdk/run/decode"
	"github.com/nextmv-io/sdk/run/encode"
)

// CLIRunner is the default CLI runner.
func CLIRunner[Input, Option, Solution any](
	algorithm Algorithm[Input, Option, Solution],
	options ...RunnerOption[CLIRunnerConfig, Input, Option, Solution],
) Runner[CLIRunnerConfig, Input, Option, Solution] {
	runner := GenericRunner(
		CliIOProducer,
		GenericDecoder[Input](decode.JSON()),
		NoopOptionsDecoder[Option],
		algorithm,
		GenericEncoder[Solution, Option](encode.JSON()),
	)

	for _, option := range options {
		option(runner)
	}

	return runner
}
