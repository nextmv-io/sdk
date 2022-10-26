package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nextmv-io/sdk/route"
)

// input describes the expected json input by the runner. Features not needed
// can simply be deleted or commented out, but make sure that the corresponding
// option in `solver` is also commented out. In case you would like to support
// a different input format you can change the struct as you see fit. You may
// need to change some code in `solver` to use the new structure.
type input struct {
	Options        *options        `json:"options,omitempty"`
	Defaults       *defaults       `json:"defaults,omitempty"`
	Vehicles       []vehicle       `json:"vehicles,omitempty"`
	Stops          []stop          `json:"stops,omitempty"`
	StopGroups     [][]string      `json:"stop_groups,omitempty"`
	AlternateStops []stop          `json:"alternate_stops,omitempty"`
	DurationGroups []durationGroup `json:"duration_groups,omitempty"`
	// internal duration groups representation
	durationGroupsDomains route.DurationGroups
}

type options struct {
	Solver *solverOptions `json:"solver,omitempty"`
}

type solverOptions struct {
	Diagram *diagram `json:"diagram,omitempty"`
	Limits  *limits  `json:"limits,omitempty"`
}

type diagram struct {
	Width     *int       `json:"width,omitempty"`
	Expansion *expansion `json:"expansion,omitempty"`
}

type expansion struct {
	Limit *int `json:"limit,omitempty"`
}

type limits struct {
	Duration  *cloudDuration `json:"duration,omitempty"`
	Nodes     *int           `json:"nodes,omitempty"`
	Solutions *int           `json:"solutions,omitempty"`
}

type cloudDuration struct {
	time.Duration
}

func (d *cloudDuration) UnmarshalJSON(b []byte) (err error) {
	// Handle duration given as a string (use time.ParseDuration)
	if b[0] == '"' {
		sd := string(b[1 : len(b)-1])
		d.Duration, err = time.ParseDuration(sd)
		return
	}

	// If duration is given as a number, fall back to duration default
	// (nanoseconds).
	var id int64
	id, err = json.Number(string(b)).Int64()
	d.Duration = time.Duration(id)
	return
}

func (d cloudDuration) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, d.String())), nil
}

type defaults struct {
	Vehicles *vehicleDefaults `json:"vehicles,omitempty"`
	Stops    *stopDefaults    `json:"stops,omitempty"`
}

type vehicleDefaults struct {
	Start                   *position  `json:"start,omitempty"`
	End                     *position  `json:"end,omitempty"`
	Speed                   *float64   `json:"speed,omitempty"`
	Capacity                any        `json:"capacity,omitempty"`
	ShiftStart              *time.Time `json:"shift_start,omitempty"`
	ShiftEnd                *time.Time `json:"shift_end,omitempty"`
	CompatibilityAttributes []string   `json:"compatibility_attributes,omitempty"`
	MaxStops                *int       `json:"max_stops,omitempty"`
	MaxDistance             *int       `json:"max_distance,omitempty"`
	MaxDuration             *int       `json:"max_duration,omitempty"`
}

type stopDefaults struct {
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

type vehicle struct {
	ID       string    `json:"id,omitempty"`
	Start    *position `json:"start,omitempty"`
	End      *position `json:"end,omitempty"`
	Speed    *float64  `json:"speed,omitempty"`
	Capacity any       `json:"capacity,omitempty"`
	// Internal representation for dynamic external one
	capacity                map[string]int
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

type stop struct {
	ID                string         `json:"id,omitempty"`
	Position          route.Position `json:"position,omitempty"`
	UnassignedPenalty *int           `json:"unassigned_penalty,omitempty"`
	Quantity          any            `json:"quantity,omitempty"`
	// Internal representation for dynamic external one
	quantity map[string]int
	Precedes any `json:"precedes,omitempty"`
	// Internal representation for dynamic external one
	precedes []string
	Succeeds any `json:"succeeds,omitempty"`
	// Internal representation for dynamic external one
	succeeds                []string
	HardWindow              *[]time.Time `json:"hard_window,omitempty"`
	MaxWait                 *int         `json:"max_wait,omitempty"`
	StopDuration            *int         `json:"stop_duration,omitempty"`
	TargetTime              *time.Time   `json:"target_time,omitempty"`
	EarlinessPenalty        *float64     `json:"earliness_penalty,omitempty"`
	LatenessPenalty         *float64     `json:"lateness_penalty,omitempty"`
	CompatibilityAttributes *[]string    `json:"compatibility_attributes"`
}

type position struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

// durationGroup represents a group of stops to which a stop duration is
// assigned.
type durationGroup struct {
	Group    []string `json:"group,omitempty"`
	Duration int      `json:"duration,omitempty"`
}

// This is the routerInput representation for the router engine.
type routerInput struct {
	Stops               []route.Stop         `json:"stops"`
	Vehicles            []string             `json:"vehicles"`
	InitializationCosts []float64            `json:"initialization_costs"`
	Starts              []route.Position     `json:"starts"`
	Ends                []route.Position     `json:"ends"`
	Quantities          map[string][]int     `json:"quantities"`
	Capacities          map[string][]int     `json:"capacities"`
	Precedences         []route.Job          `json:"precedences"`
	Windows             []route.Window       `json:"windows"`
	Shifts              []route.TimeWindow   `json:"shifts"`
	Penalties           []int                `json:"penalties"`
	Backlogs            []route.Backlog      `json:"backlogs"`
	VehicleAttributes   []route.Attributes   `json:"vehicle_attributes"`
	StopAttributes      []route.Attributes   `json:"stop_attributes"`
	Velocities          []float64            `json:"velocities"`
	Groups              [][]string           `json:"groups"`
	ServiceTimes        []route.Service      `json:"service_times"`
	AlternateStops      []route.Alternate    `json:"alternate_stops"`
	Limits              []route.Limit        `json:"limits"`
	DurationLimits      []float64            `json:"duration_limits"`
	DistanceLimits      []float64            `json:"distance_limits"`
	ServiceGroups       []route.ServiceGroup `json:"service_groups"`
	Kinds               map[string]bool
	EarlinessPenalty    []int
	LatenessPenalty     []int
	TargetTimes         []int
	MaxStops            []int
}

// transform transforms the information from the JSON input and translates it to
// router inputs.
func (i *input) transform() routerInput {
	stops := make([]route.Stop, len(i.Stops))
	vehicles := make([]string, len(i.Vehicles))
	velocities := make([]float64, len(i.Vehicles))
	maxDistance := make([]float64, len(i.Vehicles))
	maxDuration := make([]float64, len(i.Vehicles))
	penalties := make([]int, len(i.Stops))
	serviceTimes := make([]route.Service, len(i.Stops))
	backlogs := make([]route.Backlog, 0)
	vehicleAttributes := make([]route.Attributes, 0)
	stopAttributes := make([]route.Attributes, 0)
	stopGroups := make([][]string, len(i.StopGroups))
	windows := make([]route.Window, len(i.Stops))
	starts := make([]route.Position, len(i.Vehicles))
	ends := make([]route.Position, len(i.Vehicles))
	shifts := make([]route.TimeWindow, 0)
	initializationCosts := make([]float64, len(i.Vehicles))
	quantities := make(map[string][]int)
	capacities := make(map[string][]int)
	earlinessPenalty := make([]int, len(i.Stops)+len(i.Vehicles)*2)
	latenessPenalty := make([]int, len(i.Stops)+len(i.Vehicles)*2)
	targetTimes := make([]int, len(i.Stops)+len(i.Vehicles)*2)
	serviceTimeGroups := make([]route.ServiceGroup, len(i.DurationGroups))
	precedences := make([]route.Job, 0)
	alternateStops := make([]route.Alternate, 0)
	maxStops := make([]int, len(i.Vehicles))

	copy(stopGroups, i.StopGroups)

	for s, stops := range i.DurationGroups {
		serviceTimeGroups[s] = route.ServiceGroup{
			Group:    stops.Group,
			Duration: stops.Duration,
		}
	}

	kinds := map[string]bool{}
	for _, stop := range i.Stops {
		for kind := range stop.quantity {
			kinds[kind] = true
		}
	}
	for _, vehicle := range i.Vehicles {
		for kind := range vehicle.capacity {
			kinds[kind] = true
		}
	}

	for kind := range kinds {
		quantities[kind] = make([]int, len(i.Stops))
		capacities[kind] = make([]int, len(i.Vehicles))
	}
	hasEarlinessPenalty := false
	hasLatenessPenalty := false

	for s, stop := range i.Stops {
		stops[s] = route.Stop{ID: stop.ID, Position: stop.Position}
		penalties[s] = *stop.UnassignedPenalty
		if stop.EarlinessPenalty != nil {
			hasEarlinessPenalty = true
			earlinessPenalty[s] = int(*stop.EarlinessPenalty)
		}
		if stop.LatenessPenalty != nil {
			hasLatenessPenalty = true
			latenessPenalty[s] = int(*stop.LatenessPenalty)
		}
		if stop.TargetTime != nil {
			targetTimes[s] = int(stop.TargetTime.Unix())
		} else {
			targetTimes[s] = -1
		}
		for kind, quant := range stop.quantity {
			quantities[kind][s] = quant
		}
		if len(*stop.HardWindow) > 0 {
			windows[s] = route.Window{
				MaxWait: *stop.MaxWait,
				TimeWindow: route.TimeWindow{
					Start: (*stop.HardWindow)[0],
					End:   (*stop.HardWindow)[1],
				},
			}
		}

		for _, p := range stop.precedes {
			precedences = append(
				precedences,
				route.Job{PickUp: stop.ID, DropOff: p},
			)
		}

		for _, p := range stop.succeeds {
			precedences = append(
				precedences,
				route.Job{PickUp: p, DropOff: stop.ID},
			)
		}
		serviceTimes[s] = route.Service{
			ID:       stop.ID,
			Duration: *stop.StopDuration,
		}
		if len(*stop.CompatibilityAttributes) > 0 {
			stopAttributes = append(stopAttributes, route.Attributes{
				ID:         stop.ID,
				Attributes: *stop.CompatibilityAttributes,
			})
		}
	}

	for v, vehicle := range i.Vehicles {
		maxStops[v] = *vehicle.MaxStops
		vehicles[v] = vehicle.ID
		targetTimes[len(i.Stops)+v*2] = -1
		targetTimes[len(i.Stops)+v*2+1] = -1
		for kind, capacity := range vehicle.capacity {
			capacities[kind][v] = capacity
		}
		if vehicle.Start != nil {
			starts[v] = route.Position{
				Lon: vehicle.Start.Lon,
				Lat: vehicle.Start.Lat,
			}
		}
		if vehicle.End != nil {
			ends[v] = route.Position{
				Lon: vehicle.End.Lon,
				Lat: vehicle.End.Lat,
			}
		}
		start := time.Time{}
		end := time.Time{}
		if vehicle.ShiftStart != nil {
			start = *vehicle.ShiftStart
		}
		if vehicle.ShiftEnd != nil {
			end = *vehicle.ShiftEnd
		}
		if !start.IsZero() || !end.IsZero() {
			shifts = append(shifts, route.TimeWindow{
				Start: start,
				End:   end,
			})
		}
		initializationCosts[v] = float64(vehicle.InitializationCost)
		velocities[v] = *vehicle.Speed
		if len(vehicle.Backlog) > 0 {
			backlogs = append(backlogs, route.Backlog{
				VehicleID: vehicle.ID,
				Stops:     vehicle.Backlog,
			})
		}
		if len(vehicle.CompatibilityAttributes) > 0 {
			vehicleAttributes = append(vehicleAttributes, route.Attributes{
				ID:         vehicle.ID,
				Attributes: vehicle.CompatibilityAttributes,
			})
		}
		if len(vehicle.AlternateStops) > 0 {
			alternateStops = append(alternateStops, route.Alternate{
				VehicleID: vehicle.ID,
				Stops:     vehicle.AlternateStops,
			})
		}
		maxDistance[v] = float64(*vehicle.MaxDistance)
		maxDuration[v] = float64(*vehicle.MaxDuration)
	}

	if !hasEarlinessPenalty {
		earlinessPenalty = []int{}
	}
	if !hasLatenessPenalty {
		latenessPenalty = []int{}
	}

	return routerInput{
		Stops:               stops,
		Vehicles:            vehicles,
		InitializationCosts: initializationCosts,
		Starts:              starts,
		Ends:                ends,
		Quantities:          quantities,
		Capacities:          capacities,
		Precedences:         precedences,
		Windows:             windows,
		Shifts:              shifts,
		Penalties:           penalties,
		Backlogs:            backlogs,
		VehicleAttributes:   vehicleAttributes,
		StopAttributes:      stopAttributes,
		Velocities:          velocities,
		Groups:              stopGroups,
		ServiceTimes:        serviceTimes,
		AlternateStops:      alternateStops,
		DurationLimits:      maxDuration,
		DistanceLimits:      maxDistance,
		Kinds:               kinds,
		ServiceGroups:       serviceTimeGroups,
		EarlinessPenalty:    earlinessPenalty,
		LatenessPenalty:     latenessPenalty,
		TargetTimes:         targetTimes,
		MaxStops:            maxStops,
	}
}
