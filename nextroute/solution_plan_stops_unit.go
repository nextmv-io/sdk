package nextroute

// SolutionPlanStopsUnit is a set of stops that are planned to be visited by
// a vehicle.
type SolutionPlanStopsUnit interface {
	SolutionPlanUnit
	// ModelPlanStopsUnit returns the [ModelPlanStopsUnit] this unit is
	// based upon.
	ModelPlanStopsUnit() ModelPlanStopsUnit

	// SolutionStop returns the solution stop for the given model stop.
	// Will panic if the stop is not part of the unit.
	SolutionStop(stop ModelStop) SolutionStop
	// SolutionStops returns the solution stops in this unit.
	SolutionStops() SolutionStops
	// StopPositions returns the stop positions of the invoking plan unit.
	// The stop positions are the positions of the stops in the solution.
	// If the unit is unplanned, the stop positions will be empty.
	StopPositions() StopPositions
}

// SolutionPlanStopsUnits is a slice of [SolutionPlanStopsUnit].
type SolutionPlanStopsUnits []SolutionPlanStopsUnit
