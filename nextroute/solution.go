package nextroute

import (
	"time"

	"github.com/nextmv-io/sdk/alns"
	"github.com/nextmv-io/sdk/connect"
)

// NewSolution creates a new solution. The solution is created from the given
// model. The solution is empty and has no vehicles. Vehicles can be added to
// the solution using the NewVehicle() method.
func NewSolution(
	m Model,
) (Solution, error) {
	connect.Connect(con, &newSolution)
	return newSolution(m)
}

// Solution is a solution to a model.
type Solution interface {
	alns.Solution[Solution]

	// BestMove returns the best move for the given solution plan cluster. The
	// best move is the move that has the lowest score. If there are no moves
	// available for the given solution plan cluster, a move is returned which
	// is not executable, Move.IsExecutable.
	BestMove(SolutionPlanCluster) Move

	// ConstraintScore returns the constraint score of the solution. The
	// constraint score is the sum of the constraint scores. A constraint score
	// is zero if the constraint is not violated. Only soft constraints are
	// considered. Hard constraints are not considered as they should have a
	// zero score by definition.
	ConstraintScore() float64

	// EstimateDeltaScore estimates the delta score of the solution if the given
	// stop positions are moved. The delta score is the difference between the
	// score of the solution before the move and the score of the solution after
	// the move. The delta score is an estimate as the score of the solution
	// after the move is not calculated but estimated. The estimate is based on
	// ModelConstraint.EstimateDeltaScore and ModelObjective.EstimateDeltaScore.
	EstimateDeltaScore(
		stopPositions StopPositions,
	) (deltaScore float64,
		feasible bool,
		planPositionsHint StopPositionsHint,
	)

	// Marshal returns the JSON representation of the solution.
	Marshal() ([]byte, error)
	// Model returns the model of the solution.
	Model() Model

	// NewVehicle creates a new vehicle for the solution. The vehicle is added
	// to the solution. The vehicle is created from the given vehicle tye, start
	// time, start stop and end stop. The start time is the time the vehicle
	// starts at the start stop. The start stop is the stop the vehicle starts
	// at. The end stop is the stop the vehicle ends at. The start stop and end
	// stop must be part of the model. The start time must be after the start
	// time of the model.
	NewVehicle(
		vehicle ModelVehicleType,
		startTime time.Time,
		startStop ModelStop,
		endStop ModelStop,
	) (SolutionVehicle, error)

	// ObjectiveValue returns the objective value of the solution. The objective
	// value is the sum of the objective values defined in Model.Objective.
	ObjectiveValue() float64

	// PlannedPlanClusters returns the solution plan clusters that are planned.
	PlannedPlanClusters() SolutionPlanClusters

	// UnplannedPlanClusters returns the solution plan clusters that are not
	// planned.
	UnplannedPlanClusters() SolutionPlanClusters

	// Vehicles returns the vehicles of the solution.
	Vehicles() SolutionVehicles
}

// Solutions is a slice of solutions.
type Solutions []Solution

// Marshaller is a marshaller for a solution.
type Marshaller interface {
	// Marshal marshals the solution. The solution is marshaled to JSON.
	Marshal(s Solution) ([]byte, error)
}
