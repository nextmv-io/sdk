package run

import (
	"github.com/nextmv-io/sdk/run/decode"
	"github.com/nextmv-io/sdk/run/encode"
)

// NewCLIRunner is the default CLI runner.
func NewCLIRunner[Input, Option, Solution any](
	algorithm Algorithm[Input, Option, Solution],
	options ...RunnerOption[Input, Option, Solution],
) Runner[Input, Option, Solution] {
	runner := &genericRunner[Input, Option, Solution]{
		IOProducer:    CliIOProducer,
		InputDecoder:  GenericDecoder[Input](decode.JSON()),
		OptionDecoder: NoopOptionsDecoder[Option],
		Algorithm:     algorithm,
		Encoder:       GenericEncoder[Solution, Option](encode.JSON()),
	}

	runnerConfig, decodedOption, err := FlagParser[
		Option, CLIRunnerConfig,
	]()
	runner.runnerConfig = runnerConfig
	runner.decodedOption = decodedOption
	if err != nil {
		panic(err)
	}

	for _, option := range options {
		option(runner)
	}

	return runner
}
