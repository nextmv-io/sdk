// Package schema provides the input and output schema for nextroute.
package schema

import (
	"time"

	"github.com/nextmv-io/sdk/route"
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
	Capacity                any             `json:"capacity,omitempty"`
	Start                   *route.Position `json:"start,omitempty"`
	End                     *route.Position `json:"end,omitempty"`
	Speed                   *float64        `json:"speed,omitempty"`
	ShiftStart              *time.Time      `json:"shift_start,omitempty"`
	ShiftEnd                *time.Time      `json:"shift_end,omitempty"`
	MaxStops                *int            `json:"max_stops,omitempty"`
	MaxDistance             *int            `json:"max_distance,omitempty"`
	MaxDuration             *int            `json:"max_duration,omitempty"`
	CompatibilityAttributes *[]string       `json:"compatibility_attributes,omitempty"`
}

// StopDefaults contains default values for stops.
type StopDefaults struct {
	UnassignedPenalty       *int         `json:"unassigned_penalty,omitempty"`
	Quantity                any          `json:"quantity,omitempty"`
	HardWindow              *[]time.Time `json:"hard_window,omitempty"`
	MaxWait                 *int         `json:"max_wait,omitempty"`
	StopDuration            *int         `json:"stop_duration,omitempty"`
	TargetTime              *time.Time   `json:"target_time,omitempty"`
	EarlinessPenalty        *float64     `json:"earliness_penalty,omitempty"`
	LatenessPenalty         *float64     `json:"lateness_penalty,omitempty"`
	CompatibilityAttributes *[]string    `json:"compatibility_attributes"`
}

// Vehicle represents a vehicle.
type Vehicle struct {
	Capacity                any             `json:"capacity,omitempty"`
	Start                   *route.Position `json:"start,omitempty"`
	End                     *route.Position `json:"end,omitempty"`
	Speed                   *float64        `json:"speed,omitempty"`
	ID                      string          `json:"id,omitempty"`
	ShiftStart              *time.Time      `json:"shift_start,omitempty"`
	ShiftEnd                *time.Time      `json:"shift_end,omitempty"`
	CompatibilityAttributes *[]string       `json:"compatibility_attributes,omitempty"`
	MaxStops                *int            `json:"max_stops,omitempty"`
	MaxDistance             *int            `json:"max_distance,omitempty"`
	MaxDuration             *int            `json:"max_duration,omitempty"`
	InitializationCost      *int            `json:"initialization_cost,omitempty"`
	CustomData              any             `json:"custom_data,omitempty"`
}

// Stop represents a stop.
type Stop struct {
	Precedes                any            `json:"precedes,omitempty"`
	Quantity                any            `json:"quantity,omitempty"`
	Succeeds                any            `json:"succeeds,omitempty"`
	TargetTime              *time.Time     `json:"target_time,omitempty"`
	HardWindow              *[]time.Time   `json:"hard_window,omitempty"`
	MaxWait                 *int           `json:"max_wait,omitempty"`
	StopDuration            *int           `json:"stop_duration,omitempty"`
	UnassignedPenalty       *int           `json:"unassigned_penalty,omitempty"`
	EarlinessPenalty        *float64       `json:"earliness_penalty,omitempty"`
	LatenessPenalty         *float64       `json:"lateness_penalty,omitempty"`
	CompatibilityAttributes *[]string      `json:"compatibility_attributes"`
	ID                      string         `json:"id,omitempty"`
	Position                route.Position `json:"position,omitempty"`
	CustomData              any            `json:"custom_data,omitempty"`
}
