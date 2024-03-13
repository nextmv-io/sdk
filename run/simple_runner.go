package run

import (
	"context"
)

// CLI instantiates a CLIRunner and runs it. This is a wrapper function that
// allows for simple usage of the CLIRunner.
func CLI[Input, Option, Output any](solver func(
	ctx context.Context, input Input, option Option) (solutions Output, err error),
	options ...RunnerOption[CLIRunnerConfig, Input, Option, Output],
) Runner[CLIRunnerConfig, Input, Option, Output] {
	algorithm := func(
		ctx context.Context,
		input Input, option Option, out chan<- Output,
	) error {
		output, err := solver(ctx, input, option)
		if err != nil {
			return err
		}
		out <- output
		return nil
	}
	return NewCLIRunner(algorithm, options...)
}

// HTTP instantiates an HTTPRunner and runs it. The default port is 9000 and
// protocol is HTTP. Pass HTTPRunnerOptions to change these settings.
func HTTP[Input, Option, Output any](solver func(
	ctx context.Context, input Input, option Option) (solutions Output, err error),
	options ...HTTPRunnerOption[Input, Option, Output],
) HTTPRunner[HTTPRunnerConfig, Input, Option, Output] {
	algorithm := func(
		ctx context.Context,
		input Input, option Option, out chan<- Output,
	) error {
		output, err := solver(ctx, input, option)
		if err != nil {
			return err
		}
		out <- output

		return nil
	}
	return NewHTTPRunner(algorithm, options...)
}
