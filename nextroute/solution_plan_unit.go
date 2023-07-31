package nextroute

// SolutionPlanUnit is a set of stops that are planned to be visited by
// a vehicle.
type SolutionPlanUnit interface {
	// IsFixed returns true if any of stops are fixed.
	IsFixed() bool

	// IsPlanned returns true if all the stops are planned.
	IsPlanned() bool

	// ModelPlanUnit returns the model plan unit associated with the
	// solution plan unit.
	ModelPlanUnit() ModelPlanUnit

	// Solution returns the solution this unit is part of.
	Solution() Solution

	// UnPlan un-plans the unit by removing the underlying solution stops
	// from the solution. Returns true if the unit was unplanned
	// successfully, false if the unit was not unplanned successfully. A
	// unit is not successful if it did not result in a change in the
	// solution without violating any hard constraints.
	UnPlan() (bool, error)
}

// SolutionPlanUnits is a slice of [SolutionPlanUnit].
type SolutionPlanUnits []SolutionPlanUnit
