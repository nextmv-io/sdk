package nextroute

import (
	"context"

	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/run/schema"
)

// FormatOptions are the options that influence the format of the output.
type FormatOptions struct {
	Disable struct {
		Progression bool `json:"progression" usage:"disable the progression series"`
	} `json:"disable"`
}

// Format formats a solution in basic format using the map function
// toSolutionOutputFn to map a solution to a user specific format.
func Format(
	ctx context.Context,
	options any,
	progressioner Progressioner,
	toSolutionOutputFn func(Solution) any,
	solutions ...Solution,
) schema.Output {
	connect.Connect(con, &format)
	return format(
		ctx,
		options,
		progressioner,
		toSolutionOutputFn,
		solutions...,
	)
}
