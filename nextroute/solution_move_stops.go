package nextroute

// Move is a type alias for SolutionMoveStops. It is used to make the
// custom constraints and observers backward compatible.
type Move = SolutionMoveStops

// SolutionMoveStops is a move in a solution. A move is a change in the solution that
// can be executed. A move can be executed which may or may not result in a
// change in the solution. A move can be asked if it is executable, it is not
// executable if it is incomplete or if it is not executable because the move is
// not allowed. A move can be asked if it is an improvement, it is an
// improvement if it is executable and the move has a value less than zero.
// A move describes the change in the solution. The change in the solution
// is described by the stop positions. The stop positions describe for each
// stop in the associated units where in the existing solution the stop
// is supposed to be placed. The stop positions are ordered by the order
// of the stops in the unit. The first stop in the unit is the first
// stop in the stop positions. The last stop in the unit is the last
// stop in the stop positions.
type SolutionMoveStops interface {
	SolutionMove
	// Previous returns previous stop of the first to be planned
	// stop if it would be planned. Previous is the same stop as the
	// previous stop of the first stop position.
	Previous() SolutionStop

	// Next returns the next stop of the last to be planned
	// stop if it would be planned. Next is the same stop as the
	// next stop of the last stop position.
	Next() SolutionStop

	// PlanStopsUnit returns the [SolutionPlanStopsUnit] that is affected by the move.
	PlanStopsUnit() SolutionPlanStopsUnit

	// StopPositions returns the [StopPositions] that define the move and
	// how it will change the solution.
	StopPositions() StopPositions

	// Vehicle returns the vehicle, if known, that is affected by the move. If
	// not known, nil is returned.
	Vehicle() SolutionVehicle

	// Solution returns the solution that is affected by the move.
	Solution() Solution
}

// StopPosition is the definition of the change in the solution for a
// specific stop. The change is defined by a Next and a Stop. The
// Next is a stop which is already part of the solution (it is planned)
// and the Stop is a stop which is not yet part of the solution (it is not
// planned). A stop position states that the stop should be moved from the
// unplanned set to the planned set by positioning it directly before the
// Next.
type StopPosition interface {
	// Previous denotes the upcoming stop's previous stop if the associated move
	// involving the stop position is executed. It's worth noting that
	// the previous stop may not have been planned yet.
	Previous() SolutionStop

	// Next denotes the upcoming stop's next stop if the associated move
	// involving the stop position is executed. It's worth noting that
	// the next stop may not have been planned yet.
	Next() SolutionStop

	// Stop returns the stop which is not yet part of the solution. This stop
	// is not planned yet if the move where the invoking stop position belongs
	// to, has not been executed yet.
	Stop() SolutionStop
}

// StopPositions is a slice of stop positions.
type StopPositions []StopPosition
