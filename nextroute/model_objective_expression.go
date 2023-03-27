package nextroute

import "github.com/nextmv-io/sdk/connect"

// NewExpressionObjective returns a new ExpressionObjective inteface.
func NewExpressionObjective(e ModelExpression) ExpressionObjective {
	connect.Connect(con, &newExpressionObjective)
	return newExpressionObjective(e)
}

// ExpressionObjective is an objective that uses an expression to calculate an
// objective.
type ExpressionObjective interface {
	ModelObjective

	// Expression returns the expression that is used to calculate the
	// objective.
	Expression() ModelExpression
}
