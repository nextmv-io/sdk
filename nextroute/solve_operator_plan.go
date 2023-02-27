package nextroute

import "github.com/nextmv-io/sdk/connect"

// NewSolveOperatorPlan creates a new solve-operator-plan.
func NewSolveOperatorPlan(
	groupSize SolveParameter,
) SolveOperatorPlan {
	connect.Connect(solverConnect, &newSolveOperatorPlan)
	return newSolveOperatorPlan(groupSize)
}

// SolveOperatorPlan is a solve-operator that tries to plan all unplanned
// plan-clusters in each iteration. The group-size is a solve-parameter which
// can be configured by the user. Group-size determines how many random
// plan-clusters are selected for which the best move is determined. The best
// move is then executed. In one iteration of the solve run, the operator will
// continue to select a random group-size number of unplanned plan-clusters
// and execute the best move until all unplanned plan-clusters are planned or
// no more moves can be executed. In an unconstrained model all plan-clusters
// will be planned after one iteration of this operator.
type SolveOperatorPlan interface {
	SolveOperator

	GroupSize() SolveParameter
}
