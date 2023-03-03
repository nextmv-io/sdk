package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
)

// InwardnessConstraint is a constraint that limits the vehicles a plan cluster
// can be added to. A plan cluster can only be added to a vehicles whose
// centroid is closer to the plan cluster than the centroid of any other
// vehicle.
type InwardnessConstraint interface {
	ConstraintDataUpdater
	ModelConstraint
}

// NewInwardnessConstraint creates a new inwardness constraint. The constraint
// needs to be added to the model to be taken into account.
func NewInwardnessConstraint() (InwardnessConstraint, error) {
	connect.Connect(con, &newInwardnessConstraint)
	return newInwardnessConstraint()
}
