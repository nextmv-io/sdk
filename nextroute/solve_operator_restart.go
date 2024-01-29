package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
)

// SolveOperatorRestart is a solve operator that restarts the solver.
// The operator will set the working solution to the best solution found so far
// after MaximumIterations number of iterations without finding a better
// solution.
type SolveOperatorRestart interface {
	SolveOperator

	// MaximumIterations returns the maximum iterations of the solve operator.
	MaximumIterations() SolveParameter
}

// NewSolveOperatorRestart creates a new solve operator for nextroute that
// restarts the solver.
func NewSolveOperatorRestart(
	maximumIterations SolveParameter,
) (SolveOperatorRestart, error) {
	connect.Connect(con, &newSolveOperatorRestart)
	return newSolveOperatorRestart(maximumIterations)
}
