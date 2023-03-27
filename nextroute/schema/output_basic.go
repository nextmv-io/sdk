package schema

import (
	"time"

	"github.com/nextmv-io/sdk/route"
)

// VehicleOutput holds the solution of the ModelVehicle Routing Problem.
type VehicleOutput struct {
	ID            string       `json:"id"`
	Route         []StopOutput `json:"route"`
	RouteDuration int          `json:"route_duration"`
	RouteDistance int          `json:"route_distance"`
}

// StopOutput adds information to the input stop.
type StopOutput struct {
	route.Stop
	EstimatedArrival   *time.Time `json:"estimated_arrival,omitempty"`
	EstimatedDeparture *time.Time `json:"estimated_departure,omitempty"`
	EstimatedService   *time.Time `json:"estimated_service,omitempty"`
}
