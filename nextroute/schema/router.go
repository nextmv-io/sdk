package schema

import (
	"fmt"
	"reflect"

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
// - Defaults (optional: collapse same values back to defaults)
// - InitializationCosts
// - Precedences
// - Windows
// - Shifts
// - Backlogs
// - VehicleAttributes
// - StopAttributes

// RouterToNextRoute transforms router input to nextroute input.
func RouterToNextRoute(routerInput any) Input {
	internalInput := internalizeRouterInput(routerInput)

	// Convert vehicles
	vehicles := make([]Vehicle, len(internalInput.Vehicles))
	for i, v := range internalInput.Vehicles {
		vehicles[i] = Vehicle{
			ID:       v,
			Capacity: internalInput.Capacities[i],
			Start: &route.Position{
				Lon: internalInput.Starts[i].Lon,
				Lat: internalInput.Starts[i].Lat,
			},
			End: &route.Position{
				Lon: internalInput.Ends[i].Lon,
				Lat: internalInput.Ends[i].Lat,
			},
			Speed: &internalInput.Velocities[i],
		}
	}

	// Convert stops
	stops := make([]Stop, len(internalInput.Stops))
	for i, s := range internalInput.Stops {
		stops[i] = Stop{
			ID: s.ID,
			Position: route.Position{
				Lon: s.Position.Lon,
				Lat: s.Position.Lat,
			},
			Quantity:          &internalInput.Quantities[i],
			UnassignedPenalty: &internalInput.Penalties[i],
		}
	}

	return Input{
		Vehicles: vehicles,
		Stops:    stops,
	}
}

func internalizeRouterInput(routerInput any) RouterInput {
	var internalInput RouterInput
	// Get Stops element via reflection
	stops, err := getField[route.Stop](routerInput, "Stops")
	if err != nil {
		panic(err)
	}
	internalInput.Stops = stops
	// TODO: internalize via reflection
	return internalInput
}

// getField returns the field with the given name as the provided type.
func getField[T any](v any, field string) ([]T, error) {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	if !f.IsValid() {
		return nil, fmt.Errorf("no such field: %s in %v", field, v)
	}
	return f.Interface().([]T), nil
}
