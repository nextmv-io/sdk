package nextroute

import (
	"time"

	"github.com/nextmv-io/sdk/nextroute/common"
)

// ModelStop is a stop to be assigned to a vehicle.
type ModelStop interface {
	ModelData

	// ClosestStops returns a slice containing the closest stops to the
	// invoking stop. The slice is sorted by increasing distance to the
	// location. The slice first stop is the stop itself. The distance used
	// is the common.Haversine distance between the stops. All the stops
	// in the model are used in the slice. Slice with similar distance are
	// sorted by their index (increasing).
	ClosestStops() ModelStops

	// HasPlanUnit returns true if the stop belongs to a plan unit. For example,
	// start and end stops of a vehicle do not belong to a plan unit.
	HasPlanUnit() bool

	// Index returns the index of the stop.
	Index() int

	// Location returns the location of the stop.
	Location() common.Location

	// Model returns the model of the stop.
	Model() Model

	// Name returns the name of the stop.
	Name() string

	// EarliestStart returns the earliest start time of the stop. Can be set
	// using the EarliestStart StopOption or using SetEarliestStart.
	EarliestStart() time.Time

	// PlanUnit returns the [ModelPlanUnit] associated with the stop. A stop
	// is associated with at most one plan unit. Can be nil if the stop is not
	// part of a plan unit.
	PlanUnit() ModelPlanUnit

	// SetEarliestStart sets the earliest start time of the stop.
	SetEarliestStart(time time.Time)
	// SetName sets the name of the stop.
	SetName(name string)
}

// ModelStops is a slice of stops.
type ModelStops []ModelStop
