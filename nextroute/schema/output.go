package schema

import (
	"time"
)

// SolutionOutput represents a solutions as JSON.
type SolutionOutput struct {
	Unplanned []StopOutput    `json:"unplanned"`
	Vehicles  []VehicleOutput `json:"vehicles"`
	Objective ObjectiveOutput `json:"objective"`
}

// StopOutput is the basic struct for a stop.
type StopOutput struct {
	ID         string   `json:"id"`
	Location   Location `json:"location"`
	CustomData any      `json:"custom_data,omitempty"`
}

// VehicleOutput holds the solution of the ModelVehicle Routing Problem.
type VehicleOutput struct {
	ID                   string              `json:"id"`
	Route                []PlannedStopOutput `json:"route"`
	RouteTravelDuration  int                 `json:"route_travel_duration"`
	RouteStopsDuration   int                 `json:"route_stops_duration,omitempty"`
	RouteWaitingDuration int                 `json:"route_waiting_duration,omitempty"`
	RouteDuration        int                 `json:"route_duration"`
	CustomData           any                 `json:"custom_data,omitempty"`
}

// PlannedStopOutput adds information to the input stop.
type PlannedStopOutput struct {
	Stop                     StopOutput `json:"stop"`
	TravelDuration           int        `json:"travel_duration"`
	CumulativeTravelDuration int        `json:"cumulative_travel_duration"`
	TargetArrivalTime        *time.Time `json:"target_arrival_time,omitempty"`
	ArrivalTime              *time.Time `json:"arrival_time,omitempty"`
	WaitingDuration          int        `json:"waiting_duration,omitempty"`
	StartTime                *time.Time `json:"start_time,omitempty"`
	Duration                 int        `json:"duration,omitempty"`
	EndTime                  *time.Time `json:"end_time,omitempty"`
	EarlyArrivalDuration     int        `json:"early_arrival_duration,omitempty"`
	LateArrivalDuration      int        `json:"late_arrival_duration,omitempty"`
	CustomData               any        `json:"custom_data,omitempty"`
}

// ObjectiveOutput represents an objective as JSON.
type ObjectiveOutput struct {
	Name       string            `json:"name"`
	Objectives []ObjectiveOutput `json:"objectives,omitempty"`
	Factor     float64           `json:"factor,omitempty"`
	Base       float64           `json:"base,omitempty"`
	Value      float64           `json:"value"`
	CustomData any               `json:"custom_data,omitempty"`
}
