package run

import (
	"context"

	"github.com/nextmv-io/sdk/store"
)

// Run runs the runner.
func Run[Input, Option any](solver func(
	input Input, option Option,
) (store.Solver, error),
) error {
	algorithm := func(
		ctx context.Context,
		input Input, option Option, solutions chan<- store.Solution,
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
	runner := CliRunner(algorithm)
	return runner.Run(context.Background())
}
