package alns

import (
	sdkAlns "github.com/nextmv-io/sdk/alns"
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute"
)

// NewSolverOperatorAnd creates a new solve-and-operator for nextroute.
func NewSolverOperatorAnd(
	probability float64,
	operators sdkAlns.SolveOperators[nextroute.Solution],
) (sdkAlns.SolveOperatorAnd[nextroute.Solution], error) {
	connect.Connect(con, &newSolverOperatorAnd)
	return newSolverOperatorAnd(probability, operators)
}
