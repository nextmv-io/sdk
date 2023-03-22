package nextroute

import (
	"context"

	"github.com/nextmv-io/sdk/alns"
	"github.com/nextmv-io/sdk/connect"
)

// NewSolution creates a new solution. The solution is created from the given
// model. The solution starts with all plan clusters unplanned. Once a solution
// has been created the model can no longer be changed, it becomes immutable.
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
	BestMove(context.Context, SolutionPlanCluster) Move

	// EstimateDeltaScore estimates the delta score of the solution if the given
	// stop positions are moved. The delta score is the difference between the
	// score of the solution before the move and the score of the solution after
	// the move. The delta score is an estimate as the score of the solution
	// after the move is not calculated but estimated. The estimate is based on
	// ModelObjective.EstimateDeltaScore.
	EstimateDeltaScore(
		stopPositions StopPositions,
	) (deltaScore float64,
		feasible bool,
		planPositionsHint StopPositionsHint,
	)

	// Model returns the model of the solution.
	Model() Model

	// ObjectiveValue returns the objective value of the solution. The objective
	// value is the sum of the objective values defined in Model.Objective.
	ObjectiveValue() float64

	// PlannedPlanClusters returns the solution plan clusters that are planned.
	PlannedPlanClusters() SolutionPlanClusters

	// SolutionPlanCluster returns the solution plan cluster for the given
	// model plan cluster.
	SolutionPlanCluster(planCluster ModelPlanCluster) SolutionPlanCluster
	// SolutionStop returns the solution stop for the given model stop.
	SolutionStop(stop ModelStop) SolutionStop

	// UnplannedPlanClusters returns the solution plan clusters that are not
	// planned.
	UnplannedPlanClusters() SolutionPlanClusters

	// Vehicles returns the vehicles of the solution.
	Vehicles() SolutionVehicles
}

// Solutions is a slice of solutions.
type Solutions []Solution
