package alns

import (
	sdkAlns "github.com/nextmv-io/sdk/alns"
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute"
)

// NewSolverOperatorOr creates a new solve-or-operator for nextroute.
func NewSolverOperatorOr(
	loops int,
	probability float64,
	operators sdkAlns.SolveOperators[nextroute.Solution],
) (sdkAlns.SolveOperatorOr[nextroute.Solution], error) {
	connect.Connect(con, &newSolverOperatorAnd)
	return newSolverOperatorOr(loops, probability, operators)
}
