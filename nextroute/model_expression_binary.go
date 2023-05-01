package nextroute

import (
	"math"

	"github.com/nextmv-io/sdk/connect"
)

// BinaryFunction is a function that takes two float64 values and returns a
// float64 value.
type BinaryFunction func(float64, float64) float64

// BinaryExpression is an expression that takes two expressions as input and
// returns a value.
type BinaryExpression interface {
	ModelExpression
	// Left returns the left expression.
	Left() ModelExpression
	// Right returns the right expression.
	Right() ModelExpression
}

// NewAddExpression returns a new BinaryExpression that adds the values of the
// two expressions.
func NewAddExpression(
	left ModelExpression,
	right ModelExpression,
) BinaryExpression {
	return NewOperatorExpression(
		left,
		right,
		func(left float64, right float64) float64 {
			return left + right
		},
	)
}

// NewMultiplyExpression returns a new BinaryExpression that multiplies the
// values of the two expressions.
func NewMultiplyExpression(
	left ModelExpression,
	right ModelExpression,
) BinaryExpression {
	return NewOperatorExpression(
		left,
		right,
		func(left float64, right float64) float64 {
			return left * right
		},
	)
}

// NewMaximumExpression returns a new BinaryExpression that returns the maximum
// of the values of the two expressions.
func NewMaximumExpression(
	left ModelExpression,
	right ModelExpression,
) BinaryExpression {
	return NewOperatorExpression(
		left,
		right,
		math.Max,
	)
}

// NewMinimumExpression returns a new BinaryExpression that returns the minimum
// of the values of the two expressions.
func NewMinimumExpression(
	left ModelExpression,
	right ModelExpression,
) BinaryExpression {
	return NewOperatorExpression(
		left,
		right,
		math.Min,
	)
}

// NewOperatorExpression returns a new BinaryExpression that uses the given
// operator function.
func NewOperatorExpression(
	left ModelExpression,
	right ModelExpression,
	operator BinaryFunction,
) BinaryExpression {
	connect.Connect(con, &newOperatorExpression)
	return newOperatorExpression(
		left,
		right,
		operator,
	)
}
