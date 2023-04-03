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
	Centroid() (common.Location, error)

	// Duration returns the duration of the vehicle. The duration is the
	// time the vehicle is on the road. The duration is the time between
	// the start time and the end time.
	Duration() time.Duration
	// DurationValue returns the duration value of the vehicle. The duration
	// value is the value of the duration of the vehicle. The duration value
	// is the value in model duration units.
	DurationValue() float64

	// End returns the end time of the vehicle. The end time is the time
	// the vehicle ends at the end stop.
	End() time.Time
	// EndValue returns the end value of the vehicle. The end value is the
	// value of the end of the last stop. The end value is the value in
	// model duration units since the model epoch.
	EndValue() float64

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

	// ModelVehicle returns the modeled vehicle type of the vehicle.
	ModelVehicle() ModelVehicle

	// NumberOfStops returns the number of stops in the vehicle. The start
	// and end stops are not considered.
	NumberOfStops() int

	// SolutionStops returns the stops in the vehicle. The start and end
	// stops are included in the returned stops.
	SolutionStops() SolutionStops
	// Start returns the start time of the vehicle. The start time is
	// the time the vehicle starts at the start stop, it has been set
	// in the factory method of the vehicle Solution.NewVehicle.
	Start() time.Time
	// StartValue returns the start value of the vehicle. The start value
	// is the value of the start of the first stop. The start value is
	// the value in model duration units since the model epoch.
	StartValue() float64
}

// SolutionVehicles is a slice of solution vehicles.
type SolutionVehicles []SolutionVehicle
