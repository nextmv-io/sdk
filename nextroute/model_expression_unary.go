package nextroute

import "github.com/nextmv-io/sdk/connect"

// NewTermExpression returns a new TermExpression that calculates the product of
// the given factor and the value of the given expression.
func NewTermExpression(
	factor float64,
	expression ModelExpression,
) TermExpression {
	connect.Connect(con, &newTermExpression)
	return newTermExpression(factor, expression)
}

// TermExpression is an expression that returns the product of the given factor
// and the value of the given expression.
type TermExpression interface {
	ModelExpression

	// Expression returns the expression.
	Expression() ModelExpression

	// Factor returns the factor.
	Factor() float64
}
