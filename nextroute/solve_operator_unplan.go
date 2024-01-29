package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
)

// SolveOperatorUnPlan is a solve operator that un-plans units.
type SolveOperatorUnPlan interface {
	SolveOperator

	// NumberOfUnits returns the number of units of the solve operator.
	NumberOfUnits() SolveParameter
}

// NewSolveOperatorUnPlan creates a new solve operator for nextroute that
// un-plans units.
func NewSolveOperatorUnPlan(
	numberOfUnits SolveParameter,
) (SolveOperatorUnPlan, error) {
	connect.Connect(con, &newSolveOperatorUnPlan)
	return newSolveOperatorUnPlan(numberOfUnits)
}
