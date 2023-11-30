package nextroute

import (
	"github.com/nextmv-io/sdk/alns"
	"github.com/nextmv-io/sdk/connect"
)

// NewSolveOperatorRestart creates a new solve operator for nextroute that
// restarts the solver.
func NewSolveOperatorRestart(
	maximumIterations alns.SolveParameter[Solution],
) (alns.SolveOperator[Solution], error) {
	connect.Connect(con, &newSolveOperatorRestart)
	return newSolveOperatorRestart(maximumIterations)
}
