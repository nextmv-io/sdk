package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
)

// TemporalReference is a representation of OnArrival, OnEnd or OnStart as an
// enum.
type TemporalReference int

const (
	// OnStart refers to the Start at a stop.
	OnStart TemporalReference = iota
	// OnEnd refers to the End at a stop.
	OnEnd
	// OnArrival refers to the Arrival at a stop.
	OnArrival = 2
)

// EarlinessObjective is a construct that can be added to the model as an
// objective. It uses to the difference of Arrival, Start or End to the target
// time to penalize.
type EarlinessObjective interface {
	ModelObjective

	// TargetTime returns the target time expression which defines target time
	// that is compared to either arrival, start or end at the stop - depending
	// on the given TemporalReference.
	TargetTime() StopTimeExpression

	// Earliness returns the earliness of a stop. The earliness is the
	// difference between target time and the actual arrival, start or stop of a
	// stop. Depending on the TemporalReference.
	Earliness(stop SolutionStop) float64

	// TemporalReference represents the arrival, start or stop.
	TemporalReference() TemporalReference
}

// NewEarlinessObjective creates a construct that can be added to the model as
// an objective. The target time of a stop is the time the vehicle should either
// arrive, start or end. The latter is defined by temporal reference.
func NewEarlinessObjective(
	targetTime StopTimeExpression,
	temporalReference TemporalReference,
) (EarlinessObjective, error) {
	connect.Connect(con, &newEarlinessObjective)
	return newEarlinessObjective(targetTime, temporalReference)
}
