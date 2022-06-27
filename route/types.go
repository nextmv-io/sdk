// Package types holds type definitions.
package route

import (
	"time"

	"github.com/nextmv-io/sdk/store"
)

// A Router is an engine that solves Vehicle Routing Problems.
type Router interface {
	// Options configures the router with the given options. An error is
	// returned if validation issues exist.
	Options(opts ...Option) error
	// Solver receives solve options and returns a Solver interface.
	Solver(opt store.Options) (store.Solver, error)
}

// An Option configures a router.
type Option func(r Router) error

// Stop to service in a Vehicle Routing Problem.
type Stop struct {
	ID       string   `json:"id"`
	Position Position `json:"position"`
}

// TimeWindow represents a time window for a shift of a vehicle or a stop.
type TimeWindow struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// Backlog represents the backlog, a list of stops for a vehicle.
type Backlog stopsToVehicle

// Alternate represents alternate stops, a list of stops for a vehicle.
type Alternate stopsToVehicle

// stopsToVehicle represents a relation between stops and a vehicle.
type stopsToVehicle struct {
	VehicleID string   `json:"id"` //nolint:tagliatelle
	Stops     []string `json:"stops"`
}

// Window represents a fixed timeframe in which the stop must be served. The
// duration represents the time it takes to service the stop. The max wait
// attribute defines what the allowed time is for vehicles arriving to stops
// before the window is open.
type Window struct {
	TimeWindow TimeWindow `json:"time_window"`
	MaxWait    int        `json:"max_wait"`
}

// ServiceGroup holds a group of stops and the service time duration (in
// seconds) to be added for every approach to one of the stops in the group.
type ServiceGroup struct {
	Group    []string `json:"group,omitempty"`
	Duration int      `json:"duration,omitempty"`
}

// Position represents a geographical location.
type Position struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

// Job represents a combination of one pick-up and one drop-off that must be
// served together with the pick-up preceding the drop-off.
type Job struct {
	PickUp  string `json:"pick_up,omitempty"`
	DropOff string `json:"drop_off,omitempty"`
}

// Attributes holds the ID of a vehicle or stop and corresponding compatibility
// attributes for that vehicle/stop.
type Attributes struct {
	ID         string   `json:"id"`
	Attributes []string `json:"attributes"`
}

// Service holds the ID of a stop and corresponding time to service the stop
// in seconds.
type Service struct {
	ID       string `json:"id"`
	Duration int    `json:"duration"`
}

// Limit holds a measure which will be limited by the given value.
type Limit struct {
	Measure ByIndex
	Value   float64
}

// Point represents a point in space. It may have any dimension.
type Point []float64

// ByIndex estimates the cost of going from one index to another.
type ByIndex interface {
	// Cost estimates the cost of going from one index to another.
	Cost(from, to int) float64
}

// ByPoint estimates the cost of going from one point to another.
type ByPoint interface {
	// Cost estimates the cost of going from one point to another.
	Cost(from, to Point) float64
}

// Triangular indicates that the triangle inequality holds for every
// measure that implements it.
type Triangular interface {
	Triangular() bool
}
