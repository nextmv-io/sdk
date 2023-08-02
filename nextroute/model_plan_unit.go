package nextroute

import "github.com/nextmv-io/sdk/nextroute/common"

// ModelPlanUnit is a set of stops. It is a set of stops
// that are required to be planned together on the same vehicle. For example,
// a unit can be a pickup and a delivery stop that are required to be planned
// together on the same vehicle.
type ModelPlanUnit interface {
	ModelData

	// Centroid returns the centroid of the unit. The centroid is the
	// average location of all stops in the unit.
	Centroid() (common.Location, error)

	// DirectedAcyclicGraph returns the [DirectedAcyclicGraph] of the plan
	// unit.
	DirectedAcyclicGraph() DirectedAcyclicGraph

	// Index returns the index of the invoking unit.
	Index() int

	// IsFixed returns true if the PlanUnit is fixed.
	IsFixed() bool

	// NumberOfStops returns the number of stops in the invoking unit.
	NumberOfStops() int

	// Stops returns the stops in the invoking unit.
	Stops() ModelStops
}

// ModelPlanUnits is a slice of plan units .
type ModelPlanUnits []ModelPlanUnit
