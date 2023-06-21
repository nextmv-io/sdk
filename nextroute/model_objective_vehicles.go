package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute/common"
)

// NewVehiclesObjective returns a new VehiclesObjective that uses the number of
// vehicles as an objective. Each vehicle that is not empty is scored by the
// given expression. A vehicle is empty if it has no stops assigned to it
// (except for the first and last visit).
func NewVehiclesObjective(expression VehicleTypeExpression) VehiclesObjective {
	connect.Connect(con, &newVehiclesObjective)
	return newVehiclesObjective(expression)
}

// VehiclesObjective is an objective that uses the number of vehicles as an
// objective. Each vehicle that is not empty is scored by the given expression.
// A vehicle is empty if it has no stops assigned to it (except for the first
// and last visit).
type VehiclesObjective interface {
	ModelObjective
	// ActivationPenalty returns the activation penalty expression.
	ActivationPenalty() VehicleTypeExpression
}

// NewCompareByModelVehiclesObjective returns a new CompareFunction
// for the given objective. The returned function compares two model vehicles by their
// activation penalty.
func NewCompareByModelVehiclesObjective(
	objective VehiclesObjective,
) common.CompareFunction[ModelVehicle] {
	return func(a, b ModelVehicle) int {
		return common.Compare(
			objective.ActivationPenalty().ValueForVehicleType(a.VehicleType()),
			objective.ActivationPenalty().ValueForVehicleType(b.VehicleType()),
		)
	}
}

// NewCompareBySolutionVehiclesObjective returns a new CompareFunction for the given
// objective. The returned function compares two solution vehicles by their
// activation penalty.
func NewCompareBySolutionVehiclesObjective(
	objective VehiclesObjective,
) common.CompareFunction[SolutionVehicle] {
	return func(a, b SolutionVehicle) int {
		return common.Compare(
			objective.ActivationPenalty().ValueForVehicleType(a.ModelVehicle().VehicleType()),
			objective.ActivationPenalty().ValueForVehicleType(b.ModelVehicle().VehicleType()),
		)
	}
}
