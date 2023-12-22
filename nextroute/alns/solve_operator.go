package alns

import (
	sdkAlns "github.com/nextmv-io/sdk/alns"
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute"
)

// NewSolveOperator returns a new solve operator.
func NewSolveOperator(
	probability float64,
	canResultInImprovement bool,
	parameters sdkAlns.SolveParameters[nextroute.Solution],
) sdkAlns.SolveOperator[nextroute.Solution] {
	connect.Connect(con, &newSolveOperator)
	return newSolveOperator(
		probability,
		canResultInImprovement,
		parameters,
	)
}
