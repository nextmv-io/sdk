package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
)

// SolutionFilterComposition is a SolutionMoveFilter composed of other
// SolutionMoveFilters. It is a SolutionMoveFilter itself and filters moves
// that are filtered by any of the composed filters.
type SolutionFilterComposition interface {
	// FilterMove returns true if the move should be filtered out.
	FilterMove(move Move) bool

	// FilterVehicle returns true no moves should be generated for
	// vehicle.
	FilterVehicle(vehicle SolutionVehicle) bool

	SolutionMoveFilters() SolutionMoveFilters

	SolutionVehicleFilters() SolutionVehicleFilters
}

// NewSolutionFilterComposition returns a new SolutionFilterComposition. It
// takes a slice of potential SolutionMoveFilter and SolutionVehicleFilter
// instances. It will only use the filters that are actually of the correct
// type.
func NewSolutionFilterComposition(filters ...interface{}) SolutionFilterComposition {
	connect.Connect(con, &newSolutionFilterComposition)
	return newSolutionFilterComposition(filters...)
}
