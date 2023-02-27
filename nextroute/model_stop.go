package nextroute

import (
	"time"

	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute/common"
)

// ModelStop is a stop to be routed.
type ModelStop interface {
	// Data returns the arbitrary data associated with the stop. Can be set
	// using the StopData StopOption.
	Data() interface{}

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

	// SetEarliestStart sets the earliest start time of the stop.
	SetEarliestStart(time time.Time)
}

// ModelStops is a slice of stops.
type ModelStops []ModelStop

// StopOption is an option for a stop. Can be used in the factory method of
// a stop Model.NewStop.
type StopOption func(ModelStop) error

// StopData is an option for a stop. Can be used in the factory method of
// a stop Model.NewStop. The data is arbitrary data associated with the
// stop.
func StopData(
	data any,
) StopOption {
	connect.Connect(con, &stopDataStopOption)
	return stopDataStopOption(data)
}

// EarliestStart is an option for a stop. Can be used in the factory method of
// a stop Model.NewStop. The earliest start time is the earliest time at which
// the stop can be started.
func EarliestStart(earliestStart time.Time) StopOption {
	connect.Connect(con, &earliestStartStopOption)
	return earliestStartStopOption(earliestStart)
}

// Name is an option for a stop. Can be used in the factory method of a stop
// Model.NewStop. The name is the name of the stop and is used for debugging
// and reporting.
func Name(name string) StopOption {
	connect.Connect(con, &nameStopOption)
	return nameStopOption(name)
}
