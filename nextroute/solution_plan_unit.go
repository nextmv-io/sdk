package nextroute

// SolutionPlanUnit is a set of stops that are planned to be visited by
// a vehicle.
type SolutionPlanUnit interface {
	// IsPlanned returns true if all the stops are planned.
	IsPlanned() bool

	// ModelPlanUnit returns the [ModelPlanUnit] this unit is
	// based upon.
	ModelPlanUnit() ModelPlanUnit

	// Solution returns the solution this unit is part of.
	Solution() Solution
	// SolutionStop returns the solution stop for the given model stop.
	// Will panic if the stop is not part of the unit.
	SolutionStop(stop ModelStop) SolutionStop
	// SolutionStops returns the solution stops in this unit.
	SolutionStops() SolutionStops
	// StopPositions returns the stop positions of the invoking plan unit.
	// The stop positions are the positions of the stops in the solution.
	// If the unit is unplanned, the stop positions will be empty.
	StopPositions() StopPositions

	// UnPlan un-plans the unit by removing the underlying solution stops
	// from the solution. Returns true if the unit was unplanned
	// successfully, false if the unit was not unplanned successfully. A
	// unit is not successful if it did not result in a change in the
	// solution without violating any hard constraints.
	UnPlan() (bool, error)
}

// SolutionPlanUnits is a slice of [SolutionPlanUnit].
type SolutionPlanUnits []SolutionPlanUnit
