package main

import (
	"strconv"
	"time"

	"github.com/nextmv-io/sdk/mip"
)

// output holds the output data of the solution.
type output struct {
	AssignedShifts        []outputAssignment `json:"assigned_shifts"`
	NumberAssignedWorkers int                `json:"number_assigned_workers"`
}

// options holds custom configuration data.
type options struct {
	Penalty            penalty       `json:"penalty" usage:"set penalties for over and under supply of workers"`
	MaxHoursPerDay     time.Duration `json:"max_hours_per_day" default:"10h" usage:"maximum number of hours per day"`
	MaxHoursPerWeek    time.Duration `json:"max_hours_per_week" default:"40h" usage:"maximum number of hours per week"`
	MinHoursPerShift   time.Duration `json:"min_hours_per_shift" default:"2h" usage:"minimum number of hours per shift"`
	MaxHoursPerShift   time.Duration `json:"max_hours_per_shift" default:"8h" usage:"maximum number of hours per shift"`
	HoursBetweenShifts time.Duration `json:"hours_between_shifts" default:"8h" usage:"minimum number of hours between shifts"`
	Limits             mip.Limits    `json:"limits" usage:"holds a field to set the maximum duration of the solver"`
}

type penalty struct {
	OverSupply  float64 `json:"over_supply_penalty" default:"1000" usage:"penalty for over-supplying a demand"`
	UnderSupply float64 `json:"under_supply_penalty" default:"500" usage:"penalty for over-supplying a demand"`
}

// input represents a struct definition that can read input.json.
type input struct {
	Workers         []worker         `json:"workers"`
	RequiredWorkers []requiredWorker `json:"required_workers"`
}

// worker holds worker specific data.
type worker struct {
	Availability []availability `json:"availability"`
	ID           int            `json:"id"`
}

// availability holds available times for a worker.
type availability struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// requiredWorker holds data about times and number of required workers per time window.
type requiredWorker struct {
	RequiredWorkerID int       `json:"required_worker_id,omitempty"`
	Start            time.Time `json:"start"`
	End              time.Time `json:"end"`
	Count            int       `json:"count"`
}

// ID returned the RequiredWorker ID.
func (r requiredWorker) ID() string {
	return strconv.Itoa(r.RequiredWorkerID)
}

// outputAssignment holds an assignment for a worker.
type outputAssignment struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	WorkerID int       `json:"worker_id"`
}

// assignment represents a shift assignment.
type assignment struct {
	DemandsCovered []requiredWorker `json:"demands_covered"`
	Start          time.Time        `json:"start"`
	End            time.Time        `json:"end"`
	Worker         worker           `json:"worker"`
	Duration       time.Duration    `json:"duration"`
	AssignmentID   int              `json:"assignment_id"`
}

// DurationApart calculates the time to assignments are apart from each other.
func (a assignment) DurationApart(other assignment) time.Duration {
	if a.Start.After(other.End) {
		return a.Start.Sub(other.End)
	}
	if a.End.Before(other.Start) {
		return other.Start.Sub(a.End)
	}
	return 0
}

// ID returns the assignment id.
func (a assignment) ID() string {
	return strconv.Itoa(a.AssignmentID)
}
