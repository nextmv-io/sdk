package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute/common"
)

// MaximumConstraint is a constraint that limits the maximum cumulative
// value can be assigned to a vehicle type. The maximum cumulative value is
// defined by the expression and the maximum value is defined by the
// maximum expression.
type MaximumConstraint interface {
	ModelConstraint

	// Expression returns the expression which is used to calculate the
	// cumulative value of each stop which is required to stay below the
	// maximum value and above zero.
	Expression() ModelExpression

	// Maximum returns the maximum expression which defines the maximum
	// cumulative value that can be assigned to a vehicle type.
	Maximum() VehicleTypeExpression
}

// NewMaximumConstraint creates a new maximum constraint. The constraint
// needs to be added to the model to be taken into account.
func NewMaximumConstraint(
	expression ModelExpression,
	maximum VehicleTypeExpression,
) (MaximumConstraint, error) {
	connect.Connect(con, &newMaximumConstraint)
	return newMaximumConstraint(expression, maximum)
}

// NewCompareByModelVehicleMaximumConstraint returns a new CompareFunction
// for the given constraint. The returned function compares two model vehicles
// by their maximum value.
func NewCompareByModelVehicleMaximumConstraint(
	constraint MaximumConstraint,
) common.CompareFunction[ModelVehicle] {
	return func(a, b ModelVehicle) int {
		return common.Compare(
			constraint.Maximum().ValueForVehicleType(a.VehicleType()),
			constraint.Maximum().ValueForVehicleType(b.VehicleType()),
		)
	}
}

// NewCompareBySolutionVehicleMaximumConstraint returns a new CompareFunction
// for the given constraint. The returned function compares two solution
// vehicles by their maximum value.
func NewCompareBySolutionVehicleMaximumConstraint(
	constraint MaximumConstraint,
) common.CompareFunction[SolutionVehicle] {
	return func(a, b SolutionVehicle) int {
		return common.Compare(
			constraint.Maximum().ValueForVehicleType(a.ModelVehicle().VehicleType()),
			constraint.Maximum().ValueForVehicleType(b.ModelVehicle().VehicleType()),
		)
	}
}
