// Package schema provides the input and output schema for nextroute.
package schema

import (
	"time"
)

// Input is the default input schema for nextroute.
type Input struct {
	Options        any          `json:"options,omitempty"`
	Defaults       *Defaults    `json:"defaults,omitempty"`
	StopGroups     *[][]string  `json:"stop_groups,omitempty"`
	DurationMatrix *[][]float64 `json:"duration_matrix,omitempty"`
	DistanceMatrix *[][]float64 `json:"distance_matrix,omitempty"`
	Vehicles       []Vehicle    `json:"vehicles,omitempty"`
	Stops          []Stop       `json:"stops,omitempty"`
	CustomData     any          `json:"custom_data,omitempty"`
}

// Defaults contains default values for vehicles and stops.
type Defaults struct {
	Vehicles *VehicleDefaults `json:"vehicles,omitempty"`
	Stops    *StopDefaults    `json:"stops,omitempty"`
}

// VehicleDefaults contains default values for vehicles.
type VehicleDefaults struct {
	Capacity                any        `json:"capacity,omitempty"`
	StartLevel              any        `json:"start_level,omitempty"`
	StartLocation           *Location  `json:"start_location,omitempty"`
	EndLocation             *Location  `json:"end_location,omitempty"`
	Speed                   *float64   `json:"speed,omitempty"`
	StartTime               *time.Time `json:"start_time,omitempty"`
	EndTime                 *time.Time `json:"end_time,omitempty"`
	MaxStops                *int       `json:"max_stops,omitempty"`
	MaxDistance             *int       `json:"max_distance,omitempty"`
	MaxDuration             *int       `json:"max_duration,omitempty"`
	MaxWait                 *int       `json:"max_wait,omitempty"`
	CompatibilityAttributes *[]string  `json:"compatibility_attributes,omitempty"`
}

// StopDefaults contains default values for stops.
type StopDefaults struct {
	UnplannedPenalty        *int         `json:"unplanned_penalty,omitempty"`
	Quantity                any          `json:"quantity,omitempty"`
	StartTimeWindow         *[]time.Time `json:"start_time_window,omitempty"`
	MaxWait                 *int         `json:"max_wait,omitempty"`
	Duration                *int         `json:"duration,omitempty"`
	TargetArrivalTime       *time.Time   `json:"target_arrival_time,omitempty"`
	EarlyArrivalTimePenalty *float64     `json:"early_arrival_time_penalty,omitempty"`
	LateArrivalTimePenalty  *float64     `json:"late_arrival_time_penalty,omitempty"`
	CompatibilityAttributes *[]string    `json:"compatibility_attributes"`
}

// Vehicle represents a vehicle.
type Vehicle struct {
	Capacity                any        `json:"capacity,omitempty"`
	StartLevel              any        `json:"start_level,omitempty"`
	StartLocation           *Location  `json:"start_location,omitempty"`
	EndLocation             *Location  `json:"end_location,omitempty"`
	Speed                   *float64   `json:"speed,omitempty"`
	ID                      string     `json:"id,omitempty"`
	StartTime               *time.Time `json:"start_time,omitempty"`
	EndTime                 *time.Time `json:"end_time,omitempty"`
	CompatibilityAttributes *[]string  `json:"compatibility_attributes,omitempty"`
	MaxStops                *int       `json:"max_stops,omitempty"`
	MaxDistance             *int       `json:"max_distance,omitempty"`
	MaxDuration             *int       `json:"max_duration,omitempty"`
	MaxWait                 *int       `json:"max_wait,omitempty"`
	InitializationCost      *int       `json:"initialization_cost,omitempty"`
	CustomData              any        `json:"custom_data,omitempty"`
}

// Stop represents a stop.
type Stop struct {
	Precedes                any          `json:"precedes,omitempty"`
	Quantity                any          `json:"quantity,omitempty"`
	Succeeds                any          `json:"succeeds,omitempty"`
	TargetArrivalTime       *time.Time   `json:"target_arrival_time,omitempty"`
	StartTimeWindow         *[]time.Time `json:"start_time_window,omitempty"`
	MaxWait                 *int         `json:"max_wait,omitempty"`
	Duration                *int         `json:"duration,omitempty"`
	UnplannedPenalty        *int         `json:"unplanned_penalty,omitempty"`
	EarlyArrivalTimePenalty *float64     `json:"early_arrival_time_penalty,omitempty"`
	LateArrivalTimePenalty  *float64     `json:"late_arrival_time_penalty,omitempty"`
	CompatibilityAttributes *[]string    `json:"compatibility_attributes"`
	ID                      string       `json:"id,omitempty"`
	Location                Location     `json:"location,omitempty"`
	CustomData              any          `json:"custom_data,omitempty"`
}

// Location represents a geographical location.
type Location struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}
