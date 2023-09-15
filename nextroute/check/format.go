package check

import (
	"context"

	"github.com/nextmv-io/sdk/alns"
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute"
	runSchema "github.com/nextmv-io/sdk/run/schema"
)

// Format formats a solution in a basic format using factory.ToSolutionOutput
// to format each solution and also allows to check the solutions and add the
// check to the output of each solution.
func Format(
	ctx context.Context,
	options any,
	checkOptions Options,
	progressioner alns.Progressioner,
	solutions ...nextroute.Solution,
) (runSchema.Output, error) {
	connect.Connect(con, &format)
	return format(ctx, options, checkOptions, progressioner, solutions...)
}
