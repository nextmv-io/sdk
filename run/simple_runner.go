package run

import (
	"context"

	"github.com/nextmv-io/sdk/store"
)

// CLI instantiates a CLIRunner and runs it. This is a wrapper function that
// allows for simple usage of the CLIRunner.
func CLI[Input, Option, Solution any](solver func(
	input Input, option Option) (solutions []Solution, err error),
	options ...RunnerOption[CLIRunnerConfig, Input, Option, Solution],
) Runner[CLIRunnerConfig, Input, Option, Solution] {
	algorithm := func(
		_ context.Context,
		input Input, option Option, sols chan<- Solution,
	) error {
		solutions, err := solver(input, option)
		if err != nil {
			return err
		}
		for _, sol := range solutions {
			sols <- sol
		}
		return nil
	}
	return NewCLIRunner(algorithm, options...)
}

// HTTP instantiates an HTTPRunner and runs it. The default port is 9000 and
// protocol is HTTP. Pass HTTPRunnerOptions to change these settings.
func HTTP[Input, Option, Solution any](solver func(
	input Input, option Option) (solutions []Solution, err error),
	options ...HTTPRunnerOption[Input, Option, Solution],
) HTTPRunner[HTTPRunnerConfig, Input, Option, Solution] {
	algorithm := func(
		_ context.Context,
		input Input, option Option, sols chan<- Solution,
	) error {
		solutions, err := solver(input, option)
		if err != nil {
			return err
		}
		for _, sol := range solutions {
			sols <- sol
		}
		return nil
	}
	return NewHTTPRunner(algorithm, options...)
}

// AllSolutions is a helper function that unwraps a (store.Solver,error) into
// ([]store.Solution, error).
func AllSolutions(
	solver store.Solver, err error,
) (solutions []store.Solution, retErr error) {
	if err != nil {
		return nil, err
	}
	for solution := range solver.All(context.Background()) {
		solutions = append(solutions, solution)
	}
	return solutions, nil
}
