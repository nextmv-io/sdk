package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute/common"
)

// NewUnPlannedObjective returns a new UnPlannedObjective that uses the
// un-planned stops as an objective. Each unplanned stop is scored by the
// given expression.
func NewUnPlannedObjective(expression StopExpression) UnPlannedObjective {
	connect.Connect(con, &newUnPlannedObjective)
	return newUnPlannedObjective(expression)
}

// UnPlannedObjective is an objective that uses the un-planned stops as an
// objective. Each unplanned stop is scored by the given expression.
type UnPlannedObjective interface {
	ModelObjective

	// UnplannedPenalty returns the unplanned penalty expression.
	UnplannedPenalty() StopExpression
}

// NewCompareModelStopsByUnPlannedObjective returns a new CompareFunction
// for the given objective. The returned function compares two stop by their
// unplanned penalty.
func NewCompareModelStopsByUnPlannedObjective(
	objective UnPlannedObjective,
) common.CompareFunction[ModelStop] {
	return func(a, b ModelStop) int {
		return common.Compare(
			objective.UnplannedPenalty().ValueForStop(a),
			objective.UnplannedPenalty().ValueForStop(b),
		)
	}
}

// NewCompareSolutionStopsByUnPlannedObjective returns a new CompareFunction
// for the given objective. The returned function compares two solution stops
// by their unplanned penalty.
func NewCompareSolutionStopsByUnPlannedObjective(
	objective UnPlannedObjective,
) common.CompareFunction[SolutionStop] {
	return func(a, b SolutionStop) int {
		return common.Compare(
			objective.UnplannedPenalty().ValueForStop(a.ModelStop()),
			objective.UnplannedPenalty().ValueForStop(b.ModelStop()),
		)
	}
}

// NewCompareModelPlanUnitsByUnPlannedObjective returns a new CompareFunction
// for the given objective. The returned function compares two plan units by
// their unplanned penalty.
func NewCompareModelPlanUnitsByUnPlannedObjective(
	objective UnPlannedObjective,
) common.CompareFunction[ModelPlanUnit] {
	return func(a, b ModelPlanUnit) int {
		return common.Compare(
			common.SumDefined(a.Stops(), func(t ModelStop) float64 {
				return objective.UnplannedPenalty().ValueForStop(t)
			}),
			common.SumDefined(b.Stops(), func(t ModelStop) float64 {
				return objective.UnplannedPenalty().ValueForStop(t)
			}),
		)
	}
}

// NewCompareSolutionPlanUnitsByUnPlannedObjective returns a new CompareFunction
// for the given objective. The returned function compares two plan units by
// their unplanned penalty.
func NewCompareSolutionPlanUnitsByUnPlannedObjective(
	objective UnPlannedObjective,
) common.CompareFunction[SolutionPlanUnit] {
	return func(a, b SolutionPlanUnit) int {
		return common.Compare(
			common.SumDefined(a.ModelPlanUnit().Stops(), func(t ModelStop) float64 {
				return objective.UnplannedPenalty().ValueForStop(t)
			}),
			common.SumDefined(b.ModelPlanUnit().Stops(), func(t ModelStop) float64 {
				return objective.UnplannedPenalty().ValueForStop(t)
			}),
		)
	}
}
