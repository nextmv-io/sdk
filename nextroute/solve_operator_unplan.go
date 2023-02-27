package nextroute

import "github.com/nextmv-io/sdk/connect"

// NewSolveOperatorUnPlan creates a new solve-operator which un-plans planned
// plan-clusters.
func NewSolveOperatorUnPlan(
	numberOfUnits SolveParameter,
) SolveOperatorUnPlan {
	connect.Connect(solverConnect, &newSolveOperatorUnPlan)
	return newSolveOperatorUnPlan(numberOfUnits)
}

// SolveOperatorUnPlan is a solve-operator which un-plans planned plan-clusters.
// It is used to remove planned plan-clusters from the solution.
// In each iteration of the solve run, the number of plan-clusters to un-plan
// is determined by the number of units. The number of units is a
// solve-parameter which can be configured by the user. In each iteration, the
// number of units is sampled from a uniform distribution. The number of units
// is always an integer between 1 and the number of units.
type SolveOperatorUnPlan interface {
	SolveOperator

	// NumberOfUnits returns the number of units to unplan as a solve-parameter.
	// Solve-parameters can change value during the solve run.
	NumberOfUnits() SolveParameter
}
