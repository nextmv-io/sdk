package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
)

// LatestStartConstraint is a constraint that limits the latest start of a stop.
// The latest start of a stop is the latest time a stop can start at the location
// of the stop.
type LatestStartConstraint interface {
	ModelConstraint
	ModelObjective

	// Latest returns the latest start expression which defines the latest
	// start of a stop.
	Latest() StopExpression
}

// NewLatestStartConstraint creates a new latest start constraint. The constraint
// needs to be added to the model to be taken into account.
func NewLatestStartConstraint(
	latest StopExpression,
) (LatestStartConstraint, error) {
	connect.Connect(con, &newLatestStartConstraint)
	return newLatestStartConstraint(latest)
}
