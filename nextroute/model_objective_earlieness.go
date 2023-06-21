package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute/common"
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
// arrive, start or end. The latter is defined by temporal reference. The
// earliness factor is used to apply a penalty factor per stop to the deviating
// seconds.
func NewEarlinessObjective(
	targetTime StopTimeExpression,
	earlinessFactor StopExpression,
	temporalReference TemporalReference,
) (EarlinessObjective, error) {
	connect.Connect(con, &newEarlinessObjective)
	return newEarlinessObjective(targetTime, earlinessFactor, temporalReference)
}

// NewCompareModelStopByEarlinessObjective returns a new CompareFunction
// for the given objective. The returned function compares two model stops by
// their target time.
func NewCompareModelStopByEarlinessObjective(
	objective EarlinessObjective,
) common.CompareFunction[ModelStop] {
	return func(a, b ModelStop) int {
		return common.Compare(
			objective.TargetTime().Time(a).Unix(),
			objective.TargetTime().Time(b).Unix(),
		)
	}
}

// NewCompareSolutionStopByEarlinessObjective returns a new CompareFunction
// for the given objective. The returned function compares two solution stops by
// their target time.
func NewCompareSolutionStopByEarlinessObjective(
	objective EarlinessObjective,
) common.CompareFunction[SolutionStop] {
	return func(a, b SolutionStop) int {
		return common.Compare(
			objective.TargetTime().Time(a.ModelStop()).Unix(),
			objective.TargetTime().Time(b.ModelStop()).Unix(),
		)
	}
}

// NewCompareModelPlanUnitsByEarlinessObjective returns a new CompareFunction
// for the given objective. The returned function compares two plan units by
// the sum of the target times.
func NewCompareModelPlanUnitsByEarlinessObjective(
	objective EarlinessObjective,
) common.CompareFunction[ModelPlanUnit] {
	return func(a, b ModelPlanUnit) int {
		return common.Compare(
			common.SumDefined(a.Stops(), func(t ModelStop) float64 {
				return float64(objective.TargetTime().Time(t).Unix()) / 93600.0
			}),
			common.SumDefined(b.Stops(), func(t ModelStop) float64 {
				return float64(objective.TargetTime().Time(t).Unix()) / 93600.0
			}),
		)
	}
}

// NewCompareSolutionPlanUnitsByEarlinessObjective returns a new CompareFunction
// for the given objective. The returned function compares two plan units by
// the sum of the target times.
func NewCompareSolutionPlanUnitsByEarlinessObjective(
	objective EarlinessObjective,
) common.CompareFunction[SolutionPlanUnit] {
	return func(a, b SolutionPlanUnit) int {
		return common.Compare(
			common.SumDefined(a.ModelPlanUnit().Stops(), func(t ModelStop) float64 {
				return float64(objective.TargetTime().Time(t).Unix()) / 93600.0
			}),
			common.SumDefined(b.ModelPlanUnit().Stops(), func(t ModelStop) float64 {
				return float64(objective.TargetTime().Time(t).Unix()) / 93600.0
			}),
		)
	}
}
