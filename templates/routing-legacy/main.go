// package main holds the implementation of the nextroute template.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nextmv-io/sdk/nextroute"
	"github.com/nextmv-io/sdk/nextroute/factory"
	"github.com/nextmv-io/sdk/nextroute/schema"
	"github.com/nextmv-io/sdk/run"
)

func main() {
	runner := run.CLI(solver,
		run.InputValidate[run.CLIRunnerConfig, schema.FleetInput, options, FleetOutput](
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
) (FleetOutput, error) {
	nextrouteInput, err := input.ToNextRoute()
	if err != nil {
		panic(err)
	}

	if input.Options != nil && input.Options.Solver != nil &&
		input.Options.Solver.Limits != nil {
		if input.Options.Solver.Limits.Duration != "" {
			duration, err := time.ParseDuration(input.Options.Solver.Limits.Duration)
			if err != nil {
				return FleetOutput{}, err
			}
			options.Solve.Duration = duration
		}
	}

	model, err := factory.NewModel(nextrouteInput, options.Model)
	if err != nil {
		return FleetOutput{}, err
	}

	distanceExpression := makeDistanceExpression(input)
	for _, v := range model.VehicleTypes() {
		v.SetData(distanceData{
			distance: distanceExpression,
		})
	}

	solver, err := nextroute.NewParallelSolver(model)
	if err != nil {
		return FleetOutput{}, err
	}

	solutions, err := solver.Solve(ctx, options.Solve)
	if err != nil {
		return FleetOutput{}, err
	}

	runSolutions := run.Last
	solutionArray := make([]nextroute.Solution, 0)
	switch runSolutions {
	case run.Last:
		solutionArray = append(solutionArray, solutions.Last())
	case run.All:
		for s := range solutions {
			solutionArray = append(solutionArray, s)
		}
	default:
		return FleetOutput{},
			fmt.Errorf("%s is an invalid value for parameter runner.output.solutions. it must be 'all' or 'last'", runSolutions)
	}
	output, err := format(ctx, options.Solve.Duration, solver, ToFleetSolutionOutput, solutionArray...)
	if err != nil {
		return output, err
	}
	return output, nil
}
