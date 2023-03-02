package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
)

// LatestStart is a construct that can be added to the model as a constraint or
// as an objective. The latest start of a stop is the latest time a stop can
// start at the location of the stop.
type LatestStart interface {
	ConstraintReporter
	ModelConstraint
	ModelObjective

	// Latest returns the latest start expression which defines the latest
	// start of a stop.
	Latest() StopExpression
}

// NewLatestStart creates a construct that can be added to the model as a
// constraint or as an objective. The latest start of a stop is the latest time
// a stop can start at the location of the stop.
func NewLatestStart(
	latest StopExpression,
) (LatestStart, error) {
	connect.Connect(con, &newLatestStart)
	return newLatestStart(latest)
}
