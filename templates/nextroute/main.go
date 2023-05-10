// package main holds the implementation of the nextroute template.
package main

import (
	"context"
	"log"

	"github.com/nextmv-io/sdk/nextroute"
	"github.com/nextmv-io/sdk/nextroute/factory"
	"github.com/nextmv-io/sdk/nextroute/schema"
	"github.com/nextmv-io/sdk/run"
)

func main() {
	runner := run.NewCLIRunner(solver)
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
	solutions chan<- any,
) error {
	model, err := factory.NewModel(input, options.Model)
	if err != nil {
		return err
	}

	solver, err := nextroute.NewParallelSolver(model)
	if err != nil {
		return err
	}

	solverSolutions, err := solver.Solve(ctx, options.Solve)
	if err != nil {
		return err
	}

	last := solverSolutions.Last()
	if last != nil {
		solutions <- nextroute.NewBasicFormatter().ToOutput(last)
	}

	return nil
}
