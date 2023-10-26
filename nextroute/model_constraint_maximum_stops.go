package nextroute

import "github.com/nextmv-io/sdk/connect"

// MaximumStopsConstraint is a constraint that limits the maximum number of
// stops a vehicle type can have. The maximum number of stops is defined by
// the maximum stops expression. The first stop of a vehicle is not counted
// as a stop and the last stop of a vehicle is not counted as a stop.
type MaximumStopsConstraint interface {
	ModelConstraint

	// MaximumStops returns the maximum stops expression which defines the
	// maximum number of stops a vehicle type can have.
	MaximumStops() VehicleTypeExpression
}

// NewMaximumStopsConstraint creates a new maximum stops constraint. The
// constraint needs to be added to the model to be taken into account.
func NewMaximumStopsConstraint(
	maximumStops VehicleTypeExpression,
) (MaximumStopsConstraint, error) {
	connect.Connect(con, &newMaximumStopsConstraint)
	return newMaximumStopsConstraint(maximumStops)
}
