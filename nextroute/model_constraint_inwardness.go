package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
)

// Inwardness is a constraint that limits the vehicles a plan cluster
// can be added to. A plan cluster can only be added to a vehicles whose
// centroid is closer to the plan cluster than the centroid of any other
// vehicle.
type Inwardness interface {
	ConstraintDataUpdater
	ModelConstraint
	ModelObjective
}

// NewInwardness creates a new inwardness constraint. The constraint
// needs to be added to the model to be taken into account.
func NewInwardness() (Inwardness, error) {
	connect.Connect(con, &newInwardness)
	return newInwardness()
}
