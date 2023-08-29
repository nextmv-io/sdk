package main

import (
	"strconv"
	"time"
)

// Output holds the output data of the solution.
type Output struct {
	Status         string             `json:"status"`
	AssignedShifts []OutputAssignment `json:"assigned_shifts"`
	Value          float64            `json:"value"`
}

// CustomResultStatistics add custom stats.
type CustomResultStatistics struct {
	NumberOfWorkers int `json:"number_of_workers"`
}

// Options holds custom configuration data.
type Options struct {
	OverSupplyPenalty  float64       `json:"over_supply_penalty" default:"1000" usage:"penalty for over-supplying a demand"`
	UnderSupplyPenalty float64       `json:"under_supply_penalty" default:"500" usage:"penalty for over-supplying a demand"`
	MaxHoursPerDay     time.Duration `json:"max_hours_per_day" default:"10h" usage:"maximum number of hours per day"`
	MaxHoursPerWeek    time.Duration `json:"max_hours_per_week" default:"40h" usage:"maximum number of hours per week"`
	MinHoursPerShift   time.Duration `json:"min_hours_per_shift" default:"2h" usage:"minimum number of hours per shift"`
	MaxHoursPerShift   time.Duration `json:"max_hours_per_shift" default:"8h" usage:"maximum number of hours per shift"`
	HoursBetweenShifts time.Duration `json:"hours_between_shifts" default:"8h" usage:"minimum number of hours between shifts"`
	SolverDuration     time.Duration `json:"solver_duration" default:"4m" usage:"maximum runtime for the solver"`
}

// Input represents a struct definition that can read input.json.
type Input struct {
	Workers         []Worker         `json:"workers"`
	RequiredWorkers []RequiredWorker `json:"required_workers"`
}

// Worker holds worker specific data.
type Worker struct {
	Availability []Availability `json:"availability"`
	ID           int            `json:"id"`
}

// Availability holds available times for a worker.
type Availability struct {
	Start CustomTime `json:"start"`
	End   CustomTime `json:"end"`
}

// RequiredWorker holds data about times and number of required workers per time window.
type RequiredWorker struct {
	RequiredWorkerID int        `json:"required_worker_id,omitempty"`
	Start            CustomTime `json:"start"`
	End              CustomTime `json:"end"`
	Count            int        `json:"count"`
}

// ID returned the RequiredWorker ID.
func (r RequiredWorker) ID() string {
	return strconv.Itoa(r.RequiredWorkerID)
}

// CustomTime represents a time.Time.
type CustomTime struct {
	time.Time
}

// UnmarshalJSON unmarshals a CustomTime.
func (t *CustomTime) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse(`"2006-01-02 15:04:05"`, string(b))
	if err != nil {
		return err
	}
	t.Time = date
	return
}

// MarshalJSON marshals a CustomTime.
func (t *CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(t.Format(`"2006-01-02 15:04:05"`)), nil
}

// OutputAssignment holds an assignment for a driver.
type OutputAssignment struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	WorkerID int       `json:"worker_id"`
}

// Assignment represents a shift assignment.
type Assignment struct {
	DemandsCovered []RequiredWorker `json:"demands_covered"`
	Start          time.Time        `json:"start"`
	End            time.Time        `json:"end"`
	Worker         Worker           `json:"worker"`
	Duration       time.Duration    `json:"duration"`
	AssignmentID   int              `json:"assignment_id"`
}

// DurationApart calculates the time to assignments are apart from each other.
func (a Assignment) DurationApart(other Assignment) time.Duration {
	if a.Start.After(other.End) {
		return a.Start.Sub(other.End)
	}
	if a.End.Before(other.Start) {
		return other.Start.Sub(a.End)
	}
	return 0
}

// ID returns the assignment id.
func (a Assignment) ID() string {
	return strconv.Itoa(a.AssignmentID)
}
