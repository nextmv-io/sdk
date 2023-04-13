package nextroute

import "github.com/nextmv-io/sdk/connect"

// NewUnPlannedObjective returns a new UnPlannedObjective that uses the
// un-planned stops as an objective. Each unplanned stop is scored by the
// given expression.
func NewUnPlannedObjective(expression StopExpression) UnPlannedObjective {
	connect.Connect(con, &newUnPlannedObjective)
	return newUnPlannedObjective(expression)
}

// UnPlannedObjective is an objective that uses the un-planned stops as an
// objective. Each unplanned stop is scored by the given expression.
type UnPlannedObjective interface {
	ModelObjective
}
