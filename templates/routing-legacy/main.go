// package main holds the implementation of the nextroute template.
package main

import (
	"context"
	"log"
	"time"

	"github.com/nextmv-io/sdk/nextroute"
	"github.com/nextmv-io/sdk/nextroute/factory"
	"github.com/nextmv-io/sdk/nextroute/schema"
	"github.com/nextmv-io/sdk/run"
	runSchema "github.com/nextmv-io/sdk/run/schema"
)

func main() {
	runner := run.CLI(solver,
		run.InputValidate[run.CLIRunnerConfig, schema.FleetInput, options, runSchema.Output](
			nil,
		),
	)
	err := runner.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

type options struct {
	Model  factory.Options                `json:"model,omitempty"`
	Solve  nextroute.ParallelSolveOptions `json:"solve,omitempty"`
	Format nextroute.FormatOptions        `json:"format,omitempty"`
}

func solver(
	ctx context.Context,
	input schema.FleetInput,
	options options,
) (runSchema.Output, error) {
	nextrouteInput, err := input.ToNextRoute()
	if err != nil {
		panic(err)
	}

	solveOptions := options.Solve
	if input.Options != nil && input.Options.Solver != nil &&
		input.Options.Solver.Limits != nil {
		duration, err := time.ParseDuration(input.Options.Solver.Limits.Duration)
		if err != nil {
			return runSchema.Output{}, err
		}

		solveOptions.Duration = duration
	}

	model, err := factory.NewModel(nextrouteInput, options.Model)
	if err != nil {
		return runSchema.Output{}, err
	}

	solver, err := nextroute.NewParallelSolver(model)
	if err != nil {
		return runSchema.Output{}, err
	}

	solutions, err := solver.Solve(ctx, solveOptions)
	if err != nil {
		return runSchema.Output{}, err
	}
	last := solutions.Last()

	output := factory.Format(ctx, options, solver, last)
	output.Statistics.Result.Custom = factory.DefaultCustomResultStatistics(last)

	return output, nil
}
