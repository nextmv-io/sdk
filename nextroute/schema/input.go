// Package schema provides the input and output schema for nextroute.
package schema

import (
	"time"
)

// Input is the default input schema for nextroute.
type Input struct {
	Options        any              `json:"options,omitempty"`
	CustomData     any              `json:"custom_data,omitempty"`
	Defaults       *Defaults        `json:"defaults,omitempty"`
	StopGroups     *[][]string      `json:"stop_groups,omitempty"`
	DurationMatrix *[][]float64     `json:"duration_matrix,omitempty"`
	DistanceMatrix *[][]float64     `json:"distance_matrix,omitempty"`
	DurationGroups *[]DurationGroup `json:"duration_groups,omitempty"`
	Vehicles       []Vehicle        `json:"vehicles,omitempty"`
	Stops          []Stop           `json:"stops,omitempty"`
	AlternateStops *[]AlternateStop `json:"alternate_stops,omitempty"`
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
	Speed                   *float64   `json:"speed,omitempty" minimumExclusive:"0"`
	StartTime               *time.Time `json:"start_time,omitempty"`
	EndTime                 *time.Time `json:"end_time,omitempty"`
	MinStops                *int       `json:"min_stops,omitempty" minimum:"0"`
	MinStopsPenalty         *float64   `json:"min_stops_penalty,omitempty" minimum:"0"`
	MaxStops                *int       `json:"max_stops,omitempty" minimumExclusive:"0"`
	MaxDistance             *int       `json:"max_distance,omitempty" minimumExclusive:"0"`
	MaxDuration             *int       `json:"max_duration,omitempty" minimumExclusive:"0"`
	MaxWait                 *int       `json:"max_wait,omitempty" minimum:"0"`
	CompatibilityAttributes *[]string  `json:"compatibility_attributes,omitempty" uniqueItems:"true"`
	ActivationPenalty       *int       `json:"activation_penalty,omitempty" minimum:"0"`
	AlternateStops          *[]string  `json:"alternate_stops,omitempty" uniqueItems:"true"`
}

// StopDefaults contains default values for stops.
type StopDefaults struct {
	UnplannedPenalty        *int       `json:"unplanned_penalty,omitempty" minimum:"0"`
	Quantity                any        `json:"quantity,omitempty"`
	StartTimeWindow         any        `json:"start_time_window,omitempty"`
	MaxWait                 *int       `json:"max_wait,omitempty" minimum:"0"`
	Duration                *int       `json:"duration,omitempty" minimum:"0"`
	TargetArrivalTime       *time.Time `json:"target_arrival_time,omitempty"`
	EarlyArrivalTimePenalty *float64   `json:"early_arrival_time_penalty,omitempty" minimum:"0"`
	LateArrivalTimePenalty  *float64   `json:"late_arrival_time_penalty,omitempty" minimum:"0"`
	CompatibilityAttributes *[]string  `json:"compatibility_attributes,omitempty" uniqueItems:"true"`
}

// Vehicle represents a vehicle.
type Vehicle struct {
	Capacity                any            `json:"capacity,omitempty"`
	StartLevel              any            `json:"start_level,omitempty"`
	CustomData              any            `json:"custom_data,omitempty"`
	CompatibilityAttributes *[]string      `json:"compatibility_attributes,omitempty" uniqueItems:"true"`
	MaxDistance             *int           `json:"max_distance,omitempty" minimumExclusive:"0"`
	StopDurationMultiplier  *float64       `json:"stop_duration_multiplier,omitempty"`
	StartTime               *time.Time     `json:"start_time,omitempty"`
	EndTime                 *time.Time     `json:"end_time,omitempty"`
	EndLocation             *Location      `json:"end_location,omitempty"`
	MinStops                *int           `json:"min_stops,omitempty" minimum:"0"`
	MinStopsPenalty         *float64       `json:"min_stops_penalty,omitempty" minimum:"0"`
	MaxStops                *int           `json:"max_stops,omitempty" minimumExclusive:"0"`
	Speed                   *float64       `json:"speed,omitempty" minimumExclusive:"0"`
	MaxDuration             *int           `json:"max_duration,omitempty" minimumExclusive:"0"`
	MaxWait                 *int           `json:"max_wait,omitempty" minimum:"0"`
	ActivationPenalty       *int           `json:"activation_penalty,omitempty" minimum:"0"`
	StartLocation           *Location      `json:"start_location,omitempty"`
	AlternateStops          *[]string      `json:"alternate_stops,omitempty" uniqueItems:"true"`
	InitialStops            *[]InitialStop `json:"initial_stops,omitempty" uniqueItems:"true"`
	ID                      string         `json:"id,omitempty"`
}

// InitialStop represents an initial stop.
type InitialStop struct {
	Fixed *bool  `json:"fixed,omitempty"`
	ID    string `json:"id"`
}

// AlternateStop represents an alternate stop.
type AlternateStop struct {
	Quantity                any        `json:"quantity,omitempty"`
	Duration                *int       `json:"duration,omitempty" minimum:"0"`
	CustomData              any        `json:"custom_data,omitempty"`
	MaxWait                 *int       `json:"max_wait,omitempty" minimum:"0"`
	StartTimeWindow         any        `json:"start_time_window,omitempty"`
	UnplannedPenalty        *int       `json:"unplanned_penalty,omitempty" minimum:"0"`
	EarlyArrivalTimePenalty *float64   `json:"early_arrival_time_penalty,omitempty" minimum:"0"`
	LateArrivalTimePenalty  *float64   `json:"late_arrival_time_penalty,omitempty" minimum:"0"`
	TargetArrivalTime       *time.Time `json:"target_arrival_time,omitempty"`
	ID                      string     `json:"id,omitempty"`
	Location                Location   `json:"location,omitempty"`
}

// Stop represents a stop.
type Stop struct {
	Precedes                any        `json:"precedes,omitempty"`
	Quantity                any        `json:"quantity,omitempty"`
	Succeeds                any        `json:"succeeds,omitempty"`
	CustomData              any        `json:"custom_data,omitempty"`
	Duration                *int       `json:"duration,omitempty" minimum:"0"`
	MaxWait                 *int       `json:"max_wait,omitempty" minimum:"0"`
	StartTimeWindow         any        `json:"start_time_window,omitempty"`
	UnplannedPenalty        *int       `json:"unplanned_penalty,omitempty" minimum:"0"`
	EarlyArrivalTimePenalty *float64   `json:"early_arrival_time_penalty,omitempty" minimum:"0"`
	LateArrivalTimePenalty  *float64   `json:"late_arrival_time_penalty,omitempty" minimum:"0"`
	CompatibilityAttributes *[]string  `json:"compatibility_attributes,omitempty" uniqueItems:"true"`
	TargetArrivalTime       *time.Time `json:"target_arrival_time,omitempty"`
	ID                      string     `json:"id,omitempty"`
	Location                Location   `json:"location,omitempty"`
}

// Location represents a geographical location.
type Location struct {
	Lon float64 `json:"lon" minimum:"-180" maximum:"180`
	Lat float64 `json:"lat" minimum:"-90" maximum:"90"`
}

// DurationGroup represents a group of stops that get additional duration
// whenever a stop of the group is approached for the first time.
type DurationGroup struct {
	Group    []string `json:"group,omitempty" uniqueItems:"true"`
	Duration int      `json:"duration,omitempty" minimum:"0"`
}
