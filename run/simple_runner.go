package run

import "context"

// CLI runs the runner in a simple way returning all solutions.
func CLI[Input, Option, Solution any](solver func(
	input Input, option Option) ([]Solution, error),
	options ...RunnerOption[Input, Option, Solution],
) error {
	algorithm := func(
		_ context.Context,
		input Input, option Option, solutions chan<- Solution,
	) error {
		sols, err := solver(input, option)
		if err != nil {
			return err
		}
		for _, sol := range sols {
			solutions <- sol
		}
		return nil
	}
	runner := NewCLIRunner(algorithm, options...)
	return runner.Run(context.Background())
}

// HTTP runs the HTTPRunner in a simple way returning all solutions.
func HTTP[Input, Option, Solution any](solver func(
	input Input, option Option) ([]Solution, error),
	options ...HTTPRunnerOption[Input, Option, Solution],
) error {
	algorithm := func(
		_ context.Context,
		input Input, option Option, solutions chan<- Solution,
	) error {
		sols, err := solver(input, option)
		if err != nil {
			return err
		}
		for _, sol := range sols {
			solutions <- sol
		}
		return nil
	}
	runner := NewHTTPRunner(algorithm, options...)
	return runner.Run(context.Background())
}
