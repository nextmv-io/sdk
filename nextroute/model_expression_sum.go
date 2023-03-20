package nextroute

import "github.com/nextmv-io/sdk/connect"

// NewSumExpression returns a new SumExpression that calculates the sum of the
// values of the given expressions.
func NewSumExpression(expressions ModelExpressions) SumExpression {
	connect.Connect(con, &newSumExpression)
	return newSumExpression(expressions)
}

// SumExpression is an expression that returns the sum of the values of the
// given expressions.
type SumExpression interface {
	ModelExpression

	// AddExpression adds an expression to the sum.
	AddExpression(expression ModelExpression)

	// Expressions returns the expressions that are part of the sum.
	Expressions() ModelExpressions
}
