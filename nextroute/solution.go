package nextroute

import (
	"context"

	"github.com/nextmv-io/sdk/alns"
	"github.com/nextmv-io/sdk/connect"
)

// NewSolution creates a new solution. The solution is created from the given
// model. The solution starts with all plan units unplanned. Once a solution
// has been created the model can no longer be changed, it becomes immutable.
func NewSolution(
	m Model,
) (Solution, error) {
	connect.Connect(con, &newSolution)
	return newSolution(m)
}

// NewRandomSolution creates a new solution. The solution is created from the
// given model. The solution starts with an empty solution and will assign
// a random plan unit to a random vehicle. The remaining plan units
// are added to the solution in a random order at the best possible position.
func NewRandomSolution(
	ctx context.Context,
	m Model,
) (Solution, error) {
	connect.Connect(con, &newRandomSolution)
	return newRandomSolution(ctx, m)
}

// NewSweepSolution creates a new solution. The solution is created from the
// given model using a sweep construction heuristic.
func NewSweepSolution(
	ctx context.Context,
	m Model,
) (Solution, error) {
	connect.Connect(con, &newSweepSolution)
	return newSweepSolution(ctx, m)
}

// Solution is a solution to a model.
type Solution interface {
	alns.Solution[Solution]

	// BestMove returns the best move for the given solution plan unit. The
	// best move is the move that has the lowest score. If there are no moves
	// available for the given solution plan unit, a move is returned which
	// is not executable, Move.IsExecutable. The third argument of type
	// interface{} is used to pass in filters defined by the interfaces
	// [SolutionMoveFilter] and [SolutionVehicleFilter].
	BestMove(context.Context, SolutionPlanUnit, interface{}) Move

	// FixedPlanUnits returns the solution plan units that are fixed.
	// Fixed plan units are plan units that are not allowed to be planned or
	// unplanned. The union of fixed, planned and unplanned plan units
	// is the set of all plan units in the model.
	FixedPlanUnits() ImmutableSolutionPlanUnitCollection

	// Model returns the model of the solution.
	Model() Model

	// ObjectiveValue returns the objective value for the objective in the
	// solution. Also returns 0.0 if the objective is not part of the solution.
	ObjectiveValue(objective ModelObjective) float64

	// PlannedPlanUnits returns the solution plan units that are planned as
	// a collection of solution plan units.
	PlannedPlanUnits() ImmutableSolutionPlanUnitCollection

	// SolutionPlanUnit returns the [SolutionPlanUnit] for the given
	// model plan unit.
	SolutionPlanUnit(planUnit ModelPlanUnit) SolutionPlanUnit
	// SolutionStop returns the solution stop for the given model stop.
	SolutionStop(stop ModelStop) SolutionStop
	// SolutionVehicle returns the solution vehicle for the given model vehicle.
	SolutionVehicle(vehicle ModelVehicle) SolutionVehicle

	// UnPlannedPlanUnits returns the solution plan units that are not
	// planned.
	UnPlannedPlanUnits() ImmutableSolutionPlanUnitCollection

	// Vehicles returns the vehicles of the solution.
	Vehicles() SolutionVehicles
}

// Solutions is a slice of solutions.
type Solutions []Solution
