package factory

import (
	"context"
	"github.com/nextmv-io/sdk/alns"
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute"
	"github.com/nextmv-io/sdk/nextroute/schema"
	runSchema "github.com/nextmv-io/sdk/run/schema"
)

var (
	con = connect.NewConnector("sdk", "NextRouteFactory")

	newModel func(
		schema.Input,
		Options,
		...nextroute.Option,
	) (nextroute.Model, error)

	format func(
		context.Context,
		any,
		alns.Progressioner,
		...nextroute.Solution,
	) runSchema.Output
)
