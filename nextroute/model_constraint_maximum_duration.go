package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute/common"
)

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

// NewCompareByModelVehicleMaximumDurationConstraint returns a new CompareFunction
// for the given constraint. The returned function compares two model vehicles
// by their maximum duration.
func NewCompareByModelVehicleMaximumDurationConstraint(
	constraint MaximumDurationConstraint,
) common.CompareFunction[ModelVehicle] {
	return func(a, b ModelVehicle) int {
		return common.Compare(
			constraint.Maximum().ValueForVehicleType(a.VehicleType()),
			constraint.Maximum().ValueForVehicleType(b.VehicleType()),
		)
	}
}

// NewCompareBySolutionVehicleMaximumDurationConstraint returns a new CompareFunction
// for the given constraint. The returned function compares two solution
// vehicles by their maximum duration.
func NewCompareBySolutionVehicleMaximumDurationConstraint(
	constraint MaximumDurationConstraint,
) common.CompareFunction[SolutionVehicle] {
	return func(a, b SolutionVehicle) int {
		return common.Compare(
			constraint.Maximum().ValueForVehicleType(a.ModelVehicle().VehicleType()),
			constraint.Maximum().ValueForVehicleType(b.ModelVehicle().VehicleType()),
		)
	}
}
