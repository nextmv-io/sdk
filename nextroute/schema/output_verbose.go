package schema

import (
	"time"

	"github.com/nextmv-io/sdk/alns"
)

// JSONSolution represents a solutions as JSON.
type JSONSolution struct {
	Epoch       time.Time               `json:"epoch"`
	Unplanned   []JSONModelStop         `json:"unplanned"`
	Vehicles    []JSONVehicle           `json:"vehicles"`
	Progression []alns.ProgressionEntry `json:"progression"`
	Objective   JSONObjective           `json:"objective"`
}

// JSONObjective represents an objective as JSON.
type JSONObjective struct {
	Name       string          `json:"name"`
	Objectives []JSONObjective `json:"objectives,omitempty"`
	Factor     float64         `json:"factor,omitempty"`
	Base       float64         `json:"base,omitempty"`
	Value      float64         `json:"value"`
}

// JSONVehicle represents a vehicle as JSON.
type JSONVehicle struct {
	Start time.Time  `json:"start"`
	End   time.Time  `json:"end"`
	Name  string     `json:"name"`
	Stops []JSONStop `json:"stops"`
	Index int        `json:"index"`
}

// JSONModelStop represents a model stop as JSON.
type JSONModelStop struct {
	Name     string       `json:"name"`
	Location JSONLocation `json:"location"`
	Index    int          `json:"index"`
}

// JSONStop represents a stop as JSON.
type JSONStop struct {
	Slack             time.Duration         `json:"slack"`
	Arrival           time.Time             `json:"arrival"`
	EarliestStart     time.Time             `json:"earliest_start"`
	Start             time.Time             `json:"start"`
	End               time.Time             `json:"end"`
	ExpressionValues  []JSONExpressionValue `json:"expression_values"`
	ConstraintValues  []JSONConstraintValue `json:"constraint_values"`
	ConstraintReports []map[string]any      `json:"constraint_reports"`
	JSONModelStop
	Position       int           `json:"position"`
	TravelDuration time.Duration `json:"travel_duration"`
}

// JSONExpressionValue represents an expression value as JSON.
type JSONExpressionValue struct {
	Name       string  `json:"name"`
	Value      float64 `json:"value"`
	Cumulative float64 `json:"cumulative"`
}

// JSONConstraintValue represents a value of a constraint.
type JSONConstraintValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// JSONLocation represents a location as JSON.
type JSONLocation struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
