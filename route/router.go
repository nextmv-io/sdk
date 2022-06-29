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
type Backlog struct {
	VehicleID string   `json:"vehicle_id"`
	Stops     []string `json:"stops"`
}

// Alternate represents alternate stops, a list of stops for a vehicle.
type Alternate struct {
	VehicleID string   `json:"vehicle_id"`
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

// NewRouter returns a router interface. It receives a set of stops that must be
// serviced by a fleet of vehicles and a list of options. When an option is
// applied, an error is returned if there are validation issues. The router is
// composable, meaning that several options may be used or none at all. The
// options, unless otherwise noted, can be used independently of each other.
func NewRouter(
	stops []Stop,
	vehicles []string,
	opts ...Option,
) (Router, error) {
	connect()
	return newRouterFunc(stops, vehicles, opts...)
}

var newRouterFunc func([]Stop, []string, ...Option) (Router, error)
