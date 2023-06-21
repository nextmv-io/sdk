package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute/common"
)

// MaximumWaitStopConstraint is a constraint that limits the time a vehicle can
// wait between two stops. Wait is defined as the time between arriving at a
// location of a stop and starting (to work),
// [SolutionStop.StartValue()] - [SolutionStop.ArrivalValue()].
type MaximumWaitStopConstraint interface {
	ModelConstraint

	// Maximum returns the maximum expression which defines the maximum time a
	// vehicle can wait at a stop. Returns nil if not set.
	Maximum() StopDurationExpression
}

// NewMaximumWaitStopConstraint returns a new MaximumWaitStopConstraint. The
// maximum wait constraint limits the time a vehicle can wait between two stops.
// Wait is defined as the time between arriving at a location of a stop and
// starting (to work), [SolutionStop.StartValue()] -
// [SolutionStop.ArrivalValue()].
func NewMaximumWaitStopConstraint() (MaximumWaitStopConstraint, error) {
	connect.Connect(con, &newMaximumWaitStopConstraint)
	return newMaximumWaitStopConstraint()
}

// NewCompareModelStopMaximumWaitStopConstraint returns a new CompareFunction
// for the given constraint. The returned function compares two
// model stops by their maximum wait.
func NewCompareModelStopMaximumWaitStopConstraint(
	constraint MaximumWaitStopConstraint,
) common.CompareFunction[ModelStop] {
	return func(a, b ModelStop) int {
		return common.Compare(
			constraint.Maximum().DurationForStop(a).Seconds(),
			constraint.Maximum().DurationForStop(b).Seconds(),
		)
	}
}

// NewCompareSolutionStopMaximumWaitStopConstraint returns a new CompareFunction
// for the given constraint. The returned function compares two
// solution stops by their maximum wait.
func NewCompareSolutionStopMaximumWaitStopConstraint(
	constraint MaximumWaitStopConstraint,
) common.CompareFunction[SolutionStop] {
	return func(a, b SolutionStop) int {
		return common.Compare(
			constraint.Maximum().DurationForStop(a.ModelStop()).Seconds(),
			constraint.Maximum().DurationForStop(b.ModelStop()).Seconds(),
		)
	}
}

// NewCompareModelPlanUnitsMaximumWaitStopConstraint returns a new CompareFunction
// for the given constraint. The returned function compares two plan units
// by the sum of the max wait times.
func NewCompareModelPlanUnitsMaximumWaitStopConstraint(
	constraint MaximumWaitStopConstraint,
) common.CompareFunction[ModelPlanUnit] {
	return func(a, b ModelPlanUnit) int {
		return common.Compare(
			common.SumDefined(a.Stops(), func(t ModelStop) float64 {
				return constraint.Maximum().DurationForStop(t).Seconds()
			}),
			common.SumDefined(b.Stops(), func(t ModelStop) float64 {
				return constraint.Maximum().DurationForStop(t).Seconds()
			}),
		)
	}
}

// NewCompareSolutionPlanUnitsMaximumWaitStopConstraint returns a new CompareFunction
// for the given latest construct. The returned function compares two plan units
// by the sum of the max wait times.
func NewCompareSolutionPlanUnitsMaximumWaitStopConstraint(
	constraint MaximumWaitStopConstraint,
) common.CompareFunction[SolutionPlanUnit] {
	return func(a, b SolutionPlanUnit) int {
		return common.Compare(
			common.SumDefined(a.ModelPlanUnit().Stops(), func(t ModelStop) float64 {
				return constraint.Maximum().DurationForStop(t).Seconds()
			}),
			common.SumDefined(b.ModelPlanUnit().Stops(), func(t ModelStop) float64 {
				return constraint.Maximum().DurationForStop(t).Seconds()
			}),
		)
	}
}
