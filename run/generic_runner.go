package run

import (
	"context"
	"log"
	"os"
	"reflect"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
)

type start string
type data string

// Start is the key for the start time of the run.
const Start start = "start"

// Data is the key for additional data of the run.
const Data data = "data"

// GenericRunner creates a new runner from the given components.
func GenericRunner[RunnerConfig, Input, Option, Solution any](
	ioHandler IOProducer[RunnerConfig],
	inputDecoder Decoder[Input],
	inputValidator Validator[Input],
	optionDecoder Decoder[Option],
	handler Algorithm[Input, Option, Solution],
	encoder Encoder[Solution, Option],
) Runner[RunnerConfig, Input, Option, Solution] {
	runnerConfig, option, err := FlagParser[
		Option, RunnerConfig,
	]()
	if err != nil {
		log.Fatal(err)
	}
	return &genericRunner[RunnerConfig, Input, Option, Solution]{
		IOProducer:       ioHandler,
		InputDecoder:     inputDecoder,
		InputValidator:   inputValidator,
		OptionDecoder:    optionDecoder,
		Algorithm:        handler,
		Encoder:          encoder,
		runnerConfig:     runnerConfig,
		flagParsedOption: option,
	}
}

type genericRunner[RunnerConfig, Input, Option, Solution any] struct {
	IOProducer       IOProducer[RunnerConfig]
	InputDecoder     Decoder[Input]
	InputValidator   Validator[Input]
	OptionDecoder    Decoder[Option]
	Algorithm        Algorithm[Input, Option, Solution]
	Encoder          Encoder[Solution, Option]
	runnerConfig     RunnerConfig
	flagParsedOption Option
}

func (r *genericRunner[RunnerConfig, Input, Option, Solution]) handleCPUProfile(
	runnerConfig any,
) (deferFunc func() error, err error) {
	deferFunc = func() error {
		return nil
	}
	if cpuProfiler, ok := runnerConfig.(CPUProfiler); ok &&
		cpuProfiler.CPUProfilePath() != "" {
		// CPU profiler.
		f, err := os.Create(cpuProfiler.CPUProfilePath())
		if err != nil {
			return deferFunc, err
		}
		deferFunc = func() error {
			return f.Close()
		}

		if err := pprof.StartCPUProfile(f); err != nil {
			return deferFunc, err
		}
		deferFunc = func() error {
			pprof.StopCPUProfile()
			return f.Close()
		}
	}
	return deferFunc, nil
}

func (r *genericRunner[RunnerConfig, Input, Option, Solution],
) handleMemoryProfile(runnerConfig any,
) (deferFunc func() error, err error) {
	deferFunc = func() error {
		return nil
	}
	// Memory profile.
	if memoryProfiler, ok := runnerConfig.(MemoryProfiler); ok &&
		memoryProfiler.MemoryProfilePath() != "" {
		f, err := os.Create(memoryProfiler.MemoryProfilePath())
		if err != nil {
			return deferFunc, err
		}
		deferFunc = func() error {
			return f.Close()
		}

		// Clean up unused objects from the heap before profiling. But do not
		// garbage collect the runner, so we can see in-use memory.
		runtime.GC()
		runtime.KeepAlive(r)

		if err := pprof.WriteHeapProfile(f); err != nil {
			return deferFunc, err
		}
	}
	return deferFunc, nil
}

func (r *genericRunner[RunnerConfig, Input, Option, Solution]) Run(
	ctx context.Context,
) (retErr error) {
	start := time.Now()
	ctx = context.WithValue(ctx, Start, start)
	ctx = context.WithValue(ctx, Data, &sync.Map{})
	// handle CPU profile
	deferFuncCPU, retErr := r.handleCPUProfile(r.runnerConfig)
	if retErr != nil {
		return retErr
	}
	defer func() {
		err := deferFuncCPU()
		// the first error is more important
		if retErr == nil {
			retErr = err
		}
	}()
	// get IO
	ioData, retErr := r.IOProducer(ctx, r.runnerConfig)
	if retErr != nil {
		return retErr
	}

	if r.InputValidator != nil {
		retErr = r.InputValidator(ctx, ioData.Input())
		if retErr != nil {
			return retErr
		}
	}

	// decode input
	decodedInput, retErr := r.InputDecoder(ctx, ioData.Input())
	if retErr != nil {
		return retErr
	}

	// use options configured in runner via flags and environment variables
	decodedOption := r.flagParsedOption
	// decode option if provided
	tempOption, err := r.OptionDecoder(ctx, ioData.Option())
	if err != nil {
		return err
	}
	var defaultOption Option
	// if option is not default, use it
	if !reflect.DeepEqual(tempOption, defaultOption) {
		decodedOption = tempOption
	}

	// run algorithm
	solutions := make(chan Solution)
	errs := make(chan error, 1)
	go func() {
		defer close(solutions)
		defer close(errs)
		retErr = r.Algorithm(ctx, decodedInput, decodedOption, solutions)
		if retErr != nil {
			errs <- retErr
			return
		}
	}()

	// encode solutions
	retErr = r.Encoder.Encode(
		ctx, solutions, ioData.Writer(), r.runnerConfig, decodedOption,
	)
	if retErr != nil {
		return retErr
	}

	// handle memory profile
	deferFuncMemory, retErr := r.handleMemoryProfile(r.runnerConfig)
	if retErr != nil {
		return retErr
	}

	defer func() {
		err := deferFuncMemory()
		// the first error is more important
		if retErr == nil {
			retErr = err
		}
	}()

	// return potential errors
	return <-errs
}

func (r *genericRunner[RunnerConfig, Input, Option, Solution]) SetIOProducer(
	ioProducer IOProducer[RunnerConfig],
) {
	r.IOProducer = ioProducer
}

func (r *genericRunner[RunnerConfig, Input, Option, Solution]) SetInputDecoder(
	decoder Decoder[Input],
) {
	r.InputDecoder = decoder
}

func (r *genericRunner[RunnerConfig, Input, Option, Solution]) SetInputValidator(
	validator Validator[Input],
) {
	r.InputValidator = validator
}

func (r *genericRunner[RunnerConfig, Input, Option, Solution]) SetOptionDecoder(
	decoder Decoder[Option],
) {
	r.OptionDecoder = decoder
}

func (r *genericRunner[RunnerConfig, Input, Option, Solution]) SetAlgorithm(
	algorithm Algorithm[Input, Option, Solution],
) {
	r.Algorithm = algorithm
}

func (r *genericRunner[RunnerConfig, Input, Option, Solution]) SetEncoder(
	encoder Encoder[Solution, Option],
) {
	r.Encoder = encoder
}

func (r *genericRunner[
	RunnerConfig, Input, Option, Solution,
]) GetEncoder() Encoder[Solution, Option] {
	return r.Encoder
}

func (r *genericRunner[
	RunnerConfig, Input, Option, Solution],
) RunnerConfig() RunnerConfig {
	return r.runnerConfig
}
