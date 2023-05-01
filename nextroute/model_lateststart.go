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
	Latest() StopTimeExpression

	// Lateness returns the lateness of a stop. The lateness is the difference
	// between the actual start and its target start time.
	Lateness(stop SolutionStop) float64

	// SetFactor adds a factor with which a deviating stop is multiplied. This
	// is only taken into account if the construct is used as an objective.
	SetFactor(factor float64, stop ModelStop)

	// Factor returns the multiplication factor for the given stop expression.
	Factor(stop ModelStop) float64
}

// NewLatestStart creates a construct that can be added to the model as a
// constraint or as an objective. The latest start of a stop is the latest time
// a stop can start at the location of the stop.
func NewLatestStart(
	latest StopTimeExpression,
) (LatestStart, error) {
	connect.Connect(con, &newLatestStart)
	return newLatestStart(latest)
}
