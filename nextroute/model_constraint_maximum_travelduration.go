package nextroute

import "github.com/nextmv-io/sdk/connect"

// MaximumTravelDurationConstraint is a constraint that limits the
// total travel duration of a vehicle.
type MaximumTravelDurationConstraint interface {
	ModelConstraint

	// Maximum returns the maximum expression which defines the maximum
	// travel duration of a vehicle type.
	Maximum() VehicleTypeDurationExpression
}

// NewMaximumTravelDurationConstraint creates a new maximum travel duration
// constraint. The constraint needs to be added to the model to be taken into
// account.
func NewMaximumTravelDurationConstraint(
	maximum VehicleTypeDurationExpression,
) (MaximumTravelDurationConstraint, error) {
	connect.Connect(con, &newMaximumTravelDurationConstraint)
	return newMaximumTravelDurationConstraint(maximum)
}
