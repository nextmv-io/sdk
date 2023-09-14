package nextroute

import (
	"context"
	"time"

	"github.com/nextmv-io/sdk/connect"
)

// NextCheckModel is the check of a model returning a [NextCheckOutput].
func NextCheckModel(
	ctx context.Context,
	model Model,
	duration time.Duration,
	verbosity NextCheckVerbosity,
) (NextCheckOutput, error) {
	connect.Connect(con, &nextCheck)
	return nextCheckModel(ctx, model, duration, verbosity)
}

// NextCheckSolution is the check of a solution returning a [NextCheckOutput].
func NextCheckSolution(
	ctx context.Context,
	solution Solution,
	duration time.Duration,
	verbosity NextCheckVerbosity,
) (NextCheckOutput, error) {
	connect.Connect(con, &nextCheck)
	return nextCheckSolution(ctx, solution, duration, verbosity)
}

// NextCheckOptions are the options for a check.
type NextCheckOptions struct {
	// Duration is the maximum duration allowed for the sanity check to run.
	Duration  time.Duration `json:"duration" usage:"maximum duration of the sanity check" default:"30s"`
	Verbosity int           `json:"verbosity" usage:"verbosity of the sanity check, current available options are 0 (off),1 (low), 2 (medium), 3 (high)" default:"0"`
}

// NextCheckVerbosity is the verbosity of the check.
type NextCheckVerbosity int

const (
	// NextCheckOff does not run the check.
	NextCheckOff NextCheckVerbosity = iota
	// NextCheckLow checks if there is at least one move per plan unit.
	NextCheckLow
	// NextCheckMedium checks the number of moves per plan unit and the
	// number of vehicles that have moves.
	NextCheckMedium
	// NextCheckHigh count for each plan unit if it fits on a vehicle and
	// therefor how many units fit on a vehicle and if a plan unit has no
	// move what the constraints are that make the moves infeasible.
	NextCheckHigh
	// NextCheckVeryHigh checks the objective of the best move for each plan
	// unit and vehicle.
	NextCheckVeryHigh
)

// NextCheckOutput is the output of the  check.
type NextCheckOutput struct {
	// Error is the error raised during the check.
	Error *string `json:"error,omitempty"`
	// Remark is the remark of the check. It can be "ok", "timeout" or
	// anything else that should explain itself.
	Remark string `json:"remark"`
	// Verbosity is the verbosity of the check.
	Verbosity int `json:"verbosity"`
	// Duration is the input maximum duration in seconds of the check.
	DurationMaximum float64 `json:"duration_maximum"`
	// DurationUsed is the duration in seconds used for the check.
	DurationUsed float64 `json:"duration_used"`
	// Solution is the start solution of the check.
	Solution NextCheckStartSolution `json:"solution"`
	// Summary is the summary of the check.
	Summary NextCheckSummary `json:"summary"`
	// PlanUnits is the check of the individual plan units.
	PlanUnits []NextCheckPlanUnit `json:"plan_units"`
	// Vehicles is the check of the vehicles.
	Vehicles []NextCheckVehicle `json:"vehicles"`
}

// NextCheckStartSolution is the start solution of the check.
type NextCheckStartSolution struct {
	// StopsPlanned is the number of stops planned in the start solution.
	StopsPlanned int `json:"stops_planned"`
	// PlanUnitsPlanned is the number of units planned in the start
	// solution.
	PlanUnitsPlanned int `json:"plan_units_planned"`
	// PlanUnitsUnplanned is the number of units unplanned in the start
	// solution.
	PlanUnitsUnplanned int `json:"plan_units_unplanned"`
	// VehiclesUsed is the number of vehicles used in the start solution.
	VehiclesUsed int `json:"vehicles_used"`
	// VehiclesNotUsed is the number of vehicles not used in the start solution,
	// the empty vehicles.
	VehiclesNotUsed int `json:"vehicles_not_used"`
	// Objective is the objective of the start solution.
	Objective NextCheckObjective `json:"objective"`
}

// NextCheckSummary is the summary of the check.
type NextCheckSummary struct {
	// PlanUnitsToBeChecked is the number of plan units to be checked.
	PlanUnitsToBeChecked int `json:"plan_units_to_be_checked"`
	// PlanUnitsChecked is the number of plan units checked. If this is less
	// than [PlanUnitsToBeChecked] the check timed out.
	PlanUnitsChecked int `json:"plan_units_checked"`
	// PlanUnitsMoveFoundExecutable is the number of plan units for which at
	// least one move has been found and the move is executable.
	PlanUnitsBestMoveFound int `json:"plan_units_best_move_found"`
	// PlanUnitsHaveNoMove is the number of plan units for which no feasible
	// move has been found. This implies there is no move that can be executed
	// without violating a constraint.
	PlanUnitsHaveNoMove int `json:"plan_units_have_no_move"`
	// PlanUnitsMoveFoundTooExpensive is the number of plan units for which the
	// best move is executable but would increase the objective value instead
	// of decreasing it.
	PlanUnitsBestMoveIncreasesObjective int `json:"plan_units_best_move_increases_objective"`
	// PlanUnitsBestMoveFailed is the number of plan units for which the best
	// move can not be planned. This should not happen if all the constraints
	// are implemented correct.
	PlanUnitsBestMoveFailed int `json:"plan_units_best_move_failed"`
	// MovesFailed is the number of moves that failed. A move can fail if the
	// estimate of a constraint is incorrect. A constraint is incorrect if
	// [ModelConstraint.EstimateIsViolated] returns true and one of the
	// violation checks returns false. Violation checks are implementations of
	// one or more of the interfaces [SolutionStopViolationCheck],
	// [SolutionVehicleViolationCheck] or [SolutionViolationCheck] on the same
	// constraint. Most constraints do not need and do not have violation
	// checks as the estimate is perfect. The number of moves failed can be more
	// than one per plan unit as we continue to try moves on different vehicles
	// until we find a move that is executable or all vehicles have been
	// visited.
	MovesFailed int `json:"moves_failed"`
}

// NextCheckPlanUnit is the check of a plan unit.
type NextCheckPlanUnit struct {
	// ID is the ID of the plan unit. The ID of the plan unit is the slice of
	// ID's of the stops in the plan unit.
	Stops []string `json:"stops"`
	// HasBestMove is true if a move is found for the plan unit. A plan unit
	// has no move found if the plan unit is over-constrained or the move found
	// is too expensive.
	HasBestMove bool `json:"has_best_move"`
	// BestMoveIncreasesObjective is true if the best move for the plan unit
	// increases the objective.
	BestMoveIncreasesObjective bool `json:"best_move_increases_objective"`
	// BestMoveFailed is true if the plan unit best move failed to execute.
	BestMoveFailed bool `json:"best_move_failed"`
	// VehiclesHaveMoves is the number of vehicles that have moves for the plan
	// unit. Only calculated if the verbosity is medium or higher.
	VehiclesHaveMoves *int `json:"vehicles_have_moves,omitempty"`
	// VehicleWithMoves is the ID of the vehicles that have moves for the plan
	// unit. Only calculated if the verbosity is very high.
	VehiclesWithMoves *[]string `json:"vehicles_with_moves,omitempty"`
	// BestMoveObjective is the estimate of the objective of the best move if
	// the plan unit has a best move.
	BestMoveObjective *NextCheckObjective `json:"best_move_objective"`

	// Constraints is the constraints that are violated for the plan unit.
	Constraints *map[string]int `json:"constraints,omitempty"`
}

// NextCheckObjectiveTerm is the check of the individual terms of
// the objective for a move.
type NextCheckObjectiveTerm struct {
	// Name is the name of the objective term.
	Name string `json:"name"`
	// Factor is the factor of the objective term.
	Factor float64 `json:"factor"`
	// Base is the base value of the objective term.
	Base float64 `json:"base"`
	// Value is the value of the objective term which is the Factor times Base.
	Value float64 `json:"value"`
}

// NextCheckObjective is the estimate of an objective of a move.
type NextCheckObjective struct {
	// Vehicle is the ID of the vehicle for which it reports the objective.
	Vehicle *string `json:"vehicle,omitempty"`
	// Value is the value of the objective.
	Value float64 `json:"value"`
	// Terms is the check of the individual terms of the objective.
	Terms []NextCheckObjectiveTerm `json:"terms"`
}

// NextCheckVehicle is the check of a vehicle.
type NextCheckVehicle struct {
	// ID is the ID of the vehicle.
	ID string `json:"id"`
	// PlanUnitsHaveMoves is the number of plan units that have moves for the
	// vehicle. Only calculated if the depth is medium.
	PlanUnitsHaveMoves *int `json:"plan_units_have_moves,omitempty"`
}
