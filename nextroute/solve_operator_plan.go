package nextroute

import (
	"github.com/nextmv-io/sdk/alns"
	"github.com/nextmv-io/sdk/connect"
)

// NewSolveOperatorPlan creates a new solve operator for nextroute that
// plans units.
func NewSolveOperatorPlan(
	groupSize alns.SolveParameter[Solution],
) (alns.SolveOperator[Solution], error) {
	connect.Connect(con, &newSolveOperatorPlan)
	return newSolveOperatorPlan(groupSize)
}
