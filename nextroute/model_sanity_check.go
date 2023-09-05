package nextroute

import (
	"context"
	"time"

	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/run/schema"
)

// SanityCheckOptions are the options for the sanity check.
type SanityCheckOptions struct {
	Enable bool `json:"enable" usage:"enable the sanity check" default:"false"`
	// Duration is the maximum duration allowed for the sanity check to run.
	Duration time.Duration `json:"duration" usage:"maximum duration of the sanity check" default:"30s"`
	Depth    int           `json:"depth" usage:"depth of the sanity check, deeper is more checking, current available depths are 0 and 2" default:"0"`
}

// SanityCheckDepth is the depth of the sanity check.
type SanityCheckDepth int

const (
	// SanityCheckLow checks if there is at least one move per plan unit.
	SanityCheckLow SanityCheckDepth = iota
	// SanityCheckMedium checks same as SanityCheckLow.
	SanityCheckMedium
	// SanityCheckHigh count for each plan unit if it fits on a vehicle and
	// therefor how many units fit on a vehicle.
	SanityCheckHigh
)

// ModelSanityCheck is the interface for the sanity check.
type ModelSanityCheck interface {
	// Check runs the sanity check. It returns the sanity check output and an
	// error if the sanity check itself fails. Duration is the maximum time
	// allowed for the sanity check to run. If the sanity check takes longer
	// the remark will be "timeout". Depth is the depth of the sanity check.
	Check(duration time.Duration, depth SanityCheckDepth) (SanityCheckOutput, error)
}

// SanityCheckOutput is the output of the sanity check.
type SanityCheckOutput struct {
	// Remark is the remark of the sanity check. It can be "ok", "timeout" or
	// anything else that should explain itself.
	Remark string `json:"remark"`
	// Error is the error raised during the sanity check. It can be empty.
	Error string `json:"error"`
	// InitialSolution is the sanity check of the initial solution.
	InitialSolution SanityCheckInitialSolution `json:"initial_solution"`
	// PlanUnits is the sanity check of the plan units.
	PlanUnits SanityCheckPlanUnits `json:"plan_units"`
	// Vehicles is the sanity check of the vehicles.
	Vehicles SanityCheckVehicles `json:"vehicles"`
}

// SanityCheckInitialSolution is the sanity check of the initial solution. The
// initial solution is the solution created by the solver with the initial stops
// that are fixed. The non-fixed stops are not planned.
type SanityCheckInitialSolution struct {
	// Feasible is true if the initial solution is feasible.
	Feasible bool `json:"feasible"`
	// DurationToCreate is the duration to create the initial solution.
	DurationToCreate float64 `json:"duration_to_create"`
	// StopsPlanned is the number of stops planned, only the fixed initial
	// stops.
	StopsPlanned int `json:"stops_planned"`
	// StopsPlannedFixed is the number of fixed stops planned.
	StopsPlannedFixed int `json:"stops_planned_fixed"`
	// PlanUnitsUnplanned is the number of units unplanned, those are the
	// initial plan units that are not fixed.
	PlanUnitsUnplanned int `json:"plan_units_unplanned"`
	// PlanUnitsUnplannedFixed is the number of fixed stops unplanned.
	PlanUnitsUnplannedFixed int `json:"plan_units_unplanned_fixed"`
	// VehiclesUsed is the number of vehicles used in the initial solution.
	VehiclesUsed int `json:"vehicles_used"`
	// VehiclesNotUsed is the number of vehicles not used in the initial, the
	// empty vehicles.
	VehiclesNotUsed int `json:"vehicles_not_used"`
}

// SanityCheckPlanUnits is the sanity check of the plan units.
type SanityCheckPlanUnits struct {
	// PlanUnits is the number of plan units.
	PlanUnits int `json:"plan_units"`
	// PlanUnitsToBePlanned is the number of plan units to be planned.
	PlanUnitsToBePlanned int `json:"plan_units_to_be_planned"`
	// BestMoveDuration is the duration to find the best move on a vehicle.
	BestMoveDuration SanityCheckData `json:"best_move_duration"`
	// ExecuteMoveDuration is the duration to execute the best move on a
	// vehicle.
	ExecuteMoveDuration SanityCheckData `json:"execute_move_duration"`
	// MovesFailed is the number of moves that failed. A move can fail if the
	// estimate of a constraint is wrong. This should not happen if all the
	// constraints are implemented correct.This is different from
	// PlanUnitsMoveFoundFailed, this counts each vehicle for which a move
	// failed. PlanUnitsMoveFoundFailed counts the plan units for
	// which at least one move failed.
	MovesFailed int `json:"moves_failed"`
	// PlanUnitsNoMoveFound is the number of plan units that have no move found.
	// These moves can not be planned ever, they are over-constrained.
	PlanUnitsNoMoveFound int `json:"plan_units_no_move_found"`
	// PlanUnitsMoveFoundFailed is the number of plan units for which at least
	// one move has been found but the move can not be planned. This should not
	// happen if all the constraints are implemented correct.
	PlanUnitsMoveFoundFailed int `json:"plan_units_move_found_failed"`
	// SanityCheckPlanUnits is the sanity check of the plan units.
	SanityCheckPlanUnits []*SanityCheckPlanUnit `json:"sanity_check_plan_units"`
}

// SanityCheckPlanUnit is the sanity check of a plan unit.
type SanityCheckPlanUnit struct {
	// ID is the ID of the plan unit. The ID of the plan unit is the slice of
	// ID's of the stops in the plan unit.
	Stops []string `json:"stops"`
	// VehiclesHaveMoves is the number of vehicles that have moves for the plan
	// unit. Only calculated if the depth is medium.
	VehiclesHaveMoves int `json:"vehicles_have_moves"`
	// MoveFoundTooExpensive is true if the plan unit best move is too
	// expensive. A move is too expensive if the move can be executed but the
	// delta value of the objective is too high, larger than zero.
	MoveFoundTooExpensive bool `json:"move_found_too_expensive"`
	// MoveFoundObjectives is the sanity check of the objectives of the plan
	// unit best move. Only calculated if MoveFoundTooExpensive is true and
	// reports the objectives of the plan unit on the cheapest vehicles.
	MoveFoundObjectives *SanityCheckObjective `json:"move_found_objective,omitempty"`
	// MoveFoundFailed is true if the plan unit best move failed to execute.
	MoveFoundFailed bool `json:"move_found_failed"`
	// NoMoveFound is true if no move is found for the plan unit. A plan unit
	// has no move found if the plan unit is over-constrained.
	NoMoveFound bool `json:"no_move_found"`
	// Evaluated is true if the plan unit is evaluated. If the sanity check
	// timed out the plan unit might not be evaluated due to time limit.
	Evaluated bool `json:"evaluated"`
}

// SanityCheckObjectiveTerm is the sanity check of the individual terms of
// the objective.
type SanityCheckObjectiveTerm struct {
	// Name is the name of the objective term.
	Name string `json:"objective"`
	// Factor is the factor of the objective term.
	Factor float64 `json:"term"`
	// Base is the base value of the objective term.
	Base float64 `json:"base"`
	// Value is the value of the objective term which is the Factor times Base.
	Value float64 `json:"value"`
}

// SanityCheckObjective is the sanity check of the objective.
type SanityCheckObjective struct {
	// Vehicle is the ID of the vehicle for which it reports the objective.
	Vehicle string `json:"vehicle"`
	// Value is the value of the objective.
	Value float64 `json:"value"`
	// Terms is the sanity check of the individual terms of the objective.
	Terms []SanityCheckObjectiveTerm
}

// SanityCheckVehicles is the sanity check of the vehicles.
type SanityCheckVehicles struct {
	// Vehicles is the number of vehicles.
	Vehicles int `json:"vehicles"`
	// SanityCheckVehicles is the sanity check of the vehicles.
	SanityCheckVehicles []*SanityCheckVehicle `json:"sanity_check_vehicles"`
}

// SanityCheckVehicle is the sanity check of a vehicle.
type SanityCheckVehicle struct {
	// ID is the ID of the vehicle.
	ID string `json:"id"`
	// PlanUnitsHaveMoves is the number of plan units that have moves for the
	// vehicle. Only calculated if the depth is medium.
	PlanUnitsHaveMoves int `json:"plan_units_have_moves"`
}

// SanityCheckData is the data of a sanity check property.
type SanityCheckData struct {
	Minimum float64 `json:"minimum"`
	Maximum float64 `json:"maximum"`
	Average float64 `json:"average"`
}

// NewModelSanityCheck returns a new ModelSanityCheck.
//
// The sanity check is a function that takes a model and returns a
// ModelSanityCheck. The sanity check is run on the model. The sanity check
// should return the sanity check output and an error if the sanity check
// itself fails.
//
//	 func someSolver(model nextroute.Model) (runSchema.Output, error) {
//	    output, err := nextroute.NewModelSanityCheck(model).Check(
//			duration,
//			sdkNextRoute.SanityCheckMedium,
//		   )
//	    if err != nil {
//			return runSchema.Output{}, err
//		  }
//	   return nextroute.Format(
//			ctx,
//			options.Format,
//			nil,
//			func(_ sdkNextRoute.Solution) any {
//				return output
//			},
//			nil,
//		), nil
func NewModelSanityCheck(model Model) ModelSanityCheck {
	connect.Connect(con, &newModelSanityCheck)
	return newModelSanityCheck(model)
}

// SanityCheckReport is the sanity check of a model.
func SanityCheckReport(
	ctx context.Context,
	model Model,
	duration time.Duration,
	sanityCheckDepth SanityCheckDepth) (schema.Output, error) {
	connect.Connect(con, &sanityCheckReport)
	return sanityCheckReport(ctx, model, duration, sanityCheckDepth)
}
