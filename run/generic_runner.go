package run

import (
	"context"
	"os"
	"runtime"
	"runtime/pprof"
)

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

func (r *genericRunner[Input, Option, Solution]) handleCPUProfile(
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
		defer pprof.StopCPUProfile()
	}
	return deferFunc, nil
}

func (r *genericRunner[Input, Option, Solution]) handleMemoryProfile(
	runnerConfig any,
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

func (r *genericRunner[Input, Option, Solution]) Run(
	context context.Context,
) (retErr error) {
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
	ioData := r.IOProducer(context, r.runnerConfig)

	// decode input
	decodedInput, retErr := r.InputDecoder(context, ioData.Input())
	if retErr != nil {
		return retErr
	}

	// decode option
	r.decodedOption, retErr = r.OptionDecoder(
		context, ioData.Option(), r.decodedOption,
	)
	if retErr != nil {
		return retErr
	}

	// run algorithm
	solutions := make(chan Solution)
	errs := make(chan error, 1)
	go func() {
		defer close(solutions)
		defer close(errs)
		retErr = r.Algorithm(context, decodedInput, r.decodedOption, solutions)
		if retErr != nil {
			errs <- retErr
			return
		}
	}()

	// encode solutions
	retErr = r.Encoder(
		context, solutions, ioData.Writer(), r.runnerConfig, r.decodedOption,
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
