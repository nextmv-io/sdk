package nextroute

import (
	"github.com/nextmv-io/sdk/alns"
	"github.com/nextmv-io/sdk/connect"
)

// NewSolveOperatorUnPlan creates a new solve operator for nextroute that
// unplans units.
func NewSolveOperatorUnPlan(
	numberOfUnits alns.SolveParameter[Solution],
) (alns.SolveOperator[Solution], error) {
	connect.Connect(con, &newSolveOperatorUnPlan)
	return newSolveOperatorUnPlan(numberOfUnits)
}
