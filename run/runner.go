package run

import "context"

// Runner defines the interface of the runner.
type Runner[RunnerConfig, Input, Option, Solution any] interface {
	// Run runs the runner.
	Run(context.Context) error
	// SetIOProducer sets the ioProducer of a runner.
	SetIOProducer(IOProducer[RunnerConfig])
	// SetInputDecoder sets the inputDecoder of a runner.
	SetInputDecoder(Decoder[Input])
	// SetOptionDecoder sets the optionDecoder of a runner.
	SetOptionDecoder(Decoder[Option])
	// SetAlgorithm sets the algorithm of a runner.
	SetAlgorithm(Algorithm[Input, Option, Solution])
	// SetEncoder sets the encoder of a runner.
	SetEncoder(Encoder[Solution, Option])
	// GetEncoder returns the encoder of a runner.
	GetEncoder() Encoder[Solution, Option]
	// RunnerConfig returns the runnerConfig of a runner.
	RunnerConfig() RunnerConfig
}

// IOProducer is a function that produces the input, option and writer.
type IOProducer[RunnerConfig any] func(context.Context, RunnerConfig) IOData

// Decoder is a function that decodes a reader into a struct.
type Decoder[Input any] func(context.Context, any) (Input, error)

// Algorithm is a function that runs an algorithm.
type Algorithm[Input, Option, Solution any] func(
	context.Context, Input, Option, chan<- Solution,
) error

// Encoder is an interface which defines the Encode function that encodes a
// struct into a writer.
type Encoder[Solution, Option any] interface {
	Encode(context.Context, <-chan Solution, any, any, Option) error
}

// ContentTyper is an interface which defines the ContentType function that
// describes the content type of the encoder.
type ContentTyper interface {
	ContentType() string
}
