package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute/common"
)

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

// NewCompareByModelVehicleMaximumTravelDurationConstraint returns a new CompareFunction
// for the given constraint. The returned function compares two model vehicles
// by their maximum travel duration.
func NewCompareByModelVehicleMaximumTravelDurationConstraint(
	constraint MaximumTravelDurationConstraint,
) common.CompareFunction[ModelVehicle] {
	return func(a, b ModelVehicle) int {
		return common.Compare(
			constraint.Maximum().ValueForVehicleType(a.VehicleType()),
			constraint.Maximum().ValueForVehicleType(b.VehicleType()),
		)
	}
}

// NewCompareBySolutionVehicleMaximumTravelDurationConstraint returns a new CompareFunction
// for the given constraint. The returned function compares two solution
// vehicles by their maximum travel duration.
func NewCompareBySolutionVehicleMaximumTravelDurationConstraint(
	constraint MaximumTravelDurationConstraint,
) common.CompareFunction[SolutionVehicle] {
	return func(a, b SolutionVehicle) int {
		return common.Compare(
			constraint.Maximum().ValueForVehicleType(a.ModelVehicle().VehicleType()),
			constraint.Maximum().ValueForVehicleType(b.ModelVehicle().VehicleType()),
		)
	}
}
