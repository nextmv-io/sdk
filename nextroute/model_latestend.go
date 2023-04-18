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

	// Latest returns the latest end time expression which defines the latest
	// end of a stop.
	Latest() StopTimeExpression

	// Lateness returns the lateness of a stop. The lateness is the difference
	// between the actual end and its target end time.
	Lateness(stop SolutionStop) float64

	// SetFactor adds a factor with which a deviating stop is multiplied. This
	// is only taken into account if the construct is used as an objective.
	SetFactor(factor float64, stop ModelStop)

	// Factor returns the multiplication factor for the given stop expression.
	Factor(stop ModelStop) float64
}

// NewLatestEnd creates a new latest end construct. The latest end of a stop is
// the latest time a stop can end at the location of the stop. The LatestEnd
// can be added to the model as a constraint or as an objective.
func NewLatestEnd(
	latestEnd StopTimeExpression,
) (LatestEnd, error) {
	connect.Connect(con, &newLatestEnd)
	return newLatestEnd(latestEnd)
}
