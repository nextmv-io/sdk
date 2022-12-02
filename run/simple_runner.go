package run

import "context"

// Simple runs the runner in a simple way.
func Simple[Input, Option, Solution any](solver func(
	input Input, option Option) (Solution, error),
) error {
	algorithm := func(
		_ context.Context,
		input Input, option Option, solutions chan<- Solution,
	) error {
		solution, err := solver(input, option)
		if err != nil {
			return err
		}
		solutions <- solution
		return nil
	}
	runner := CliRunner(algorithm)
	return runner.Run(context.Background())
}
