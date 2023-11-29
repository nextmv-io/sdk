package nextroute

import (
	"github.com/nextmv-io/sdk/alns"
	"github.com/nextmv-io/sdk/connect"
)

// NewSolver creates a new solver. The solver can be used to solve a solution.
// The solution passed to the solver is the starting point of the solver. The
// solver will try to improve the solution.
func NewSolveOperatorUnplan(
	numberOfUnits alns.SolveParameter[Solution],
) (alns.SolveOperator[Solution], error) {
	connect.Connect(con, &newSolveOperatorUnplan)
	return newSolveOperatorUnplan(numberOfUnits)
}
