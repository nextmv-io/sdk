package run

import "context"

// CLI runs the runner in a simple way returning all solutions.
func CLI[Input, Option, Solution any](solver func(
	input Input, option Option) (solutions []Solution, err error),
	options ...RunnerOption[Input, Option, Solution],
) error {
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
	runner := CLIRunner(algorithm, options...)
	return runner.Run(context.Background())
}

// HTTP runs the HTTPRunner in a simple way returning all solutions.
func HTTP[Input, Option, Solution any](solver func(
	input Input, option Option) (solutions []Solution, err error),
	options ...HTTPRunnerOption[Input, Option, Solution],
) error {
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
	runner := HTTPRunner(algorithm, options...)
	return runner.Run(context.Background())
}
