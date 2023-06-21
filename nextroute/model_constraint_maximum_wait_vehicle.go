package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute/common"
)

// MaximumWaitVehicleConstraint is a constraint that limits the accumulated time
// a vehicle can wait at stops on its route. Wait is defined as the time between
// arriving at a the location of stop and starting (to work),
// [SolutionStop.StartValue()] - [SolutionStop.ArrivalValue()].
type MaximumWaitVehicleConstraint interface {
	ModelConstraint

	// Maximum returns the maximum expression which defines the maximum
	// accumulated time a vehicle can wait on a route. Returns nil if not set.
	Maximum() VehicleTypeDurationExpression
}

// NewMaximumWaitVehicleConstraint returns a new MaximumWaitVehicleConstraint.
// The maximum wait constraint limits the accumulated time a vehicle can wait at
// stops on its route. Wait is defined as the time between arriving at a
// location of a stop and starting (to work), [SolutionStop.StartValue()] -
// [SolutionStop.ArrivalValue()].
func NewMaximumWaitVehicleConstraint() (MaximumWaitVehicleConstraint, error) {
	connect.Connect(con, &newMaximumWaitVehicleConstraint)
	return newMaximumWaitVehicleConstraint()
}

// NewCompareByModelVehicleMaximumWaitVehicleConstraint returns a new CompareFunction
// for the given constraint. The returned function compares two model vehicles
// by their maximum wait.
func NewCompareByModelVehicleMaximumWaitVehicleConstraint(
	constraint MaximumWaitVehicleConstraint,
) common.CompareFunction[ModelVehicle] {
	return func(a, b ModelVehicle) int {
		return common.Compare(
			constraint.Maximum().ValueForVehicleType(a.VehicleType()),
			constraint.Maximum().ValueForVehicleType(b.VehicleType()),
		)
	}
}

// NewCompareBySolutionVehicleMaximumWaitVehicleConstraint returns a new CompareFunction
// for the given constraint. The returned function compares two solution
// vehicles by their maximum wait.
func NewCompareBySolutionVehicleMaximumWaitVehicleConstraint(
	constraint MaximumWaitVehicleConstraint,
) common.CompareFunction[SolutionVehicle] {
	return func(a, b SolutionVehicle) int {
		return common.Compare(
			constraint.Maximum().ValueForVehicleType(a.ModelVehicle().VehicleType()),
			constraint.Maximum().ValueForVehicleType(b.ModelVehicle().VehicleType()),
		)
	}
}
