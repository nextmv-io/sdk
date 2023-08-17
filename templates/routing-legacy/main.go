// package main holds the implementation of the nextroute template.
package main

import (
	"context"
	"log"

	"github.com/nextmv-io/sdk/nextroute"
	"github.com/nextmv-io/sdk/nextroute/factory"
	"github.com/nextmv-io/sdk/nextroute/schema"
	"github.com/nextmv-io/sdk/run"
	runSchema "github.com/nextmv-io/sdk/run/schema"
)

func main() {
	runner := run.CLI(solver)
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

	opt := options.Solve
	if input.Options != nil || input.Options.Solver != nil &&
		input.Options.Solver.Limits != nil && input.Options.Solver.Limits.Duration != nil {
		opt.Duration = input.Options.Solver.Limits.Duration.Duration
	}

	model, err := factory.NewModel(nextrouteInput, options.Model)
	if err != nil {
		return runSchema.Output{}, err
	}

	solver, err := nextroute.NewParallelSolver(model)
	if err != nil {
		return runSchema.Output{}, err
	}

	solutions, err := solver.Solve(ctx, opt)
	if err != nil {
		return runSchema.Output{}, err
	}
	last := solutions.Last()

	output := factory.Format(ctx, options, solver, last)
	output.Statistics.Result.Custom = factory.DefaultCustomResultStatistics(last)

	return output, nil
}
