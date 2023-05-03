package schema

import (
	"time"

	"github.com/nextmv-io/sdk/route"
)

// JSONBasicSolution represents a solutions as JSON.
type JSONBasicSolution struct {
	Unassigned []route.Stop    `json:"unassigned"`
	Vehicles   []VehicleOutput `json:"vehicles"`
	Objective  JSONObjective   `json:"objective"`
}

// VehicleOutput holds the solution of the ModelVehicle Routing Problem.
type VehicleOutput struct {
	ID            string       `json:"id"`
	Route         []StopOutput `json:"route"`
	RouteDuration int          `json:"route_duration"`
}

// StopOutput adds information to the input stop.
type StopOutput struct {
	route.Stop
	EstimatedArrival   *time.Time `json:"estimated_arrival,omitempty"`
	EstimatedDeparture *time.Time `json:"estimated_departure,omitempty"`
	EstimatedService   *time.Time `json:"estimated_service,omitempty"`
	TargetTime         *time.Time `json:"target_time,omitempty"`
	Earliness          *float64   `json:"earliness,omitempty"`
	Lateness           *float64   `json:"lateness,omitempty"`
}
