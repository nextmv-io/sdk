package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute/common"
)

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

// NewCompareByModelVehicleMaximumStopsConstraint returns a new CompareFunction
// for the given constraint. The returned function compares two model vehicles
// by their maximum stops.
func NewCompareByModelVehicleMaximumStopsConstraint(
	constraint MaximumStopsConstraint,
) common.CompareFunction[ModelVehicle] {
	return func(a, b ModelVehicle) int {
		return common.Compare(
			constraint.MaximumStops().ValueForVehicleType(a.VehicleType()),
			constraint.MaximumStops().ValueForVehicleType(b.VehicleType()),
		)
	}
}

// NewCompareBySolutionVehicleMaximumStopsConstraint returns a new CompareFunction
// for the given constraint. The returned function compares two solution
// vehicles by their maximum stops.
func NewCompareBySolutionVehicleMaximumStopsConstraint(
	constraint MaximumStopsConstraint,
) common.CompareFunction[SolutionVehicle] {
	return func(a, b SolutionVehicle) int {
		return common.Compare(
			constraint.MaximumStops().ValueForVehicleType(a.ModelVehicle().VehicleType()),
			constraint.MaximumStops().ValueForVehicleType(b.ModelVehicle().VehicleType()),
		)
	}
}
