package nextroute

import "github.com/nextmv-io/sdk/connect"

// StopPositionsHint is an interface that can be used to give a hint to the
// solver about the next stop position. This can be used to speed up the
// solver. The solver will use the hint if it is available. Hints are generated
// by the estimate function of a constraint.
type StopPositionsHint interface {
	// HasNextStop returns true if the hint contains a next stop. If the hint
	// does not contain a next stop the solver will try to find the next stop.
	HasNextStop() bool

	// NextStop returns the next stop. The stop must be part of the solution.
	// The solver will use the hint if HasNextStop returns true.
	NextStop() SolutionStop

	// SkipVehicle returns true if the solver should skip the vehicle. The
	// solver will use the hint if it is available.
	SkipVehicle() bool
}

// NewNoStopPositionsHint returns a new StopPositionsHint that does not skip
// the vehicle and does not contain a next stop. The solver will try to find
// the next stop.
func NewNoStopPositionsHint() StopPositionsHint {
	connect.Connect(con, &newNoStopPositionsHint)
	return newNoStopPositionsHint()
}

// NewSkipVehiclePositionsHint returns a new StopPositionsHint that skips the
// vehicle if skipVehicle is true. Is skipVehicle is false the solver will try
// to find the next stop.
func NewSkipVehiclePositionsHint(skipVehicle bool) StopPositionsHint {
	connect.Connect(con, &newSkipVehiclePositionsHint)
	return newSkipVehiclePositionsHint(skipVehicle)
}
