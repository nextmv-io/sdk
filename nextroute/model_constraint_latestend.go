package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
)

// LatestEndConstraint is a constraint that limits the latest end of a stop.
// The latest end of a stop is the latest time a stop can end at the location
// of the stop.
type LatestEndConstraint interface {
	ModelConstraint
	ModelObjective

	// Latest returns the latest end expression which defines the latest
	// end of a stop.
	Latest() StopExpression
}

// NewLatestEndConstraint creates a new latest end constraint. The constraint
// needs to be added to the model to be taken into account.
func NewLatestEndConstraint(
	latestEnd StopExpression,
) (LatestEndConstraint, error) {
	connect.Connect(con, &newLatestEndConstraint)
	return newLatestEndConstraint(latestEnd)
}
