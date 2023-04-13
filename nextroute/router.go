package nextroute

import (
	"reflect"
	"time"

	"github.com/nextmv-io/sdk/route"
)

// RouterInput is the schema for the input of router.
type RouterInput struct {
	Stops               []route.Stop         `json:"stops"`
	Vehicles            []string             `json:"vehicles"`
	InitializationCosts []float64            `json:"initialization_costs"`
	Starts              []route.Position     `json:"starts"`
	Ends                []route.Position     `json:"ends"`
	Quantities          []int                `json:"quantities"`
	Capacities          []int                `json:"capacities"`
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
}

// TODO: Conversion is currently incomplete. We need to handle the following:
// - Defaults (we should probably ignore them)
// - InitializationCosts
// - Precedences
// - Windows
// - Shifts
// - Backlogs
// - VehicleAttributes
// - StopAttributes
// - Groups
// - ServiceTimes
// - AlternateStops
// - Limits
// - DurationLimits
// - DistanceLimits
// - ServiceGroups

// FromRouter transforms router input to nextroute input.
func FromRouter(routerInput RouterInput) Input {
	// Convert stop defaults
	stopDefaults := StopDefaults{}
	if anyAndAllEqual(routerInput.Penalties) {
		stopDefaults.UnassignedPenalty = &routerInput.Penalties[0]
	}
	if anyAndAllEqual(routerInput.Quantities) {
		stopDefaults.Quantity = &routerInput.Quantities[0]
	}
	if anyAndAllEqual(routerInput.Windows) {
		stopDefaults.HardWindow = &[]time.Time{
			routerInput.Windows[0].TimeWindow.Start,
			routerInput.Windows[0].TimeWindow.End,
		}
		stopDefaults.MaxWait = &routerInput.Windows[0].MaxWait
	}
	if anyAndAllEqual(routerInput.ServiceTimes) {
		stopDefaults.StopDuration = &routerInput.ServiceTimes[0].Duration
	}
	if anyAndAllEqual(routerInput.StopAttributes) {
		stopDefaults.CompatibilityAttributes = &routerInput.StopAttributes[0].Attributes
	}

	// Convert vehicle defaults
	vehicleDefaults := VehicleDefaults{}
	if anyAndAllEqual(routerInput.Capacities) {
		vehicleDefaults.Capacity = &routerInput.Capacities[0]
	}
	if anyAndAllEqual(routerInput.Starts) {
		vehicleDefaults.Start = &routerInput.Starts[0]
	}
	if anyAndAllEqual(routerInput.Ends) {
		vehicleDefaults.End = &routerInput.Ends[0]
	}
	if anyAndAllEqual(routerInput.Velocities) {
		vehicleDefaults.Speed = &routerInput.Velocities[0]
	}
	if anyAndAllEqual(routerInput.Shifts) {
		vehicleDefaults.ShiftStart = &routerInput.Shifts[0].Start
		vehicleDefaults.ShiftEnd = &routerInput.Shifts[0].End
	}
	if anyAndAllEqual(routerInput.VehicleAttributes) {
		vehicleDefaults.CompatibilityAttributes = &routerInput.VehicleAttributes[0].Attributes
	}

	defaults := Defaults{
		Stops:    &stopDefaults,
		Vehicles: &vehicleDefaults,
	}

	// Convert vehicles
	vehicles := make([]Vehicle, len(routerInput.Vehicles))
	for i, v := range routerInput.Vehicles {
		vehicles[i] = Vehicle{
			ID:       v,
			Capacity: routerInput.Capacities[i],
			Start: &route.Position{
				Lon: routerInput.Starts[i].Lon,
				Lat: routerInput.Starts[i].Lat,
			},
			End: &route.Position{
				Lon: routerInput.Ends[i].Lon,
				Lat: routerInput.Ends[i].Lat,
			},
			Speed: &routerInput.Velocities[i],
		}
	}

	// Convert stops
	stops := make([]Stop, len(routerInput.Stops))
	for i, s := range routerInput.Stops {
		stops[i] = Stop{
			ID: s.ID,
			Position: route.Position{
				Lon: s.Position.Lon,
				Lat: s.Position.Lat,
			},
			Quantity:          &routerInput.Quantities[i],
			UnassignedPenalty: &routerInput.Penalties[i],
		}
	}

	return Input{
		Vehicles: vehicles,
		Stops:    stops,
		Defaults: &defaults,
	}
}

func anyAndAllEqual[T any](v []T) bool {
	if len(v) == 0 {
		return false
	}
	for i := 1; i < len(v); i++ {
		if !reflect.DeepEqual(v[i], v[0]) {
			return false
		}
	}
	return true
}