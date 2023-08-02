package nextroute

// ModelPlanUnitsUnit is a set of plan units. A plan unit is a set of stops
// that must be visited together.
type ModelPlanUnitsUnit interface {
	ModelPlanUnit

	// PlanUnits returns the plan units in the invoking unit.
	PlanUnits() ModelPlanUnits

	// PlanOneOf returns true if the plan unit only has to plan exactly one of
	// the associated plan units. If PlanOneOf returns true, then PlanAll will
	// return false and vice versa.
	PlanOneOf() bool

	// PlanAll returns true if the plan unit has to plan all the associated
	// plan units. If PlanAll returns true, then PlanOneOf will return false
	// and vice versa.
	PlanAll() bool

	// SameVehicle returns true if all the plan units in this unit have to be
	// planned on the same vehicle. If this unit is a conjunction, then
	// this will return true if all the plan units in this unit have to be
	// planned on the same vehicle. If this unit is a disjunction, then
	// this has no semantic meaning.
	SameVehicle() bool
}

// ModelPlanUnitsUnits is a slice of model plan units units .
type ModelPlanUnitsUnits []ModelPlanUnitsUnit
