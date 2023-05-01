package nextroute

import "github.com/nextmv-io/sdk/connect"

// StopPositionsHint is an interface that can be used to give a hint to the
// solver about the next stop position. This can be used to speed up the
// solver. The solver will use the hint if it is available. Hints are generated
// by the estimate function of a constraint.
type StopPositionsHint interface {
	// HasNextStopPositions returns true if the hint contains next positions.
	HasNextStopPositions() bool

	// NextStopPositions returns the next positions.
	NextStopPositions() StopPositions

	// SkipVehicle returns true if the solver should skip the vehicle. The
	// solver will use the hint if it is available.
	SkipVehicle() bool
}

// NoPositionsHint returns a new StopPositionsHint that does not skip
// the vehicle and does not contain a next stop. The solver will try to find
// the next stop.
func NoPositionsHint() StopPositionsHint {
	connect.Connect(con, &noPositionsHint)
	return noPositionsHint()
}

// SkipVehiclePositionsHint returns a new StopPositionsHint that skips the
// vehicle.
func SkipVehiclePositionsHint() StopPositionsHint {
	connect.Connect(con, &skipVehiclePositionsHint)
	return skipVehiclePositionsHint()
}
