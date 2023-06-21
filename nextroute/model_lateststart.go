package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute/common"
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

// NewCompareModelStopByLatestStart returns a new CompareFunction
// for the given latest construct. The returned function compares two
// model stops by their latest time.
func NewCompareModelStopByLatestStart(
	latest LatestStart,
) common.CompareFunction[ModelStop] {
	return func(a, b ModelStop) int {
		return common.Compare(
			latest.Latest().Time(a).Unix(),
			latest.Latest().Time(b).Unix(),
		)
	}
}

// NewCompareSolutionStopByLatestStart returns a new CompareFunction
// for the given latest construct. The returned function compares two
// solution stops by their latest time.
func NewCompareSolutionStopByLatestStart(
	latest LatestStart,
) common.CompareFunction[SolutionStop] {
	return func(a, b SolutionStop) int {
		return common.Compare(
			latest.Latest().Time(a.ModelStop()).Unix(),
			latest.Latest().Time(b.ModelStop()).Unix(),
		)
	}
}

// NewCompareModelPlanUnitsByLatestStart returns a new CompareFunction
// for the given latest construct. The returned function compares two plan units
// by the sum of the latest times.
func NewCompareModelPlanUnitsByLatestStart(
	latest LatestStart,
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

// NewCompareSolutionPlanUnitsByLatestStart returns a new CompareFunction
// for the given latest construct. The returned function compares two plan units
// by the sum of the latest times.
func NewCompareSolutionPlanUnitsByLatestStart(
	latest LatestStart,
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
