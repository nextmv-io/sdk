package check

import (
	"context"

	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute"
)

// ModelCheck is the check of a model returning a [Output].
func ModelCheck(
	ctx context.Context,
	model nextroute.Model,
	options Options,
) (Output, error) {
	connect.Connect(con, &modelCheck)
	return modelCheck(ctx, model, options)
}

// SolutionCheck is the check of a solution returning a [Output].
func SolutionCheck(
	ctx context.Context,
	solution nextroute.Solution,
	options Options,
) (Output, error) {
	connect.Connect(con, &solutionCheck)
	return solutionCheck(ctx, solution, options)
}
