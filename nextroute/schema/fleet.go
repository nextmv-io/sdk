// Package schema provides the input and output schema for nextroute.
package schema

import (
	"fmt"
	"time"
)

// FleetInput schema.
// DEPRECATION NOTICE: this part of the API is deprecated and is no longer
// maintained. It will be deleted soon. Please use [Input] instead.
type FleetInput struct {
	Options        *Options        `json:"options,omitempty"`
	Defaults       *FleetDefaults  `json:"defaults,omitempty"`
	Vehicles       []FleetVehicle  `json:"vehicles,omitempty"`
	Stops          []FleetStop     `json:"stops,omitempty"`
	StopGroups     [][]string      `json:"stop_groups,omitempty"`
	AlternateStops []FleetStop     `json:"alternate_stops,omitempty"`
	DurationGroups []DurationGroup `json:"duration_groups,omitempty"`
}

// FleetDefaults holds the fleet input default data.
// FleetInput schema.
// DEPRECATION NOTICE: this part of the API is deprecated and is no longer
// maintained. It will be deleted soon. Please use [Defaults] instead.
type FleetDefaults struct {
	Vehicles *FleetVehicleDefaults `json:"vehicles,omitempty"`
	Stops    *FleetStopDefaults    `json:"stops,omitempty"`
}

// FleetVehicleDefaults holds the fleet input vehicle default data.
// FleetInput schema.
// DEPRECATION NOTICE: this part of the API is deprecated and is no longer
// maintained. It will be deleted soon. Please use [VehicleDefaults] instead.
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
// DEPRECATION NOTICE: this part of the API is deprecated and is no longer
// maintained. It will be deleted soon. Please use [StopDefaults] instead.
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
// DEPRECATION NOTICE: this part of the API is deprecated and is no longer
// maintained. It will be deleted soon. Please use [Vehicle] instead.
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
// DEPRECATION NOTICE: this part of the API is deprecated and is no longer
// maintained. It will be deleted soon. Please use [Stop] instead.
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

// Options adds solver options to the input.
// DEPRECATION NOTICE: this part of the API is deprecated and is no longer
// maintained. It will be deleted soon. Please use [solve.Options] instead.
type Options struct {
	Solver *SolverOptions `json:"solver,omitempty"`
	Runner *RunnerOptions `json:"runner,omitempty"`
}

// RunnerOptions represent the solver runtime duration in legacy fleet.
// DEPRECATION NOTICE: this part of the API is deprecated and is no longer
// maintained. It will be deleted soon. Please use [solve.Options] instead.
type RunnerOptions struct {
	Output struct {
		Solutions string `json:"solutions,omitempty"`
	} `json:"output,omitempty"`
}

// SolverOptions represent the solver runtime duration in legacy fleet.
// DEPRECATION NOTICE: this part of the API is deprecated and is no longer
// maintained. It will be deleted soon. Please use [solve.Options] instead.
type SolverOptions struct {
	Limits *Limits `json:"limits,omitempty"`
}

// Limits represent the solver runtime limitation in fleet.
// DEPRECATION NOTICE: this part of the API is deprecated and is no longer
// maintained. It will be deleted soon. Please use [solve.Options] instead.
type Limits struct {
	Duration string `json:"duration,omitempty"`
}

// ToNextRoute converters a legacy cloud fleet input into nextroute input format.
func (fleetInput FleetInput) ToNextRoute() (Input, error) {
	input := Input{
		Defaults: &Defaults{},
	}
	stopCompats := make([]string, 0)
	vehicleCompats := make([]string, 0)
	// Use default values and add special handling for CompatibilityAttributes
	// and HardWindows.
	if fleetInput.Defaults != nil && fleetInput.Defaults.Stops != nil {
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
	for i := range fleetInput.Vehicles {
		fleetInput.Vehicles[i].CompatibilityAttributes =
			append(fleetInput.Vehicles[i].CompatibilityAttributes, vehicleCompats...)
	}

	// Create vehicles with special logic for backlog legacy needs.
	backlogStops := make(map[string]struct{})
	vehicles := make([]Vehicle, len(fleetInput.Vehicles))
	for i, v := range fleetInput.Vehicles {
		v := v
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

	for i, s := range fleetInput.Stops {
		if _, ok := backlogStops[s.ID]; !ok && s.CompatibilityAttributes == nil {
			fleetInput.Stops[i].CompatibilityAttributes = &stopCompats
		}
	}

	// Create stops with legacy backlog feature.
	stops := createStops(fleetInput, backlogStops)

	// Put new input format together and return it.
	input.StopGroups = &fleetInput.StopGroups
	input.DurationGroups = &fleetInput.DurationGroups
	input.Vehicles = vehicles
	input.Stops = stops

	return input, nil
}

func createStops(fleetInput FleetInput, backlogStops map[string]struct{}) []Stop {
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
			timeWindow := *s.HardWindow
			for i := range timeWindow {
				timeWindow[i] = roundToMinute(timeWindow[i])
			}
			stops[i].StartTimeWindow = timeWindow
		}
	}
	return stops
}

func roundToMinute(t time.Time) time.Time {
	if t.Minute() > 29 {
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute()+1, 0, 0, t.Location())
	}
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
}
