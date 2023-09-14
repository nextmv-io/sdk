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

	// HasPlanStopsUnit returns true if the stop belongs to a plan unit. For example,
	// start and end stops of a vehicle do not belong to a plan unit.
	HasPlanStopsUnit() bool

	// ID returns the identifier of the stop.
	ID() string

	// Index returns the unique index of the stop.
	Index() int

	// IsFirstOrLast returns true if the stop is the first or last stop of one
	// or more vehicles. A stop which is the first or last stop of one or more
	// vehicles is not allowed to be part of a plan unit. A stop which is the
	// first or last stop of one or more vehicles is by definition fixed.
	IsFirstOrLast() bool

	// IsFixed returns true if fixed.
	IsFixed() bool

	// Location returns the location of the stop.
	Location() common.Location

	// Model returns the model of the stop.
	Model() Model

	// EarliestStart returns the earliest start time of the stop.
	EarliestStart() (t time.Time)

	// Windows returns the time windows of the stop.
	Windows() [][2]time.Time

	// PlanStopsUnit returns the [ModelPlanStopsUnit] associated with the stop.
	// A stop is associated with at most one plan unit. Can be nil if the stop
	// is not part of a stops plan unit.
	PlanStopsUnit() ModelPlanStopsUnit

	// MeasureIndex returns the measure index of the invoking stop . This index
	// is not necessarily unique.
	// This index is used by the model expression constructed by the factory
	// NewMeasureByIndexExpression to calculate the value of the measure
	// expression. By default, the measure index is the same as the index of
	// the stop.
	MeasureIndex() int

	// SetEarliestStart sets the earliest start time of the stop.
	SetEarliestStart(t time.Time) error

	// SetMeasureIndex sets the reference index of the stop, by default the
	// measure index is the same as the index of the stop.
	// This index is used by the model expression constructed by the factory
	// NewMeasureByIndexExpression to calculate the value of the measure
	// expression.
	SetMeasureIndex(int)

	// SetWindows sets the time windows of the stop.
	SetWindows(windows [][2]time.Time) error

	// ToEarliestStartValue returns the earliest start time if the vehicle
	// arrives at the stop at the given arrival time in seconds since
	// [Model.Epoch].
	ToEarliestStartValue(arrival float64) float64

	// SetID sets the identifier of the stop. This identifier is not used by
	// nextroute, and therefore it does not have to be unique for nextroute
	// internally. However, to make this ID useful for debugging and reporting
	// it should be made unique.
	SetID(string)
}

// ModelStops is a slice of stops.
type ModelStops []ModelStop
