package nextroute

import "github.com/nextmv-io/sdk/connect"

// MaximumWaitConstraint is a constraint that limits the wait time of a vehicle.
type MaximumWaitConstraint interface {
	ModelConstraint

	// StopMaximum returns the maximum expression which defines the maximum
	// time a vehicle can wait before serving a stop. Returns nil if not set.
	StopMaximum() StopDurationExpression

	// VehicleMaximum returns the maximum expression which defines the maximum
	// accumulated time a vehicle can wait on a route. Returns nil if not set.
	VehicleMaximum() VehicleTypeDurationExpression
}

// NewMaximumWaitConstraint returns a new MaximumWaitConstraint. The maximum
// wait constraint limits the time a vehicle can wait between two stops. Wait
// time is defined as the time a vehicle is neither moving nor serving a stop.
func NewMaximumWaitConstraint() (MaximumWaitConstraint, error) {
	connect.Connect(con, &newMaximumWaitConstraint)
	return newMaximumWaitConstraint()
}
