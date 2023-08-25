// package main holds the implementation of the nextroute template.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"github.com/nextmv-io/sdk"
	"github.com/nextmv-io/sdk/alns"
	"github.com/nextmv-io/sdk/nextroute"
	"github.com/nextmv-io/sdk/nextroute/common"
	"github.com/nextmv-io/sdk/nextroute/factory"
	"github.com/nextmv-io/sdk/nextroute/schema"
	"github.com/nextmv-io/sdk/run"
	"github.com/nextmv-io/sdk/run/statistics"
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

	b, err := json.Marshal(nextrouteInput)
	os.WriteFile("test.json", b, 0644)

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

	distanceExpression := makeDistanceExpression(input)

	model, err := factory.NewModel(nextrouteInput, options.Model)
	if err != nil {
		return FleetOutput{}, err
	}

	for _, v := range model.VehicleTypes() {
		v.SetData(distanceData{
			distance: distanceExpression,
		})
	}

	solver, err := nextroute.NewParallelSolver(model)
	if err != nil {
		return FleetOutput{}, err
	}

	solutionsCount := "last"
	runSolutions, err := run.ParseSolutions(solutionsCount)
	if err != nil {
		return FleetOutput{}, err
	}

	solutions, err := solver.Solve(ctx, options.Solve)
	if err != nil {
		return FleetOutput{}, err
	}

	// Last solution == 1, all solutions == 0
	if runSolutions == 1 {
		last := solutions.Last()
		output, err := format(ctx, options.Solve.Duration, solver, ToFleetSolutionOutput, last)
		if err != nil {
			return output, err
		}
		return output, nil
	}

	solutionArray := make([]nextroute.Solution, 0)
	for s := range solutions {
		solutionArray = append(solutionArray, s)
	}
	output, err := format(ctx, options.Solve.Duration, solver, ToFleetSolutionOutput, solutionArray...)
	if err != nil {
		return output, err
	}
	return output, nil
}

func format(
	ctx context.Context,
	duration time.Duration,
	progressioner alns.Progressioner,
	toSolutionOutputFn func(nextroute.Solution) (any, error),
	solutions ...nextroute.Solution,
) (FleetOutput, error) {
	mappedSolutions, err := MapWithError(solutions, toSolutionOutputFn)
	if err != nil {
		return FleetOutput{}, err
	}

	output := FleetOutput{
		Solutions: mappedSolutions,
		Options: schema.Options{
			Solver: &schema.SolverOptions{
				Limits: &schema.Limits{
					Duration: duration.String(),
				},
			},
		},
		Hop: struct {
			Version string "json:\"version\""
		}{
			Version: sdk.VERSION,
		},
	}

	startTime := time.Time{}
	if start, ok := ctx.Value(run.Start).(time.Time); ok {
		startTime = start
	}

	progressionValues := progressioner.Progression()

	if len(progressionValues) == 0 {
		return output, errors.New("no solution values or elapsed time values found")
	}

	seriesData := common.Map(
		progressionValues,
		func(progressionEntry alns.ProgressionEntry) statistics.DataPoint {
			return statistics.DataPoint{
				X: statistics.Float64(progressionEntry.ElapsedSeconds),
				Y: statistics.Float64(progressionEntry.Value),
			}
		},
	)

	if len(output.Solutions) == 1 {
		seriesData = seriesData[len(seriesData)-1:]
	}

	if len(output.Solutions) != len(seriesData) {
		return output, errors.New("more or less solution values than solutions found")
	}
	for idx, data := range seriesData {
		if _, ok := output.Solutions[idx].(FleetState); ok {
			output.Statistics.Time.Start = startTime
			output.Statistics.Value = int(data.Y)
			output.Statistics.Time.ElapsedSeconds = float64(data.X)
			output.Statistics.Time.Elapsed = time.Duration(data.X * statistics.Float64(time.Second)).String()
		}
	}

	return output, nil
}

// MapWithError maps a slice of type T to a slice of type R using the function f.
func MapWithError[T any, R any](v []T, f func(T) (R, error)) ([]R, error) {
	r := make([]R, len(v))
	for idx, x := range v {
		o, err := f(x)
		if err != nil {
			return r, err
		}
		r[idx] = o
	}
	return r, nil
}
