package nextroute

import "github.com/nextmv-io/sdk/connect"

// Maximum can be used as a constraint or an objective that limits the maximum
// cumulative value can be assigned to a vehicle type. The maximum cumulative
// value is defined by the expression and the maximum value is defined by the
// maximum expression.
type Maximum interface {
	ModelConstraint
	ModelObjective

	// Expression returns the expression which is used to calculate the
	// cumulative value of each stop which is required to stay below the
	// maximum value and above zero.
	Expression() ModelExpression

	// Maximum returns the maximum expression which defines the maximum
	// cumulative value that can be assigned to a vehicle type.
	Maximum() VehicleTypeExpression

	// PenaltyOffset returns the penalty offset. Penalty offset is used to
	// offset the penalty. The penalty offset is added to the penalty if there
	// is at least one violation.
	PenaltyOffset() float64

	// SetPenaltyOffset sets the penalty offset. Penalty offset is used to
	// offset the penalty. The penalty offset is added to the penalty if there
	// is at least one violation. The default penalty offset is 0.0 and it can
	// be changed by this method and must be positive.
	SetPenaltyOffset(penaltyOffset float64) error
}

// NewMaximum creates a new maximum constraint/objective. If you add it as a
// constraint, it will behave it as a constraint. If you add it as an objective,
// it will behave as an objective.
func NewMaximum(
	expression ModelExpression,
	maximum VehicleTypeExpression,
) (Maximum, error) {
	connect.Connect(con, &newMaximum)
	return newMaximum(expression, maximum)
}
