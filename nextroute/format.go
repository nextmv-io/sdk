package nextroute

import (
	"time"

	"github.com/nextmv-io/sdk/nextroute/schema"
	"github.com/nextmv-io/sdk/route"
)

// The Formatter interface is used to create custom JSON output.
type Formatter interface {
	ToOutput(Solution) any
}

// NewDefaultFormatter returns a default JSONFormatter to output solutions in
// our default way.
func NewDefaultFormatter() Formatter {
	return &defaultFormatter{}
}

type defaultFormatter struct{}

func (f *defaultFormatter) ToOutput(s Solution) any {
	// Process solutions of vehicles.
	solutionVehicles := s.Vehicles()
	vehicles := make([]VehicleOutput, len(solutionVehicles))
	for v, state := range solutionVehicles {
		vehicles[v] = output(state)
	}

	// Process unassigned stops.
	unassigned := make([]route.Stop, 0)
	for _, u := range s.UnplannedPlanClusters() {
		for _, v := range u.ModelPlanCluster().Stops() {
			unassigned = append(unassigned, route.Stop{
				ID: v.Name(),
				Position: route.Position{
					Lon: v.Location().Longitude(),
					Lat: v.Location().Latitude(),
				},
			})
		}
	}

	v := map[string]any{
		"unassigned": unassigned,
		"vehicles":   vehicles,
	}

	return v
}

// output constructs the output state of a vehicle.
func output(v SolutionVehicle) VehicleOutput {
	solutionRoute := v.SolutionStops()
	ID := v.ModelVehicle().Name()

	stops := make([]stopOutput, len(solutionRoute))
	data := v.ModelVehicle().VehicleType().Data().(schema.Vehicle)
	hasShiftStart := data.ShiftStart != nil
	hasStart := data.Start != nil
	hasEnd := data.End != nil

	// Prepare output route stops
	for i := 0; i < len(solutionRoute); i++ {
		// Determine matching input stop
		var stop route.Stop
		switch i {
		case 0:
			if hasStart {
				stop = makeStop(ID, true, data.Start)
			}
		case len(solutionRoute) - 1:
			if hasEnd {
				stop = makeStop(ID, false, data.End)
			}
		default:
			stop = route.Stop{
				ID: solutionRoute[i].ModelStop().Name(),
				Position: route.Position{
					Lon: solutionRoute[i].ModelStop().Location().Longitude(),
					Lat: solutionRoute[i].ModelStop().Location().Latitude(),
				},
			}
		}

		// Create output stop for this location
		stops[i] = stopOutput{Stop: stop}
		// Set ETA & ETD, if possible

		if hasShiftStart {
			eta := solutionRoute[i].Arrival()
			ets := solutionRoute[i].Start()
			etd := solutionRoute[i].End()

			stops[i].EstimatedArrival = &eta
			stops[i].EstimatedDeparture = &etd
			stops[i].EstimatedService = &ets
		}
	}

	// Slice output route according to whether starts/ends are present
	startIdx, endIdx := 0, len(solutionRoute)
	if !hasStart {
		startIdx = 1
	}
	if !hasEnd {
		endIdx = len(solutionRoute) - 1
	}

	return VehicleOutput{
		ID:            ID,
		Route:         stops[startIdx:endIdx],
		RouteDuration: int(v.Last().CumulativeTravelDurationValue()),
	}
}

// VehicleOutput holds the solution of the ModelVehicle Routing Problem.
type VehicleOutput struct {
	ID            string       `json:"id"`
	Route         []stopOutput `json:"route"`
	RouteDuration int          `json:"route_duration"`
	RouteDistance int          `json:"route_distance"`
}

// stopOutput adds information to the input stop.
type stopOutput struct {
	route.Stop
	EstimatedArrival   *time.Time `json:"estimated_arrival,omitempty"`
	EstimatedDeparture *time.Time `json:"estimated_departure,omitempty"`
	EstimatedService   *time.Time `json:"estimated_service,omitempty"`
}

// stop builds a Stop from a vehicle's start or end location.
func makeStop(id string, start bool, p *route.Position) route.Stop {
	if start {
		id += "-start"
	} else {
		id += "-end"
	}

	return route.Stop{ID: id, Position: *p}
}
