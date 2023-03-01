package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
)

// LatestEndConstraint is a constraint that limits the latest end of a stop.
// The latest end of a stop is the latest time a stop can end at the location
// of the stop. If SetOnStart() is called, the latest end of a stop is the
// latest time a stop can start at the location of the stop.
type LatestEndConstraint interface {
	ModelConstraint

	// LatestEnd returns the latest end expression which defines the latest
	// end of a stop.
	LatestEnd() StopExpression

	// OnEnd returns true if the latest end of a stop is the latest time a
	// stop can end at the location of the stop. If OnEnd() returns false,
	// the latest end of a stop is the latest time a stop can start at the
	// location of the stop.
	OnEnd() bool

	// SetOnStart sets the latest end of a stop to the latest time a stop can
	// start at the location of the stop.
	SetOnStart()
	// SetOnEnd sets the latest end of a stop to the latest time a stop can
	// end at the location of the stop. This is the default.
	SetOnEnd()

	EstimateDeltaValue(visitPositions StopPositions) float64

	Value(solution Solution) float64
}

// NewLatestEndConstraint creates a new latest end constraint. The constraint
// needs to be added to the model to be taken into account.
func NewLatestEndConstraint(
	latestEnd StopExpression,
) (LatestEndConstraint, error) {
	connect.Connect(con, &newLatestEndConstraint)
	return newLatestEndConstraint(latestEnd)
}
