package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/measure"
)

// NewMeasureByIndexExpression returns a new MeasureByIndexExpression.
// A MeasureByIndexExpression is a ModelExpression that uses a measure.ByIndex to
// calculate the cost between two stops.
// The index of the measure have to be the same as the index of the stops in the
// model.
func NewMeasureByIndexExpression(m measure.ByIndex) ModelExpression {
	connect.Connect(con, &newMeasureByIndexExpression)
	return newMeasureByIndexExpression(m)
}
