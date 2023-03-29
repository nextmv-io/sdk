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

// NewRandomSolution creates a new solution. The solution is created from the
// given model. The solution starts with an empty solution and will assign
//
//	a random plan cluster to a random vehicle. The remaining plan clusters
//
// are added to the solution in a random order at the best possible position.
func NewRandomSolution(
	m Model,
) (Solution, error) {
	connect.Connect(con, &newRandomSolution)
	return newRandomSolution(m)
}

// NewSweepSolution creates a new solution. The solution is created from the
// given model using a sweep construction heuristic.
func NewSweepSolution(
	m Model,
) (Solution, error) {
	connect.Connect(con, &newSweepSolution)
	return newSweepSolution(m)
}

// Solution is a solution to a model.
type Solution interface {
	alns.Solution[Solution]

	// BestMove returns the best move for the given solution plan cluster. The
	// best move is the move that has the lowest score. If there are no moves
	// available for the given solution plan cluster, a move is returned which
	// is not executable, Move.IsExecutable.
	BestMove(context.Context, SolutionPlanCluster) Move

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
	// SolutionVehicle returns the solution vehicle for the given model vehicle.
	SolutionVehicle(vehicle ModelVehicle) SolutionVehicle

	// UnplannedPlanClusters returns the solution plan clusters that are not
	// planned.
	UnplannedPlanClusters() SolutionPlanClusters

	// Vehicles returns the vehicles of the solution.
	Vehicles() SolutionVehicles
}

// Solutions is a slice of solutions.
type Solutions []Solution
