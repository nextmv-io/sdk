package route

import (
	"github.com/nextmv-io/sdk/model"
)

// An Option configures a Router.
type Option func(Router) error

// Starts sets the starting locations indexed by vehicle. The length must match
// the vehicles' length.
func Starts(starts []Position) Option {
	connect()
	return startsFunc(starts)
}

// Ends sets the ending locations indexed by vehicle. The length must match the
// vehicles' length.
func Ends(ends []Position) Option {
	connect()
	return endsFunc(ends)
}

// InitializationCosts set the vehicle initialization costs indexed by vehicle.
// The length must match the vehicles' length.
func InitializationCosts(initializationCosts []float64) Option {
	connect()
	return initializationCostsFunc(initializationCosts)
}

// Capacity adds a capacity constraint to the list of constraints and takes two
// arguments: quantities and capacities. Quantities represent the change in
// vehicle capacity indexed by stop. Capacities represent the maximum capacity
// indexed by vehicle. The quantities and capacities must match the length of
// stops and vehicles, respectively. To specify multiple capacity constraints,
// this option may be used several times with the corresponding quantities and
// capacities.
func Capacity(quantities []int, capacities []int) Option {
	connect()
	return capacityFunc(quantities, capacities)
}

// Precedence adds a precedence constraint to the list of constraints. It takes
// one argument as a slice of jobs. Each job defines a pick-up and drop-off by
// ID. The pick-up must precede the drop-off in the route.
func Precedence(precedences []Job) Option {
	connect()
	return precedenceFunc(precedences)
}

// Services adds a service time to the given stops. For stops not in the slice,
// a service time of zero seconds is applied.
func Services(serviceTimes []Service) Option {
	connect()
	return servicesFunc(serviceTimes)
}

/*
Shifts adds shifts to the vehicles. Shifts are indexed by vehicle and represent
a time window on its shift's start and end time. When using the Windows option,
using the Shifts option is required. Shifts are additionally used to:
	- enable the calculation of the estimated arrival and departure at stops
	- set the start time of the route when using the Windows option, given that
	time tracking is needed for the Window constraint.
If the Measures option is not used, a default Haversine measure will be used. If
the TimeMeasures option is not used, the measures will be scaled with a constant
velocity of 10 m/s.
*/
func Shifts(shifts []TimeWindow) Option {
	connect()
	return shiftsFunc(shifts)
}

// Windows adds a time window constraint to the list of constraints. The method
// takes in windows, which are indexed by stop and represent a fixed timeframe
// in which the stop must be served. Service times at the stops can be
// optionally added using the Services option.
//
// PLEASE NOTE: this option requires using the Shift option.
func Windows(windows []Window) Option {
	connect()
	return windowsFunc(windows)
}

// Unassigned sets the unassigned penalties indexed by stop. The length must
// match the stops' length.
func Unassigned(penalties []int) Option {
	connect()
	return unassignedFunc(penalties)
}

// Backlogs sets the backlog for the specified vehicles. A backlog is an
// unordered list of stops that a vehicle has to serve.
func Backlogs(backlogs []Backlog) Option {
	connect()
	return backlogsFunc(backlogs)
}

// Minimize sets the solver type of the router to minimize the value with a
// hybrid solver that uses decision diagrams and ALNS. This is the default
// solver that the router engine uses.
func Minimize() Option {
	connect()
	return minimizeFunc()
}

// Maximize sets the solver type of the router to maximize the value with a
// hybrid solver that uses decision diagrams and ALNS.
func Maximize() Option {
	connect()
	return maximizeFunc()
}

// Limits adds a route limit constraint to the list of constraints. The limit
// constraint can be used to limit the routes. The option takes two arguments:
// Firstly, the routeLimits struct which is indexed by vehicle and has two
// fields:
//  - The value in the unit of the given measure.
//  - The value to which the route is limited in the unit of the given measure.
//    To not limit the route to any value, use model.MaxInt
// Secondly, a flag to ignore the triangular inequality.
//
// PLEASE NOTE: If you want to limit the route's duration or length please use
// the options LimitDistance and LimitDuration, respectively.
func Limits(
	routeLimits []Limit,
	ignoreTriangular bool,
) Option {
	connect()
	return limitsFunc(routeLimits, ignoreTriangular)
}

// LimitDistances limits the distances of the routes by the given values. The
// values are indexed by vehicle and must be given in meters. To not limit a
// route to any value, use model.MaxInt.
func LimitDistances(maxDistances []float64, ignoreTriangular bool) Option {
	connect()
	return limitDistancesFunc(maxDistances, ignoreTriangular)
}

// LimitDurations limits the the durations of the routes to the given values.
// The values are indexed by vehicle and must be given in seconds. To not limit
// a route to any value, use model.MaxInt.
func LimitDurations(maxDurations []float64, ignoreTriangular bool) Option {
	connect()
	return limitDurationsFunc(maxDurations, ignoreTriangular)
}

// Grouper adds a custom location group to the list of location groups. When one
// or more groups of locations are defined, the router engine will ensure that
// all locations of a group will be assigned to the same route. If no groups are
// given, locations can be assigned together in the same routes without the need
// to assign any other locations to that route.
func Grouper(groups [][]string) Option {
	connect()
	return grouperFunc(groups)
}

// ValueFunctionMeasures sets custom measures for every vehicle to calculate the
// overall solution costs, and should be indexed as such. If no custom measures
// are provided, a default haversine measure will be used to calculate costs
// (distances) between stops.
// Note that if your value function measures do represent time and you are using
// the window option with wait times, in order to see those wait times reflected
// in your value function, you will have to override the value function by using
// the `Update()` function. In the `Update()` function you can request a
// `TimeTracker` from the `state` and use it to get access to time information
// the route.
func ValueFunctionMeasures(valueFunctionMeasures []ByIndex) Option {
	connect()
	return valueFunctionMeasuresFunc(valueFunctionMeasures)
}

/*
TravelTimeMeasures sets custom time measures for every vehicle, and should be
indexed as such. If no custom time measures are provided, a default time measure
will be used, based on haversine using a velocity of 10 m/s if no custom
velocities are given using the Velocities option. Time measures are used to:
	- calculate travel time in the Window option and check if time windows are
	met.
	- calculated route duration, the estimated time of arrival and departure at
	stops.

PLEASE NOTE: When defining a custom TravelTimeMeasure, this measure must not
account for any service times. To account for services times please use the
Services option.
*/
func TravelTimeMeasures(timeMeasures []ByIndex) Option {
	connect()
	return travelTimeMeasuresFunc(timeMeasures)
}

// Attribute sets a compatibility filter for stops and vehicles. It takes two
// arguments, vehicles and stops which define a slice of compatibility
// attributes for stops and vehicles. Stops that are not provided are compatible
// with any vehicle. Vehicles that are not provided are only compatible with
// stops without attributes.
func Attribute(vehicles []Attributes, stops []Attributes) Option {
	connect()
	return attributeFunc(vehicles, stops)
}

// Threads sets the number of threads that the internal solver uses. The router
// engine's solver is a hybrid solver that uses a Decision Diagram (DD) solver
// and various ALNS solvers with DD sub-solvers. If threads = 1, it means that
// only the first solver is used, which corresponds to pure DD. As a default,
// threads are calculated based on the number of machine processors, using
// the following expression: runtime.GOMAXPROCS(0) / 2.
func Threads(threads int) Option {
	connect()
	return threadsFunc(threads)
}

// Alternates sets a slice of alternate stops per vehicle. The vehicle will be
// assigned exactly one stop from the list of alternate stops, which are passed
// into this option, and any other stops from the list of stops that solve the
// TSP/VRP cost optimally.
func Alternates(alternates []Alternate) Option {
	connect()
	return alternatesFunc(alternates)
}

// Velocities sets the speed for all vehicles to define a corresponding
// TravelTimeMeasure based on haversine distance and is indexed by vehicle. The
// length must match the vehicles' length.
func Velocities(velocities []float64) Option {
	connect()
	return velocitiesFunc(velocities)
}

// ServiceGroups adds an additional service time for a group of stops.
// Whenever a stop in the group is visited from another stop that is not part
// of it, the specified duration is added.
func ServiceGroups(serviceGroups []ServiceGroup) Option {
	connect()
	return serviceGroupsFunc(serviceGroups)
}

// Selector sets the given custom location selector. The location selector lets
// you define a function which selects the location that will be inserted next
// into the solution. If no custom location selector is given, the location with
// the lowest index will be inserted next.
func Selector(selector func(FleetPlan) model.Domain) Option {
	connect()
	return selectorFunc(selector)
}

// FleetVehicleUpdater defines an interface that is used to override the
// vehicle's default value function. It requires the implementation of two
// functions, Value and Clone. The Clone function is used to make deep copies of
// bookkeeping data which will be used to properly update the solution's value
// in the Value function. The Value function takes a FleetVehicle as an input
// and returns two values, a new solution value and a bool value to indicate
// wether the vehicle's solution value received an update.
type FleetVehicleUpdater interface {
	Value(FleetVehicle) (int, bool)
	Clone() FleetVehicleUpdater
}

// FleetUpdater defines an interface that is used to override the fleet's
// default value function. It requires the implementation of two functions,
// Value and Clone. The Clone function is used to make deep copies of
// bookkeeping data which will be used to properly update the solution's value
// in the Value function. The Value function takes a FleetPlan and a slice of
// FleetVehicles as an input and returns two values, a new solution value and a
// bool value to indicate wether the vehicle's solution value received an
// update.
type FleetUpdater interface {
	Value(FleetPlan, []FleetVehicle) (int, bool)
	Clone() FleetUpdater
}

/*
Update sets the collection of functions that are called when transitioning from
one state to another in the router engine's Decision Diagram search for the best
solution in the time alloted. Updating information is useful for two purposes:
	- setting a custom value function (objective) that will be optimized.
	- book-keeping of custom data.

The option takes the following arguments:
	- FleetVehicleUpdater: updates the vehicle.
	- FleetUpdater: updates the fleet.

User-defined custom types must implement the interfaces.
*/
func Update(v FleetVehicleUpdater, f FleetUpdater) Option {
	connect()
	return updateFunc(v, f)
}

// FilterWithRoute adds a new VehicleFilter. Compared to the Filter option, the
// FilterWithRoute option is more flexible. It defines a function that takes an
// IntDomain of candidate vehicles, an IntDomain of locations that will be
// assigned to a particular vehicle, and a slice of routes for all vehicles. It
// returns an IntDomain representing vehicles that cannot service the domain of
// locations.
func FilterWithRoute(
	filter func(
		vehicleCandidates model.Domain,
		locations model.Domain,
		routes [][]int,
	) model.Domain,
) Option {
	connect()
	return filterWithRouteFunc(filter)
}

// Sorter sets the given custom vehicle sorter. The vehicle sorter lets you
// define a function which returns the vehicle indices in a specific order. The
// underlying engine will try to assign the locations to each vehicle in that
// returned order.
func Sorter(sorter func(
	p FleetPlan,
	locations model.Domain,
	vehicles model.Domain,
) []int,
) Option {
	connect()
	return sorterFunc(sorter)
}

// Constraint sets a custom constraint for specified vehicles. It takes two
// arguments, the constraint to be applied and a list of vehicles, indexed by
// ID, to which the constraint shall be applied.
func Constraint(constraint VehicleConstraint, ids []string) Option {
	connect()
	return constraintFunc(constraint, ids)
}

// VehicleConstraint defines an interface that needs to be implemented when
// creating a custom vehicle constraint.
type VehicleConstraint interface {
	Violated(FleetVehicle) (VehicleConstraint, bool)
}

var (
	startsFunc                func([]Position) Option
	endsFunc                  func([]Position) Option
	capacityFunc              func([]int, []int) Option
	initializationCostsFunc   func([]float64) Option
	precedenceFunc            func([]Job) Option
	servicesFunc              func([]Service) Option
	shiftsFunc                func([]TimeWindow) Option
	windowsFunc               func([]Window) Option
	unassignedFunc            func([]int) Option
	backlogsFunc              func([]Backlog) Option
	minimizeFunc              func() Option
	maximizeFunc              func() Option
	limitsFunc                func([]Limit, bool) Option
	limitDistancesFunc        func([]float64, bool) Option
	limitDurationsFunc        func([]float64, bool) Option
	grouperFunc               func([][]string) Option
	valueFunctionMeasuresFunc func([]ByIndex) Option
	travelTimeMeasuresFunc    func([]ByIndex) Option
	attributeFunc             func([]Attributes, []Attributes) Option
	threadsFunc               func(int) Option
	alternatesFunc            func([]Alternate) Option
	velocitiesFunc            func([]float64) Option
	serviceGroupsFunc         func([]ServiceGroup) Option
	selectorFunc              func(func(FleetPlan) model.Domain) Option
	updateFunc                func(FleetVehicleUpdater, FleetUpdater) Option
	filterWithRouteFunc       func(
		func(model.Domain, model.Domain, [][]int) model.Domain,
	) Option
	sorterFunc func(
		func(FleetPlan, model.Domain, model.Domain) []int,
	) Option
	constraintFunc func(VehicleConstraint, []string) Option
)
