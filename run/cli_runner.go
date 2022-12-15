package run

import (
	"github.com/nextmv-io/sdk/run/decode"
	"github.com/nextmv-io/sdk/run/encode"
)

// CLIRunner is the default CLI runner. It reads the input from stdin or a file,
// writes output to stdout or a file, decodes the input using the JSON decoder,
// accepts options from the command line, and encodes the solution using the
// JSON encoder.
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
