package nextroute

import (
	"context"
	"time"

	"github.com/nextmv-io/sdk/nextroute/common"
)

// SolutionVehicle is a vehicle in a solution.
type SolutionVehicle interface {
	// FirstMove creates a move that adds the given plan cluster to the
	// vehicle after the first solution stop of the vehicle. The move is
	// first feasible move after the first solution stop based on the
	// estimates of the constraint, this move is not necessarily executable.
	FirstMove(SolutionPlanCluster) Move

	// BestMove returns the best move for the given solution plan cluster on
	// the invoking vehicle. The best move is the move that has the lowest
	// score. If there are no moves available for the given solution plan
	// cluster, a move is returned which is not executable, Move.IsExecutable.
	BestMove(context.Context, SolutionPlanCluster) Move

	// Centroid returns the centroid of the vehicle. The centroid is the
	// average location of all stops in the vehicle excluding the start and
	// end stops.
	Centroid() common.Location

	// EndTime returns the end time of the vehicle. The end time is the time
	// the vehicle ends at the end stop.
	EndTime() time.Time

	// First returns the first stop of the vehicle. The first stop is the
	// start stop.
	First() SolutionStop

	// Index returns the index of the vehicle in the solution.
	Index() int
	// IsEmpty returns true if the vehicle is empty, false otherwise. A
	// vehicle is empty if it does not have any stops. The start and end
	// stops are not considered.
	IsEmpty() bool

	// Last returns the last stop of the vehicle. The last stop is the end
	// stop.
	Last() SolutionStop

	// NumberOfStops returns the number of stops in the vehicle. The start
	// and end stops are not considered.
	NumberOfStops() int

	// SolutionStops returns the stops in the vehicle. The start and end
	// stops are included in the returned stops.
	SolutionStops() SolutionStops
	// StartTime returns the start time of the vehicle. The start time is
	// the time the vehicle starts at the start stop, it has been set
	// in the factory method of the vehicle Solution.NewVehicle.
	StartTime() time.Time

	// ModelVehicle returns the modeled vehicle type of the vehicle.
	ModelVehicle() ModelVehicle
}

// SolutionVehicles is a slice of solution vehicles.
type SolutionVehicles []SolutionVehicle
