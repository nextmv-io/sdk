package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
)

// SolveOperatorPlan is a solve-operator that tries to plan all unplanned
// plan-units in each iteration. The group-size is a solve-parameter which
// can be configured by the user. Group-size determines how many random
// plan-units are selected for which the best move is determined. The best
// move is then executed. In one iteration of the solve run, the operator will
// continue to select a random group-size number of unplanned plan-units
// and execute the best move until all unplanned plan-units are planned or
// no more moves can be executed. In an unconstrained model all plan-units
// will be planned after one iteration of this operator.
type SolveOperatorPlan interface {
	SolveOperator

	// GroupSize returns the group size of the solve operator.
	GroupSize() SolveParameter
}

// NewSolveOperatorPlan creates a new solve operator for nextroute that
// plans units.
func NewSolveOperatorPlan(
	groupSize SolveParameter,
) (SolveOperatorPlan, error) {
	connect.Connect(con, &newSolveOperatorPlan)
	return newSolveOperatorPlan(groupSize)
}
