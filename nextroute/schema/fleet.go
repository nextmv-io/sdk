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

// FleetDefaults holds the fleet input default data.
type FleetDefaults struct {
	Vehicles *FleetVehicleDefaults `json:"vehicles,omitempty"`
	Stops    *FleetStopDefaults    `json:"stops,omitempty"`
}

// FleetVehicleDefaults holds the fleet input vehicle default data.
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

// FleetStopDefaults holds the fleet input stop default data.
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

// FleetVehicle holds the fleet input vehicle data.
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

// FleetStop holds the fleet input stop data.
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

// FleetToNextRoute takes a legacy cloud fleet input and converts it into
// nextroute input format.
func FleetToNextRoute(fleetInput FleetInput) (Input, error) {
	input := Input{}
	stopCompats := make([]string, 0)
	vehicleCompats := make([]string, 0)
	// Use default values and add special handling for CompatibilityAttributes
	// and HardWindows.
	if fleetInput.Defaults != nil && fleetInput.Defaults.Stops != nil {
		input.Defaults = &Defaults{}
		input.Defaults.Stops = &StopDefaults{
			UnplannedPenalty:        fleetInput.Defaults.Stops.UnassignedPenalty,
			Quantity:                fleetInput.Defaults.Stops.Quantity,
			MaxWait:                 fleetInput.Defaults.Stops.MaxWait,
			Duration:                fleetInput.Defaults.Stops.StopDuration,
			TargetArrivalTime:       fleetInput.Defaults.Stops.TargetTime,
			EarlyArrivalTimePenalty: fleetInput.Defaults.Stops.EarlinessPenalty,
			LateArrivalTimePenalty:  fleetInput.Defaults.Stops.LatenessPenalty,
		}
		if fleetInput.Defaults.Stops.CompatibilityAttributes != nil {
			stopCompats = *fleetInput.Defaults.Stops.CompatibilityAttributes
		}
		if fleetInput.Defaults.Stops.HardWindow != nil {
			input.Defaults.Stops.StartTimeWindow = *fleetInput.Defaults.Stops.HardWindow
		}
	}

	// Use default values. nil values are not compatible between fleet and
	// nextroute.
	if fleetInput.Defaults != nil && fleetInput.Defaults.Vehicles != nil {
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

	// Handle special case of compatibility attributes, to use them to make
	// legacy backlogs more or less work like in legacy fleet.
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

	// Create vehicles with special logic for backlog legacy needs.
	backlogStops := make(map[string]struct{})
	vehicles := make([]Vehicle, len(fleetInput.Vehicles))
	for i, v := range fleetInput.Vehicles {
		newAttributes := make([]string, len(v.CompatibilityAttributes))
		copy(newAttributes, v.CompatibilityAttributes)
		for _, ca := range v.CompatibilityAttributes {
			for _, b := range v.Backlog {
				newAttributes = append(newAttributes, fmt.Sprintf("%s_%s", ca, b))
			}
		}
		newBacklog := make([]InitialStop, len(v.Backlog))
		falseBool := false
		for i, b := range v.Backlog {
			backlogStops[b] = struct{}{}
			newBacklog[i] = InitialStop{
				Fixed: &falseBool,
				ID:    b,
			}
		}
		vehicles[i] = Vehicle{
			Capacity:                v.Capacity,
			CompatibilityAttributes: &newAttributes,
			MaxDistance:             v.MaxDistance,
			StopDurationMultiplier:  v.StopDurationMultiplier,
			StartTime:               v.ShiftStart,
			EndTime:                 v.ShiftEnd,
			StartLocation:           v.Start,
			EndLocation:             v.End,
			MaxStops:                v.MaxStops,
			Speed:                   v.Speed,
			MaxDuration:             v.MaxDuration,
			InitialStops:            &newBacklog,
			ActivationPenalty:       &v.InitializationCost,
			ID:                      v.ID,
			StartLevel:              nil,
			CustomData:              nil,
			MinStops:                nil,
			MinStopsPenalty:         nil,
			MaxWait:                 nil,
		}
	}

	// Create stops with legacy backlog feature.
	stops := make([]Stop, len(fleetInput.Stops))
	for i, s := range fleetInput.Stops {
		compats := make([]string, 0)
		if s.CompatibilityAttributes != nil {
			if _, ok := backlogStops[s.ID]; ok {
				for _, ca := range *s.CompatibilityAttributes {
					compats = append(compats, fmt.Sprintf("%s_%s", ca, s.ID))
				}
			} else {
				compats = *s.CompatibilityAttributes
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
			stops[i].StartTimeWindow = *s.HardWindow
		}
	}

	// Put new input format together and return it.
	input.StopGroups = &fleetInput.StopGroups
	input.DurationGroups = &fleetInput.DurationGroups
	input.Vehicles = vehicles
	input.Stops = stops

	return input, nil
}
