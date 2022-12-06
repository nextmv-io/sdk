package run

import (
	"github.com/nextmv-io/sdk/run/decode"
	"github.com/nextmv-io/sdk/run/encode"
)

// CliRunner is the default CLI runner.
func CliRunner[Input, Option, Solution any](
	algorithm Algorithm[Input, Option, Solution],
	options ...RunnerOption[Input, Option, Solution],
) Runner[Input, Option, Solution] {
	runner := &genericRunner[Input, Option, Solution]{
		IOProducer:    CliIOProducer,
		InputDecoder:  NewGenericDecoder[Input](decode.JSON()),
		OptionDecoder: NoopOptionsDecoder[Option],
		Algorithm:     algorithm,
		Encoder:       NewGenericEncoder[Solution, Option](encode.JSON()),
	}

	runnerConfig, decodedOption, err := FlagParser[
		Option, CliRunnerConfig,
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
