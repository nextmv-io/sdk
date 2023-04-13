package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/measure"
)

// NewMeasureByPointExpression returns a new MeasureByPointExpression.
// A MeasureByPointExpression is a ModelExpression that uses a measure.ByPoint
// to calculate the cost between two stops.
func NewMeasureByPointExpression(m measure.ByPoint) ModelExpression {
	connect.Connect(con, &newMeasureByPointExpression)
	return newMeasureByPointExpression(m)
}
