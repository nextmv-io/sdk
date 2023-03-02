package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
)

// LatestEnd is a construct that can be added to the model as a constraint or
// as an objective. The latest end of a stop is the latest time a stop can end
// at the location of the stop.
type LatestEnd interface {
	ConstraintReporter
	ModelConstraint
	ModelObjective

	// Latest returns the latest end expression which defines the latest
	// end of a stop.
	Latest() StopExpression
}

// NewLatestEnd creates a new latest end construct. The latest end of a stop is
// the latest time a stop can end at the location of the stop. The LatestEnd
// can be added to the model as a constraint or as an objective.
func NewLatestEnd(
	latestEnd StopExpression,
) (LatestEnd, error) {
	connect.Connect(con, &newLatestEnd)
	return newLatestEnd(latestEnd)
}
