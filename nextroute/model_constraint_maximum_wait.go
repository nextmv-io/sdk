package nextroute

import "github.com/nextmv-io/sdk/connect"

// MaximumWaitConstraint is a constraint that limits the wait time of a vehicle.
type MaximumWaitConstraint interface {
	ModelConstraint

	// Maximum returns the maximum expression which defines the maximum
	// wait time of a vehicle type.
	Maximum() VehicleTypeDurationExpression
}

// NewMaximumWaitConstraint returns a new MaximumWaitConstraint. The maximum
// wait constraint limits the time a vehicle can wait between two stops. Wait
// time is defined as the time a vehicle is neither moving nor serving a stop.
func NewMaximumWaitConstraint(
	maximum VehicleTypeDurationExpression,
) (MaximumWaitConstraint, error) {
	connect.Connect(con, &newMaximumWaitConstraint)
	return newMaximumWaitConstraint(maximum)
}
