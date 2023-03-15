package nextroute

import "github.com/nextmv-io/sdk/connect"

// MaximumDurationConstraint is a constraint that limits the
// duration of a vehicle.
type MaximumDurationConstraint interface {
	ModelConstraint

	// Maximum returns the maximum expression which defines the maximum
	// duration of a vehicle type.
	Maximum() VehicleTypeDurationExpression
}

// NewMaximumDurationConstraint creates a new maximum constraint. The constraint
// needs to be added to the model to be taken into account.
func NewMaximumDurationConstraint(
	maximum VehicleTypeDurationExpression,
) (MaximumDurationConstraint, error) {
	connect.Connect(con, &newMaximumDurationConstraint)
	return newMaximumDurationConstraint(maximum)
}
