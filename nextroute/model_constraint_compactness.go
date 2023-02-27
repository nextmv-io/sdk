package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
)

// CompactnessConstraint is a constraint that limits the distance a new
// plan cluster can be from the centroid of the vehicle it will be added to.
// The maximim distance is defined by the limit expression.
type CompactnessConstraint interface {
	ModelConstraint

	// Limit returns the limit expression which defines the maximum distance
	// a new plan cluster can be from the centroid of the vehicle it will be
	// added to.
	Limit() StopExpression
}

// NewCompactnessConstraint creates a new compactness constraint. The
// constraint needs to be added to the model to be taken into account.
func NewCompactnessConstraint(
	limit StopExpression,
) (CompactnessConstraint, error) {
	connect.Connect(con, &newCompactnessConstraint)
	return newCompactnessConstraint(limit)
}
