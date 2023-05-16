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
	Model factory.Options                `json:"model,omitempty"`
	Solve nextroute.ParallelSolveOptions `json:"solve,omitempty"`
}

func solver(
	ctx context.Context,
	input schema.Input,
	options options,
) (runSchema.Output, error) {
	model, err := factory.NewModel(input, options.Model)
	if err != nil {
		return runSchema.Output{}, err
	}

	solver, err := nextroute.NewParallelSolver(model)
	if err != nil {
		return runSchema.Output{}, err
	}

	solutions, err := solver.Solve(ctx, options.Solve)
	if err != nil {
		return runSchema.Output{}, err
	}

	return nextroute.Format(ctx, options, solver, solutions.Last()), nil
}
