package check

import (
	"context"

	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute"
	runSchema "github.com/nextmv-io/sdk/run/schema"
)

var (
	con        = connect.NewConnector("sdk", "NextRouteCheck")
	modelCheck func(
		context.Context,
		nextroute.Model,
		Options,
	) (Output, error)
	solutionCheck func(
		context.Context,
		nextroute.Solution,
		Options,
	) (Output, error)

	format func(
		context.Context,
		any,
		Options,
		nextroute.Progressioner,
		...nextroute.Solution,
	) (runSchema.Output, error)
)
