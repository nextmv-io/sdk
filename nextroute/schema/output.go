package schema

import (
	"time"

	"github.com/nextmv-io/sdk/nextroute"
)

// SolutionOutput represents a solutions as JSON.
type SolutionOutput struct {
	Unplanned []StopOutput               `json:"unplanned"`
	Vehicles  []VehicleOutput            `json:"vehicles"`
	Objective ObjectiveOutput            `json:"objective"`
	NextCheck *nextroute.NextCheckOutput `json:"next_check,omitempty"`
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
	RouteTravelDistance  int                 `json:"route_travel_distance,omitempty"`
	RouteStopsDuration   int                 `json:"route_stops_duration,omitempty"`
	RouteWaitingDuration int                 `json:"route_waiting_duration,omitempty"`
	RouteDuration        int                 `json:"route_duration"`
	CustomData           any                 `json:"custom_data,omitempty"`
	AlternateStops       *[]string           `json:"alternate_stops,omitempty"`
}

// PlannedStopOutput adds information to the input stop.
type PlannedStopOutput struct {
	Stop                     StopOutput `json:"stop"`
	TravelDuration           int        `json:"travel_duration"`
	CumulativeTravelDuration int        `json:"cumulative_travel_duration"`
	TravelDistance           int        `json:"travel_distance,omitempty"`
	CumulativeTravelDistance int        `json:"cumulative_travel_distance,omitempty"`
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

// CustomResultStatistics is an example of custom result statistics that can be
// added to the output and used in experiments.
type CustomResultStatistics struct {
	// ActivatedVehicles is the number of vehicles that were used in the
	// solution.
	ActivatedVehicles int `json:"activated_vehicles"`
	// UnplannedStops is the number of stops that were not planned in the
	// solution.
	UnplannedStops int `json:"unplanned_stops"`
	// MaxTravelDuration is the maximum travel duration of a vehicle in the
	// solution.
	MaxTravelDuration int `json:"max_travel_duration"`
	// MaxDuration is the maximum duration of a vehicle (including waiting
	// times) in the solution.
	MaxDuration int `json:"max_duration"`
	// MinTravelDuration is the minimum travel duration of a vehicle in the
	// solution, excluding vehicles that were not used.
	MinTravelDuration int `json:"min_travel_duration"`
	// MinDuration is the minimum duration of a vehicle (including waiting
	// times) in the solution, excluding vehicles that were not used.
	MinDuration int `json:"min_duration"`
	// MaxStopsInRoute is the maximum number of stops in a vehicle's route in
	// the solution. The start and end stops of the vehicle are not considered.
	MaxStopsInVehicle int `json:"max_stops_in_vehicle"`
	// MinStopsInRoute is the minimum number of stops in a vehicle's route in
	// the solution. The start and end stops of the vehicle are not considered.
	MinStopsInVehicle int `json:"min_stops_in_vehicle"`
}
