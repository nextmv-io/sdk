package run

import (
	"context"

	"github.com/nextmv-io/sdk/run/decode"
	"github.com/nextmv-io/sdk/run/encode"
)

// CliRunner is the default CLI runner.
func CliRunner[Input, Option, Solution any](
	handler Algorithm[Input, Option, Solution],
) Runner[Input, Option, Solution] {
	runner := &genericRunner[Input, Option, Solution]{
		IOProducer:    DefaultIOProducer,
		InputDecoder:  GenericDecoder[Input, decode.JSONDecoder],
		OptionDecoder: NoopOptionsDecoder[Option],
		Algorithm:     handler,
		Encoder:       GenericEncoder[Solution, Option, encode.JSONEncoder],
	}
	runnerConfig, decodedOption, err := DefaultFlagParser[
		Option, CliRunnerConfig,
	]()
	runner.runnerConfig = runnerConfig
	runner.decodedOption = decodedOption
	if err != nil {
		panic(err)
	}
	return runner
}

// NewGenericRunner creates a new one-off runner.
func NewGenericRunner[Input, Option, Solution any](
	ioHandler IOProducer,
	inputDecoder InputDecoder[Input],
	optionDecoder OptionDecoder[Option],
	handler Algorithm[Input, Option, Solution],
	encoder Encoder[Solution, Option],
) Runner[Input, Option, Solution] {
	return &genericRunner[Input, Option, Solution]{
		IOProducer:    ioHandler,
		InputDecoder:  inputDecoder,
		OptionDecoder: optionDecoder,
		Algorithm:     handler,
		Encoder:       encoder,
	}
}

type genericRunner[Input, Option, Solution any] struct {
	IOProducer    IOProducer
	InputDecoder  InputDecoder[Input]
	OptionDecoder OptionDecoder[Option]
	Algorithm     Algorithm[Input, Option, Solution]
	Encoder       Encoder[Solution, Option]
	runnerConfig  any
	decodedOption Option
}

func (r *genericRunner[Input, Option, Solution]) Run(
	context context.Context,
) error {
	// get IO
	ioData := r.IOProducer(context, r.runnerConfig)

	// decode input
	decodedInput, err := r.InputDecoder(context, ioData.Input())
	if err != nil {
		return err
	}

	// decode option
	r.decodedOption, err = r.OptionDecoder(
		context, ioData.Option(), r.decodedOption,
	)
	if err != nil {
		return err
	}

	// run algorithm
	solutions := make(chan Solution)
	errs := make(chan error, 1)
	go func() {
		defer close(solutions)
		defer close(errs)
		err = r.Algorithm(context, decodedInput, r.decodedOption, solutions)
		if err != nil {
			errs <- err
			return
		}
	}()

	// encode solutions
	err = r.Encoder(
		context, solutions, ioData.Writer(), r.runnerConfig, r.decodedOption,
	)
	if err != nil {
		return err
	}

	// return potential errors
	return <-errs
}

func (r *genericRunner[Input, Option, Solution]) SetIOProducer(
	ioProducer IOProducer,
) {
	r.IOProducer = ioProducer
}

func (r *genericRunner[Input, Option, Solution]) SetInputDecoder(
	decoder InputDecoder[Input],
) {
	r.InputDecoder = decoder
}

func (r *genericRunner[Input, Option, Solution]) SetOptionDecoder(
	decoder OptionDecoder[Option],
) {
	r.OptionDecoder = decoder
}

func (r *genericRunner[Input, Option, Solution]) SetAlgorithm(
	algorithm Algorithm[Input, Option, Solution],
) {
	r.Algorithm = algorithm
}

func (r *genericRunner[Input, Option, Solution]) SetEncoder(
	encoder Encoder[Solution, Option],
) {
	r.Encoder = encoder
}
