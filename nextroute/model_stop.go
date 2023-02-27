package nextroute

import (
	"time"

	"github.com/nextmv-io/sdk/nextroute/common"
)

// ModelStop is a stop to be routed.
type ModelStop interface {
	// Data returns the arbitrary data associated with the stop. Can be set
	// using the StopData StopOption.
	Data() any

	// Index returns the index of the stop.
	Index() int

	// Location returns the location of the stop.
	Location() common.Location

	// Name returns the name of the stop.
	Name() string

	// EarliestStart returns the earliest start time of the stop. Can be set
	// using the EarliestStart StopOption or using SetEarliestStart.
	EarliestStart() time.Time
	// EarliestStartValue returns the earliest start time of the stop as a
	// float64. The float64 value is the number of time units since the epoch
	// both possibly set at the construction of the Model. Can be set using
	// the EarliestStart StopOption or using SetEarliestStart.
	EarliestStartValue() float64

	// SetData sets the arbitrary data associated with the stop.
	SetData(data any)
	// SetEarliestStart sets the earliest start time of the stop.
	SetEarliestStart(time time.Time)
	// SetName sets the name of the stop.
	SetName(name string)
}

// ModelStops is a slice of stops.
type ModelStops []ModelStop
