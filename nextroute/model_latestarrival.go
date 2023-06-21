package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute/common"
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
// a stop can arrive at the location of the stop.
func NewLatestArrival(
	latest StopTimeExpression,
) (LatestArrival, error) {
	connect.Connect(con, &newLatestArrival)
	return newLatestArrival(latest)
}

// NewCompareModelStopByLatestArrival returns a new CompareFunction
// for the given latest construct. The returned function compares two
// model stops by their latest time.
func NewCompareModelStopByLatestArrival(
	latest LatestArrival,
) common.CompareFunction[ModelStop] {
	return func(a, b ModelStop) int {
		return common.Compare(
			latest.Latest().Time(a).Unix(),
			latest.Latest().Time(b).Unix(),
		)
	}
}

// NewCompareSolutionStopByLatestArrival returns a new CompareFunction
// for the given latest construct. The returned function compares two
// solution stops by their latest time.
func NewCompareSolutionStopByLatestArrival(
	latest LatestArrival,
) common.CompareFunction[SolutionStop] {
	return func(a, b SolutionStop) int {
		return common.Compare(
			latest.Latest().Time(a.ModelStop()).Unix(),
			latest.Latest().Time(b.ModelStop()).Unix(),
		)
	}
}

// NewCompareModelPlanUnitsByLatestArrival returns a new CompareFunction
// for the given latest construct. The returned function compares two plan units
// by the sum of the latest times.
func NewCompareModelPlanUnitsByLatestArrival(
	latest LatestArrival,
) common.CompareFunction[ModelPlanUnit] {
	return func(a, b ModelPlanUnit) int {
		return common.Compare(
			common.SumDefined(a.Stops(), func(t ModelStop) float64 {
				return float64(latest.Latest().Time(t).Unix()) / 93600.0
			}),
			common.SumDefined(b.Stops(), func(t ModelStop) float64 {
				return float64(latest.Latest().Time(t).Unix()) / 93600.0
			}),
		)
	}
}

// NewCompareSolutionPlanUnitsByLatestArrival returns a new CompareFunction
// for the given latest construct. The returned function compares two plan units
// by the sum of the latest times.
func NewCompareSolutionPlanUnitsByLatestArrival(
	latest LatestArrival,
) common.CompareFunction[SolutionPlanUnit] {
	return func(a, b SolutionPlanUnit) int {
		return common.Compare(
			common.SumDefined(a.ModelPlanUnit().Stops(), func(t ModelStop) float64 {
				return float64(latest.Latest().Time(t).Unix()) / 93600.0
			}),
			common.SumDefined(b.ModelPlanUnit().Stops(), func(t ModelStop) float64 {
				return float64(latest.Latest().Time(t).Unix()) / 93600.0
			}),
		)
	}
}
