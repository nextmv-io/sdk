package nextroute

// ModelPlanUnitsUnit is a set of plan units. A plan unit is a set of stops
// that must be visited together.
type ModelPlanUnitsUnit interface {
	ModelPlanUnit

	// PlanStopsUnits returns the plan units in the invoking unit.
	PlanStopsUnits() ModelPlanStopsUnits

	// IsDisjunction returns true if the plan unit is a disjunction.
	// A disjunction is a set of plan units of which exactly one must be
	// planned.
	IsDisjunction() bool

	// IsConjunction returns true if the plan unit is a conjunction.
	// A conjunction is a set of plan units of which all must be planned.
	IsConjunction() bool

	// SameVehicle returns true if all the plan units in this unit have to be
	// planned on the same vehicle. If this unit is a conjunction, then
	// this will return true if all the plan units in this unit have to be
	// planned on the same vehicle. If this unit is a disjunction, then
	// this has no semantic meaning.
	SameVehicle() bool
}

// ModelPlanUnitsUnits is a slice of model plan units units .
type ModelPlanUnitsUnits []ModelPlanUnitsUnit
