// Package main allows you to run a nextroute solver from the command line
// without the need of compiling plugins.
package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/nextmv-io/sdk/nextroute"
	"github.com/nextmv-io/sdk/nextroute/factory"
	schema "github.com/nextmv-io/sdk/nextroute/schema"
	"github.com/nextmv-io/sdk/run"
	"github.com/nextmv-io/sdk/run/encode"
)

// Options holds the options for the solver.
type Options struct {
	factory.ModelOptions
	SolveOptions nextroute.ParallelSolveOptions
}

func main() {
	runner := run.NewCLIRunner(
		solver,
		run.Encode[run.CLIRunnerConfig, schema.Input](
			GenericEncoder[schema.JSONBasicSolution, Options](encode.JSON()),
		),
	)
	err := runner.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func solver(ctx context.Context,
	input schema.Input,
	opts Options,
	returnSolutions chan<- schema.JSONBasicSolution,
) error {
	// Use options from input if provided.
	if input.Options != nil {
		optionMap := input.Options.(map[string]interface{})
		if duration, ok := optionMap["max_duration"]; ok {
			opts.SolveOptions.MaximumDuration = time.Duration(
				duration.(float64)) * time.Second
		}
	}

	// This lines toggles the usage of service durations if present in the input.
	opts.ModelOptions.IgnoreServiceDurations = false

	model, err := factory.NewModel(input, opts.ModelOptions)
	if err != nil {
		return err
	}

	parallelSolver, err := nextroute.NewParallelSolver(model)
	if err != nil {
		return err
	}

	solutions, err := parallelSolver.Solve(ctx, opts.SolveOptions)
	if err != nil {
		return err
	}

	solution := solutions.Last()

	if solution != nil {
		output := nextroute.NewBasicFormatter().ToOutput(solution)
		var jsonOutput schema.JSONBasicSolution
		b, err := json.Marshal(output)
		if err != nil {
			return err
		}
		err = json.Unmarshal(b, &jsonOutput)
		if err != nil {
			return err
		}
		returnSolutions <- jsonOutput
	}

	return nil
}
