package nextroute

// SolutionPlanUnitsUnit is a set of solution plan units.
type SolutionPlanUnitsUnit interface {
	SolutionPlanUnit
	// ModelPlanUnitsUnit returns the [ModelPlanUnitsUnit] this unit is
	// based upon.
	ModelPlanUnitsUnit() ModelPlanUnitsUnit

	// SolutionPlanUnit returns the solution plan unit for the given model
	// plan unit. Will panic if the unit is not part of the unit.
	SolutionPlanUnit(planUnit ModelPlanUnit) SolutionPlanUnit
	// SolutionPlanUnits returns the solution units in this unit.
	SolutionPlanUnits() SolutionPlanUnits
}

// SolutionPlanUnitsUnits is a slice of [SolutionPlanUnitsUnit].
type SolutionPlanUnitsUnits []SolutionPlanUnitsUnit
