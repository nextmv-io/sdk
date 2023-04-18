package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
)

// LatestArrival is a construct that can be added to the model as a constraint
// or as an objective. The latest arrival of a stop is the latest time a stop
// can arrive at the location of the stop.
type LatestArrival interface {
	ConstraintReporter
	ModelConstraint
	ModelObjective

	// Latest returns the latest arrival expression which defines the latest
	// arrival of a stop.
	Latest() StopTimeExpression

	// Lateness returns the lateness of a stop. The lateness is the difference
	// between the actual arrival and its target arrival time.
	Lateness(stop SolutionStop) float64

	// SetFactor adds a factor with which a deviating stop is multiplied. This
	// is only taken into account if the construct is used as an objective.
	SetFactor(factor float64, stop ModelStop)

	// Factor returns the multiplication factor for the given stop expression.
	Factor(stop ModelStop) float64
}

// NewLatestArrival creates a construct that can be added to the model as a
// constraint or as an objective. The latest start of a stop is the latest time
// a stop can start at the location of the stop.
func NewLatestArrival(
	latest StopTimeExpression,
) (LatestArrival, error) {
	connect.Connect(con, &newLatestArrival)
	return newLatestArrival(latest)
}
