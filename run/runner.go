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
	SetEncoder(Encoder[Solution])
}
