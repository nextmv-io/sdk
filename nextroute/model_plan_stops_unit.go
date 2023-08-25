package nextroute

import "github.com/nextmv-io/sdk/nextroute/common"

// ModelPlanStopsUnit is a set of stops. It is a set of stops
// that are required to be planned together on the same vehicle. For example,
// a unit can be a pickup and a delivery stop that are required to be planned
// together on the same vehicle.
type ModelPlanStopsUnit interface {
	ModelPlanUnit

	// Centroid returns the centroid of the unit. The centroid is the
	// average location of all stops in the unit.
	Centroid() (common.Location, error)

	// DirectedAcyclicGraph returns the [DirectedAcyclicGraph] of the plan
	// unit.
	DirectedAcyclicGraph() DirectedAcyclicGraph

	// NumberOfStops returns the number of stops in the invoking unit.
	NumberOfStops() int

	// Stops returns the stops in the invoking unit.
	Stops() ModelStops
}

// ModelPlanStopsUnits is a slice of model plan stops units .
type ModelPlanStopsUnits []ModelPlanStopsUnit
