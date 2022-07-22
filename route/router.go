package route

import (
	"time"

	"github.com/nextmv-io/sdk/model"
	"github.com/nextmv-io/sdk/store"
)

// A Router is an engine that solves Vehicle Routing Problems.
type Router interface {
	// Options configures the router with the given options. An error is
	// returned if validation issues exist.
	Options(...Option) error
	// Solver receives solve options and returns a Solver interface.
	Solver(store.Options) (store.Solver, error)
	/*
		Plan returns a variable which holds information about the current set of
		vehicles with their respective routes and any unassigned stops. The Plan
		variable can be used to retrieve the values from the Store of a
		Solution.

			router, err := route.NewRouter(
				stops,
				vehicles,
				route.Capacity(quantities, capacities),
			)
			if err != nil {
				panic(err)
			}
			solver, err := router.Solver(store.DefaultOptions())
			if err != nil {
				panic(err)
			}
			solution := solver.Last(context.Background())
			s := solution.Store
			p := router.Plan()
			vehicles, unassigned := p.Get(s).Vehicles, p.Get(s).Unassigned
	*/
	Plan() store.Variable[Plan]
}

// PartialPlan is an (incomplete) Plan that operates on the internal
// solver data structures. Certain router options that customize solver
// internals have to work with this data structure.
type PartialPlan interface {
	// Unassigned returns an Integer Domain with unassigned stop indexes.
	// These are stops explicitly excluded from being served by a vehicle.
	Unassigned() model.Domain
	// Unplanned returns an Integer Domain with not yet assigned or unassigned
	// stops indexes.
	Unplanned() model.Domain
	// Value return the value of this plan.
	Value() int
	// Vehicles returns a slice of vehicles part of this partial plan.
	Vehicles() []PartialVehicle
}

// PartialVehicle represents a Vehicle that operates on the internal solver
// data structures. Certain router options that customize solver internals have
// to work with this data structure.
type PartialVehicle interface {
	// ID returns the vehicle ID.
	ID() string
	// Updater returns either nil in case no custom VehicleUpdater was used or
	// the custom VehicleUpdater that was used for this vehicle.
	Updater() VehicleUpdater
	// Route returns the route of the vehicle represented by a sequence of stop
	// indexes. The first and last indexes are always the starting and ending
	// locations of the vehicle, respectively.
	Route() []int
	// Value return the value of vehicle. Usually this is the cost of the route.
	Value() int
}

// Plan describes a solution to a Vehicle Routing Problem.
type Plan struct {
	Unassigned []Stop           `json:"unassigned"`
	Vehicles   []PlannedVehicle `json:"vehicles"`
}

// PlannedVehicle holds information about the vehicle in a solution to a Vehicle
// Routing Problem.
type PlannedVehicle struct {
	ID            string        `json:"id"`
	Route         []PlannedStop `json:"route"`
	RouteDuration int           `json:"route_duration"`
}

// PlannedStop describes a stop as part of a Vehicle's route of solution
// to a Vehicle Routing Problem.
type PlannedStop struct {
	Stop
	EstimatedArrival   *time.Time `json:"estimated_arrival,omitempty"`
	EstimatedDeparture *time.Time `json:"estimated_departure,omitempty"`
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
