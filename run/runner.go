package run

import "context"

// Runner defines the interface of the runner.
type Runner[Input, Option, Solution any] interface {
	// Run runs the runner.
	Run(context.Context) error
	// SetIOProducer sets the ioProducer of a runner.
	SetIOProducer(IOProducer)
	// SetInputDecoder sets the inputDecoder of a runner.
	SetInputDecoder(InputDecoder[Input])
	// SetOptionDecoder sets the optionDecoder of a runner.
	SetOptionDecoder(OptionDecoder[Option])
	// SetAlgorithm sets the algorithm of a runner.
	SetAlgorithm(Algorithm[Input, Option, Solution])
	// SetEncoder sets the encoder of a runner.
	SetEncoder(Encoder[Solution, Option])
}

// IOProducer is a function that produces the input, option and writer.
type IOProducer func(context.Context, any) IOData

// InputDecoder is a function that decodes a reader into a struct.
type InputDecoder[Input any] func(context.Context, any) (Input, error)

// OptionDecoder is a function that decodes a reader into a struct.
type OptionDecoder[Option any] func(
	context.Context, any, Option,
) (Option, error)

// Algorithm is a function that runs an algorithm.
type Algorithm[Input, Option, Solution any] func(
	context.Context, Input, Option, chan<- Solution,
) error

// Encoder is a function that encodes a struct into a writer.
type Encoder[Solution, Option any] func(
	context.Context, <-chan Solution, any, any, Option) error
