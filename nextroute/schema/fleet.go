// Package schema provides the input and output schema for nextroute.
package schema

import (
	"fmt"
	"time"
)

// FleetInput schema.
type FleetInput struct {
	Defaults       *FleetDefaults  `json:"defaults,omitempty"`
	Vehicles       []FleetVehicle  `json:"vehicles,omitempty"`
	Stops          []FleetStop     `json:"stops,omitempty"`
	StopGroups     [][]string      `json:"stop_groups,omitempty"`
	AlternateStops []FleetStop     `json:"alternate_stops,omitempty"`
	DurationGroups []DurationGroup `json:"duration_groups,omitempty"`
}

type FleetDefaults struct {
	Vehicles *FleetVehicleDefaults `json:"vehicles,omitempty"`
	Stops    *FleetStopDefaults    `json:"stops,omitempty"`
}

type FleetVehicleDefaults struct {
	Start                   *Location  `json:"start,omitempty"`
	End                     *Location  `json:"end,omitempty"`
	Speed                   *float64   `json:"speed,omitempty"`
	Capacity                any        `json:"capacity,omitempty"`
	ShiftStart              *time.Time `json:"shift_start,omitempty"`
	ShiftEnd                *time.Time `json:"shift_end,omitempty"`
	CompatibilityAttributes []string   `json:"compatibility_attributes,omitempty"`
	MaxStops                *int       `json:"max_stops,omitempty"`
	MaxDistance             *int       `json:"max_distance,omitempty"`
	MaxDuration             *int       `json:"max_duration,omitempty"`
}

type FleetStopDefaults struct {
	UnassignedPenalty       *int         `json:"unassigned_penalty,omitempty"`
	Quantity                any          `json:"quantity,omitempty"`
	HardWindow              *[]time.Time `json:"hard_window,omitempty"`
	MaxWait                 *int         `json:"max_wait,omitempty"`
	StopDuration            *int         `json:"stop_duration,omitempty"`
	TargetTime              *time.Time   `json:"target_time,omitempty"`
	EarlinessPenalty        *float64     `json:"earliness_penalty,omitempty"`
	LatenessPenalty         *float64     `json:"lateness_penalty,omitempty"`
	CompatibilityAttributes *[]string    `json:"compatibility_attributes,omitempty"`
}

type FleetVehicle struct {
	ID                      string     `json:"id,omitempty"`
	Start                   *Location  `json:"start,omitempty"`
	End                     *Location  `json:"end,omitempty"`
	Speed                   *float64   `json:"speed,omitempty"`
	Capacity                any        `json:"capacity,omitempty"`
	ShiftStart              *time.Time `json:"shift_start,omitempty"`
	ShiftEnd                *time.Time `json:"shift_end,omitempty"`
	CompatibilityAttributes []string   `json:"compatibility_attributes,omitempty"`
	MaxStops                *int       `json:"max_stops,omitempty"`
	MaxDistance             *int       `json:"max_distance,omitempty"`
	MaxDuration             *int       `json:"max_duration,omitempty"`
	StopDurationMultiplier  *float64   `json:"stop_duration_multiplier,omitempty"`
	Backlog                 []string   `json:"backlog,omitempty"`
	AlternateStops          []string   `json:"alternate_stops,omitempty"`
	InitializationCost      int        `json:"initialization_cost,omitempty"`
}

type FleetStop struct {
	ID                      string       `json:"id,omitempty"`
	Position                Location     `json:"position,omitempty"`
	UnassignedPenalty       *int         `json:"unassigned_penalty,omitempty"`
	Quantity                any          `json:"quantity,omitempty"`
	Precedes                any          `json:"precedes,omitempty"`
	Succeeds                any          `json:"succeeds,omitempty"`
	HardWindow              *[]time.Time `json:"hard_window,omitempty"`
	MaxWait                 *int         `json:"max_wait,omitempty"`
	StopDuration            *int         `json:"stop_duration,omitempty"`
	TargetTime              *time.Time   `json:"target_time,omitempty"`
	EarlinessPenalty        *float64     `json:"earliness_penalty,omitempty"`
	LatenessPenalty         *float64     `json:"lateness_penalty,omitempty"`
	CompatibilityAttributes *[]string    `json:"compatibility_attributes,omitempty"`
}

func FleetToNextRoute(fleetInput FleetInput) (Input, error) {
	input := Input{}
	stopCompats := make([]string, 0)
	vehicleCompats := make([]string, 0)
	if fleetInput.Defaults != nil {
		input.Defaults = &Defaults{}
		if fleetInput.Defaults.Stops != nil {
			if fleetInput.Defaults.Stops.CompatibilityAttributes != nil {
				stopCompats = *fleetInput.Defaults.Stops.CompatibilityAttributes
			}
			input.Defaults.Stops = &StopDefaults{
				UnplannedPenalty:        fleetInput.Defaults.Stops.UnassignedPenalty,
				Quantity:                fleetInput.Defaults.Stops.Quantity,
				StartTimeWindow:         fleetInput.Defaults.Stops.HardWindow,
				MaxWait:                 fleetInput.Defaults.Stops.MaxWait,
				Duration:                fleetInput.Defaults.Vehicles.MaxDuration,
				TargetArrivalTime:       fleetInput.Defaults.Stops.TargetTime,
				EarlyArrivalTimePenalty: fleetInput.Defaults.Stops.EarlinessPenalty,
				LateArrivalTimePenalty:  fleetInput.Defaults.Stops.LatenessPenalty,
			}
		}

		if fleetInput.Defaults.Vehicles != nil {
			vehicleCompats = fleetInput.Defaults.Vehicles.CompatibilityAttributes
			input.Defaults.Vehicles = &VehicleDefaults{
				Capacity:          fleetInput.Defaults.Vehicles.Capacity,
				StartLocation:     fleetInput.Defaults.Vehicles.Start,
				EndLocation:       fleetInput.Defaults.Vehicles.End,
				Speed:             fleetInput.Defaults.Vehicles.Speed,
				StartTime:         fleetInput.Defaults.Vehicles.ShiftStart,
				EndTime:           fleetInput.Defaults.Vehicles.ShiftEnd,
				MaxStops:          fleetInput.Defaults.Vehicles.MaxStops,
				MaxDistance:       fleetInput.Defaults.Vehicles.MaxDistance,
				MaxDuration:       fleetInput.Defaults.Vehicles.MaxDuration,
				StartLevel:        nil,
				MinStops:          nil,
				MinStopsPenalty:   nil,
				MaxWait:           nil,
				ActivationPenalty: nil,
			}
		}
	}

	for _, s := range fleetInput.Stops {
		if s.CompatibilityAttributes != nil {
			ca := append(*s.CompatibilityAttributes, stopCompats...)
			s.CompatibilityAttributes = &ca
		} else {
			s.CompatibilityAttributes = &stopCompats
		}
	}

	for _, v := range fleetInput.Vehicles {
		v.CompatibilityAttributes = append(v.CompatibilityAttributes, vehicleCompats...)
	}

	vehicles := make([]Vehicle, len(fleetInput.Vehicles))
	for i, v := range fleetInput.Vehicles {
		for _, ca := range v.CompatibilityAttributes {
			for _, b := range v.Backlog {
				v.CompatibilityAttributes = append(v.CompatibilityAttributes, fmt.Sprintf("%s_%s", ca, b))
			}
		}
		vehicles[i] = Vehicle{
			Capacity:                v.Capacity,
			CompatibilityAttributes: &v.CompatibilityAttributes,
			MaxDistance:             v.MaxDistance,
			StopDurationMultiplier:  v.StopDurationMultiplier,
			StartTime:               v.ShiftStart,
			EndTime:                 v.ShiftEnd,
			StartLocation:           v.Start,
			EndLocation:             v.End,
			MaxStops:                v.MaxStops,
			Speed:                   v.Speed,
			MaxDuration:             v.MaxDuration,
			ActivationPenalty:       &v.InitializationCost,
			ID:                      v.ID,
			InitialStops:            nil,
			StartLevel:              nil,
			CustomData:              nil,
			MinStops:                nil,
			MinStopsPenalty:         nil,
			MaxWait:                 nil,
		}
	}

	stops := make([]Stop, len(fleetInput.Stops))
	for i, s := range fleetInput.Stops {
		compats := make([]string, 0)
		if s.CompatibilityAttributes != nil {
			for _, ca := range *s.CompatibilityAttributes {
				compats = append(compats, fmt.Sprintf("%s_%s", ca, s.ID))
			}
		}
		stops[i] = Stop{
			Precedes:                s.Precedes,
			Quantity:                s.Quantity,
			Succeeds:                s.Succeeds,
			Duration:                s.StopDuration,
			MaxWait:                 s.MaxWait,
			UnplannedPenalty:        s.UnassignedPenalty,
			EarlyArrivalTimePenalty: s.EarlinessPenalty,
			LateArrivalTimePenalty:  s.LatenessPenalty,
			CompatibilityAttributes: &compats,
			TargetArrivalTime:       s.TargetTime,
			ID:                      s.ID,
			Location:                s.Position,
			CustomData:              nil,
		}
		if s.HardWindow != nil {
			stops[i].StartTimeWindow = s.HardWindow
		}
	}

	input.StopGroups = &fleetInput.StopGroups
	input.DurationGroups = &fleetInput.DurationGroups
	input.Vehicles = vehicles
	input.Stops = stops

	return input, nil
}
