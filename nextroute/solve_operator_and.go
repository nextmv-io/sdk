package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
)

// SolveOperatorAnd is a solve-operator which executes a set of solve-operators
// in each iteration.
type SolveOperatorAnd interface {
	SolveOperator

	// Operators returns the solve-operators that will be executed in each
	// iteration.
	Operators() SolveOperators
}

// NewSolverOperatorAnd creates a new solve-and-operator for nextroute.
func NewSolverOperatorAnd(
	probability float64,
	operators SolveOperators,
) (SolveOperatorAnd, error) {
	connect.Connect(con, &newSolverOperatorAnd)
	return newSolverOperatorAnd(probability, operators)
}
