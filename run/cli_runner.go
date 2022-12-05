package run

import (
	"github.com/nextmv-io/sdk/run/decode"
	"github.com/nextmv-io/sdk/run/encode"
)

// CliRunner is the default CLI runner.
func CliRunner[Input, Option, Solution any](
	handler Algorithm[Input, Option, Solution],
	options ...RunnerOption[Input, Option, Solution],
) Runner[Input, Option, Solution] {
	runner := &genericRunner[Input, Option, Solution]{
		IOProducer:    CliIOProducer,
		InputDecoder:  NewGenericDecoder[Input](decode.JSON()),
		OptionDecoder: NoopOptionsDecoder[Option],
		Algorithm:     handler,
		Encoder:       GenericEncoder[Solution, Option, encode.JSONEncoder],
	}

	for _, option := range options {
		option(runner)
	}

	runnerConfig, decodedOption, err := FlagParser[
		Option, CliRunnerConfig,
	]()
	runner.runnerConfig = runnerConfig
	runner.decodedOption = decodedOption
	if err != nil {
		panic(err)
	}
	return runner
}
