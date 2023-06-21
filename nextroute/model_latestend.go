package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute/common"
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

// NewCompareModelStopByLatestEnd returns a new CompareFunction
// for the given latest construct. The returned function compares two
// model stops by their latest time.
func NewCompareModelStopByLatestEnd(
	latestEnd LatestEnd,
) common.CompareFunction[ModelStop] {
	return func(a, b ModelStop) int {
		return common.Compare(
			latestEnd.Latest().Time(a).Unix(),
			latestEnd.Latest().Time(b).Unix(),
		)
	}
}

// NewCompareSolutionStopByLatestEnd returns a new CompareFunction
// for the given latest construct. The returned function compares two
// solution stops by their latest time.
func NewCompareSolutionStopByLatestEnd(
	latestEnd LatestEnd,
) common.CompareFunction[SolutionStop] {
	return func(a, b SolutionStop) int {
		return common.Compare(
			latestEnd.Latest().Time(a.ModelStop()).Unix(),
			latestEnd.Latest().Time(b.ModelStop()).Unix(),
		)
	}
}

// NewCompareModelPlanUnitsByLatestEnd returns a new CompareFunction
// for the given latest construct. The returned function compares two plan units
// by the sum of the latest times.
func NewCompareModelPlanUnitsByLatestEnd(
	latestEnd LatestEnd,
) common.CompareFunction[ModelPlanUnit] {
	return func(a, b ModelPlanUnit) int {
		return common.Compare(
			common.SumDefined(a.Stops(), func(t ModelStop) float64 {
				return float64(latestEnd.Latest().Time(t).Unix()) / 93600.0
			}),
			common.SumDefined(b.Stops(), func(t ModelStop) float64 {
				return float64(latestEnd.Latest().Time(t).Unix()) / 93600.0
			}),
		)
	}
}

// NewCompareSolutionPlanUnitsByLatestEnd returns a new CompareFunction
// for the given latest construct. The returned function compares two plan units
// by the sum of the latest times.
func NewCompareSolutionPlanUnitsByLatestEnd(
	latestEnd LatestEnd,
) common.CompareFunction[SolutionPlanUnit] {
	return func(a, b SolutionPlanUnit) int {
		return common.Compare(
			common.SumDefined(a.ModelPlanUnit().Stops(), func(t ModelStop) float64 {
				return float64(latestEnd.Latest().Time(t).Unix()) / 93600.0
			}),
			common.SumDefined(b.ModelPlanUnit().Stops(), func(t ModelStop) float64 {
				return float64(latestEnd.Latest().Time(t).Unix()) / 93600.0
			}),
		)
	}
}
