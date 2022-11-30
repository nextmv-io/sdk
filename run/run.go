package run

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"strings"

	"github.com/itzg/go-flagsfiller"
	"github.com/nextmv-io/sdk/store"
)

// Run runs the runner.
func Run[Input, Option any](solver func(
	input Input, option Option,
) (store.Solver, error),
) error {
	algorithm := func(
		ctx context.Context,
		input Input, option Option, solutions chan<- store.Solution,
	) (bool, error) {
		solver, err := solver(input, option)
		if err != nil {
			return false, err
		}
		for solution := range solver.All(ctx) {
			solutions <- solution
		}
		return false, nil
	}
	runner := CliRunner(algorithm)
	return runner.Run(context.Background())
}

// IOData describes the data that is used in the IOProducer.
type IOData interface {
	Input() any
	Option() any
	Writer() any
}

// NewIOData creates a new IOData.
func NewIOData(input any, option any, writer any) IOData {
	return ioData{
		input:  input,
		option: option,
		writer: writer,
	}
}

type ioData struct {
	input  any
	option any
	writer any
}

func (d ioData) Input() any {
	return d.input
}

func (d ioData) Option() any {
	return d.option
}

func (d ioData) Writer() any {
	return d.writer
}

// IOProducer is a function that produces the input, option and writer.
type IOProducer func(context.Context, any) IOData

// InputDecoder is a function that decodes a reader into a struct.
type InputDecoder[Input any] func(context.Context, any) (Input, error)

// OptionDecoder is a function that decodes a reader into a struct.
type OptionDecoder[Option any] func(
	context.Context, any, Option,
) (Option, error)

// FlagParser is a function that parses flags.
type FlagParser[Input any] func() (any, Input, error)

// JSONDecoder is a Decoder that decodes a json into a struct.
func JSONDecoder[Input any](
	_ context.Context, reader any) (input Input, err error,
) {
	ioReader, ok := reader.(io.Reader)
	if !ok {
		return input, errors.New(
			"JsonDecoder is not compatible with configured IOProducer",
		)
	}
	decoder := json.NewDecoder(ioReader)
	err = decoder.Decode(&input)
	return input, err
}

// NoopOptionsDecoder is a Decoder that returns the option as is.
func NoopOptionsDecoder[Input any](
	_ context.Context, _ any, input Input,
) (Input, error) {
	return input, nil
}

// DefaultFlagParser parses flags and env vars.
func DefaultFlagParser[Option, RunnerCfg any]() (
	runnerConfig RunnerCfg, option Option, err error,
) {
	// create a FlagSetFiller
	filler := flagsfiller.New(
		flagsfiller.WithEnv(""),
		flagsfiller.WithFieldRenamer(
			func(name string) string {
				repl := strings.ReplaceAll(name, "-", ".")
				return strings.ToLower(repl)
			},
		),
	)
	err = filler.Fill(flag.CommandLine, &option)
	if err != nil {
		return runnerConfig, option, err
	}

	err = filler.Fill(flag.CommandLine, &runnerConfig)
	if err != nil {
		return runnerConfig, option, err
	}

	flag.Parse()

	return runnerConfig, option, nil
}

// Algorithm is a function that runs an algorithm.
type Algorithm[Input, Option, Solution any] func(
	context.Context, Input, Option, chan<- Solution,
) (bool, error)

// DefaultIOProducer is a test IOProducer.
func DefaultIOProducer(_ context.Context, config any) IOData {
	cfg, ok := config.(CliRunnerConfig)
	if !ok {
		log.Fatal("DefaultIOProducer is not compatible with the runner")
	}
	reader := os.Stdin
	if cfg.Runner.Input.Path != "" {
		r, err := os.Open(cfg.Runner.Input.Path)
		if err != nil {
			log.Fatal(err)
		}
		reader = r
	}
	writer := os.Stdout
	if cfg.Runner.Output.Path != "" {
		w, err := os.Create(cfg.Runner.Output.Path)
		if err != nil {
			log.Fatal(err)
		}
		writer = w
	}
	return NewIOData(
		reader,
		nil,
		writer,
	)
}

// Encoder is a function that encodes a struct into a writer.
type Encoder[Solution any] func(
	context.Context, <-chan Solution, any, any,
) error

// JSONEncoder is an Encoder that encodes a struct into a json.
func JSONEncoder[Solution any](
	_ context.Context, solutions <-chan Solution, writer any, runnerCfg any,
) error {
	ioWriter, ok := writer.(io.Writer)
	if !ok {
		return errors.New("JsonEncoder is not compatible with configured IOProducer")
	}
	runnerConfig, ok := runnerCfg.(CliRunnerConfig)
	if !ok {
		return errors.New("JsonEncoder is not compatible with configured IOProducer")
	}
	encoder := json.NewEncoder(ioWriter)
	switch runnerConfig.Runner.Output.Solutions {
	case All:
		_, err := ioWriter.Write([]byte("["))
		if err != nil {
			return err
		}
		for solution := range solutions {
			err = encoder.Encode(solution)
			if err != nil {
				return err
			}
		}
		_, err = ioWriter.Write([]byte("]"))
		if err != nil {
			return err
		}
	case Last:
		var last Solution
		for solution := range solutions {
			last = solution
		}
		err := encoder.Encode(last)
		if err != nil {
			return err
		}
	}

	return nil
}
