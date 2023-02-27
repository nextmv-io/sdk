package nextroute

import "github.com/nextmv-io/sdk/connect"

// NewSolveOperatorRestart creates a new solve-operator that restarts the solver
// after a certain number of iterations without improvement.
func NewSolveOperatorRestart(
	maximumIterations SolveParameter,
) SolveOperatorRestart {
	connect.Connect(solverConnect, &newSolveOperatorRestart)
	return newSolveOperatorRestart(maximumIterations)
}

// SolveOperatorRestart is a solve-operator that restarts the solver after a
// certain number of iterations without improvement. The restart is done by
// invoking the Restart method on the solver and replaces the current work
// solution with the best solution found so far.
type SolveOperatorRestart interface {
	SolveOperator

	// MaximumIterations returns the maximum number of iterations without
	// improvement before the solver is restarted.
	MaximumIterations() SolveParameter
}
