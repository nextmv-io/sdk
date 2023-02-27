package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
)

// NewHaversineExpression returns a new DistanceExpression that calculates the
// distance between two stops using the Haversine formula.
func NewHaversineExpression(buffer bool) DistanceExpression {
	connect.Connect(con, &newHaversineExpression)
	return newHaversineExpression(buffer)
}
