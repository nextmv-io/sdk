package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute/common"
)

// AttributesConstraint is a constraint that limits the vehicles a plan unit
// can be added to. The Attribute constraint configures compatibility
// attributes for stops and vehicles separately. This is done by specifying
// a list of attributes for stops and vehicles, respectively. Stops that
// have configured attributes are only compatible with vehicles that match
// at least one of them. Stops that do not have any specified attributes are
// compatible with any vehicle. Vehicles that do not have any specified
// attributes are only compatible with stops without attributes.
type AttributesConstraint interface {
	ModelConstraint

	// SetStopAttributes sets the attributes for the given stop. The attributes
	// are specified as a list of strings. The attributes are not interpreted
	// in any way. They are only used to determine compatibility between stops
	// and vehicle types.
	SetStopAttributes(
		stop ModelStop,
		stopAttributes []string,
	)
	// SetVehicleTypeAttributes sets the attributes for the given vehicle type.
	// The attributes are specified as a list of strings. The attributes are not
	// interpreted in any way. They are only used to determine compatibility
	// between stops and vehicle types.
	SetVehicleTypeAttributes(
		vehicle ModelVehicleType,
		vehicleAttributes []string,
	)
	// StopAttributes returns the attributes for the given stop. The attributes
	// are specified as a list of strings.
	StopAttributes(stop ModelStop) []string

	// VehicleTypeAttributes returns the attributes for the given vehicle type.
	// The attributes are specified as a list of strings.
	VehicleTypeAttributes(vehicle ModelVehicleType) []string
}

// NewAttributesConstraint creates a new attributes constraint. The constraint
// needs to be added to the model to be taken into account.
func NewAttributesConstraint() (AttributesConstraint, error) {
	connect.Connect(con, &newAttributesConstraint)
	return newAttributesConstraint()
}

// NewCompareModelStopCountAttributesConstraint returns a new CompareFunction
// for the given constraint. The returned function compares two
// model stops by the count of attributes.
func NewCompareModelStopCountAttributesConstraint(
	constraint AttributesConstraint,
) common.CompareFunction[ModelStop] {
	return func(a, b ModelStop) int {
		return common.Compare(
			len(constraint.StopAttributes(a)),
			len(constraint.StopAttributes(b)),
		)
	}
}

// NewCompareSolutionStopCountAttributesConstraint returns a new CompareFunction
// for the given constraint. The returned function compares two
// solution stops by the count of attributes.
func NewCompareSolutionStopCountAttributesConstraint(
	constraint AttributesConstraint,
) common.CompareFunction[SolutionStop] {
	return func(a, b SolutionStop) int {
		return common.Compare(
			len(constraint.StopAttributes(a.ModelStop())),
			len(constraint.StopAttributes(b.ModelStop())),
		)
	}
}

// NewCompareModelPlanUnitCountAttributesConstraint returns a new CompareFunction
// for the given constraint. The returned function compares two
// model plan units by the sum of the count of attributes of the stops of the
// unit.
func NewCompareModelPlanUnitCountAttributesConstraint(
	constraint AttributesConstraint,
) common.CompareFunction[ModelPlanUnit] {
	return func(a, b ModelPlanUnit) int {
		return common.Compare(
			common.SumDefined(a.Stops(), func(t ModelStop) int {
				return len(constraint.StopAttributes(t))
			}),
			common.SumDefined(b.Stops(), func(t ModelStop) int {
				return len(constraint.StopAttributes(t))
			}),
		)
	}
}

// NewCompareSolutionPlanUnitCountAttributesConstraint returns a new CompareFunction
// for the given constraint. The returned function compares two
// solution plan units by the sum of the count of attributes of the stops of the
// unit.
func NewCompareSolutionPlanUnitCountAttributesConstraint(
	constraint AttributesConstraint,
) common.CompareFunction[SolutionPlanUnit] {
	return func(a, b SolutionPlanUnit) int {
		return common.Compare(
			common.SumDefined(a.ModelPlanUnit().Stops(), func(t ModelStop) int {
				return len(constraint.StopAttributes(t))
			}),
			common.SumDefined(b.ModelPlanUnit().Stops(), func(t ModelStop) int {
				return len(constraint.StopAttributes(t))
			}),
		)
	}
}

// NewCompareByModelVehicleAttributesConstraint returns a new CompareFunction
// for the given constraint. The returned function compares two model vehicles
// by their number of attributes.
func NewCompareByModelVehicleAttributesConstraint(
	constraint AttributesConstraint,
) common.CompareFunction[ModelVehicle] {
	return func(a, b ModelVehicle) int {
		return common.Compare(
			len(constraint.VehicleTypeAttributes(a.VehicleType())),
			len(constraint.VehicleTypeAttributes(b.VehicleType())),
		)
	}
}

// NewCompareBySolutionVehicleAttributesConstraint returns a new CompareFunction
// for the given constraint. The returned function compares two solution
// vehicles  their number of attributes.
func NewCompareBySolutionVehicleAttributesConstraint(
	constraint AttributesConstraint,
) common.CompareFunction[SolutionVehicle] {
	return func(a, b SolutionVehicle) int {
		return common.Compare(
			len(constraint.VehicleTypeAttributes(a.ModelVehicle().VehicleType())),
			len(constraint.VehicleTypeAttributes(b.ModelVehicle().VehicleType())),
		)
	}
}
