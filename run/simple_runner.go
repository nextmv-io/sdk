package run

import "context"

// CliLast runs the runner in a simple way and returns the last solution.
func CliLast[Input, Option, Solution any](solver func(
	input Input, option Option) (Solution, error),
	options ...RunnerOption[Input, Option, Solution],
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
	runner := CliRunner(algorithm, options...)
	return runner.Run(context.Background())
}

// CliAll runs the runner in a simple way returning all solutions.
func CliAll[Input, Option, Solution any](solver func(
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
	runner := CliRunner(algorithm, options...)
	return runner.Run(context.Background())
}

// HTTPLast runs the HTTPRunner in a simple way and returns the last solution.
func HTTPLast[Input, Option, Solution any](solver func(
	input Input, option Option) (Solution, error),
	options ...HTTPRunnerOption[Input, Option, Solution],
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
	runner := NewHTTPRunner(algorithm, options...)
	return runner.Run(context.Background())
}

// HTTPAll runs the HTTPRunner in a simple way returning all solutions.
func HTTPAll[Input, Option, Solution any](solver func(
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
