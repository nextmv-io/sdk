package run

import (
	"context"
)

// DefaultOneOffRunner is the default one-off runner.
func DefaultOneOffRunner[Input, Option, Solution any](
	handler Algorithm[Input, Option, Solution],
) Runner[Input, Option, Solution] {
	return NewOneOffRunner(
		DefaultIOProducer,
		JSONDecoder[Input],
		OptionsDecoder[Option],
		handler,
		JSONEncoder[Solution],
	)
}

// NewOneOffRunner creates a new one-off runner.
func NewOneOffRunner[Input, Option, Solution any](
	ioHandler IOProducer,
	inputDecoder InputDecoder[Input],
	optionDecoder OptionDecoder[Option],
	handler Algorithm[Input, Option, Solution],
	encoder Encoder[Solution],
) Runner[Input, Option, Solution] {
	return &oneOffRunner[Input, Option, Solution]{
		IOProducer:    ioHandler,
		InputDecoder:  inputDecoder,
		OptionDecoder: optionDecoder,
		Algorithm:     handler,
		Encoder:       encoder,
		FlagParser:    DefaultFlagParser[Option],
	}
}

type oneOffRunner[Input, Option, Solution any] struct {
	FlagParser    FlagParser[Option]
	IOProducer    IOProducer
	InputDecoder  InputDecoder[Input]
	OptionDecoder OptionDecoder[Option]
	Algorithm     Algorithm[Input, Option, Solution]
	Encoder       Encoder[Solution]
}

func (r *oneOffRunner[Input, Option, Solution]) Run(
	context context.Context,
) error {
	runnerConfig, decodedOption, err := r.FlagParser()
	if err != nil {
		return err
	}
	ioData := r.IOProducer(context, runnerConfig)

	decodedInput, err := r.InputDecoder(context, ioData.Input())
	if err != nil {
		return err
	}
	decodedOption, err = r.OptionDecoder(context, ioData.Option(), decodedOption)
	if err != nil {
		return err
	}
	solutions, err := r.Algorithm(context, decodedInput, decodedOption)
	if err != nil {
		return err
	}
	err = r.Encoder(context, solutions, ioData.Writer())
	if err != nil {
		return err
	}
	return nil
}

func (r *oneOffRunner[Input, Option, Solution]) SetIOProducer(
	ioProducer IOProducer,
) {
	r.IOProducer = ioProducer
}

func (r *oneOffRunner[Input, Option, Solution]) SetInputDecoder(
	decoder InputDecoder[Input],
) {
	r.InputDecoder = decoder
}

func (r *oneOffRunner[Input, Option, Solution]) SetOptionDecoder(
	decoder OptionDecoder[Option],
) {
	r.OptionDecoder = decoder
}

func (r *oneOffRunner[Input, Option, Solution]) SetAlgorithm(
	algorithm Algorithm[Input, Option, Solution],
) {
	r.Algorithm = algorithm
}

func (r *oneOffRunner[Input, Option, Solution]) SetEncoder(
	encoder Encoder[Solution],
) {
	r.Encoder = encoder
}
