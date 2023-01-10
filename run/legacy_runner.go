package run

import (
	"context"

	"github.com/nextmv-io/sdk/store"
)

// Run runs the runner. This is a legacy function that is kept for backwards
// compatibility.
func Run[Input any](solver func(
	input Input, option store.Options,
) (store.Solver, error),
	options ...RunnerOption[CLIRunnerConfig, Input, store.Options, store.Solution],
) (Runner[CLIRunnerConfig, Input, store.Options, store.Solution], error) {
	algorithm := func(
		ctx context.Context,
		input Input, option store.Options, solutions chan<- store.Solution,
	) error {
		solver, err := solver(input, option)
		if err != nil {
			return err
		}
		for solution := range solver.All(ctx) {
			solutions <- solution
		}
		return nil
	}
	runner := NewCLIRunner(algorithm, options...)
	return runner, runner.Run(context.Background())
}
