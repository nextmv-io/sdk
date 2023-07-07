package nextroute

import "github.com/nextmv-io/sdk/connect"

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
// needs to be added to the model to be taken into account. Furthermore, a name
// identifying the constraint needs to be provided.
func NewMaximumConstraint(
	expression ModelExpression,
	maximum VehicleTypeExpression,
	name string,
) (MaximumConstraint, error) {
	connect.Connect(con, &newMaximumConstraint)
	return newMaximumConstraint(expression, maximum, name)
}
