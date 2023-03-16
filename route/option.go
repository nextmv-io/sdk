package route

import (
	"math/rand"

	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/model"
)

// An Option configures a Router.
type Option func(Router) error

// Starts sets the starting locations indexed by vehicle. The length must match
// the vehicles' length.
func Starts(starts []Position) Option {
	connect.Connect(con, &startsFunc)
	return startsFunc(starts)
}

// Ends sets the ending locations indexed by vehicle. The length must match the
// vehicles' length.
func Ends(ends []Position) Option {
	connect.Connect(con, &endsFunc)
	return endsFunc(ends)
}

// InitializationCosts set the vehicle initialization costs indexed by vehicle.
// The length must match the vehicles' length.
func InitializationCosts(initializationCosts []float64) Option {
	connect.Connect(con, &initializationCostsFunc)
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
	connect.Connect(con, &capacityFunc)
	return capacityFunc(quantities, capacities)
}

// Precedence adds a precedence constraint to the list of constraints. It takes
// one argument as a slice of jobs. Each job defines a pick-up and drop-off by
// ID. The pick-up must precede the drop-off in the route.
func Precedence(precedences []Job) Option {
	connect.Connect(con, &precedenceFunc)
	return precedenceFunc(precedences)
}

// Services adds a service time to the given stops. For stops not in the slice,
// a service time of zero seconds is applied.
func Services(serviceTimes []Service) Option {
	connect.Connect(con, &servicesFunc)
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
	connect.Connect(con, &shiftsFunc)
	return shiftsFunc(shifts)
}

// Windows adds a time window constraint to the list of constraints. The method
// takes in windows, which are indexed by stop and represent a fixed timeframe
// in which the stop must be served. Service times at the stops can be
// optionally added using the Services option.
//
// PLEASE NOTE: this option requires using the Shift option.
func Windows(windows []Window) Option {
	connect.Connect(con, &windowsFunc)
	return windowsFunc(windows)
}

// MultiWindows adds a time window constraint to the list of constraints. The
// method takes in multiple windows per stop, which are indexed by stop and
// represent several fixed time frames in which the stop must be served.
// Furthermore, a wait time per stop must be specified where -1 means that a
// vehicle may wait indefinitely until a window opens. Service times at the
// stops can be optionally added using the Services option.
//
// PLEASE NOTE: this option requires using the Shift option.
func MultiWindows(windows [][]TimeWindow, maxWaitTimes []int) Option {
	connect.Connect(con, &multiWindowsFunc)
	return multiWindowsFunc(windows, maxWaitTimes)
}

// Unassigned sets the unassigned penalties indexed by stop. The length must
// match the stops' length.
func Unassigned(penalties []int) Option {
	connect.Connect(con, &unassignedFunc)
	return unassignedFunc(penalties)
}

// Backlogs sets the backlog for the specified vehicles. A backlog is an
// unordered list of stops that a vehicle has to serve.
func Backlogs(backlogs []Backlog) Option {
	connect.Connect(con, &backlogsFunc)
	return backlogsFunc(backlogs)
}

// Minimize sets the solver type of the router to minimize the value with a
// hybrid solver that uses decision diagrams and ALNS. This is the default
// solver that the router engine uses.
func Minimize() Option {
	connect.Connect(con, &minimizeFunc)
	return minimizeFunc()
}

// Maximize sets the solver type of the router to maximize the value with a
// hybrid solver that uses decision diagrams and ALNS.
func Maximize() Option {
	connect.Connect(con, &maximizeFunc)
	return maximizeFunc()
}

// Limits adds a route limit constraint to the list of constraints. The limit
// constraint can be used to limit the routes. The option takes two arguments:
// Firstly, the routeLimits struct which is indexed by vehicle and has two
// fields:
//   - The value in the unit of the given measure.
//   - The value to which the route is limited in the unit of the given measure.
//     To not limit the route to any value, use model.MaxInt
//
// Secondly, a flag to ignore the triangular inequality.
//
// PLEASE NOTE: If you want to limit the route's duration or length please use
// the options LimitDistance and LimitDuration, respectively.
func Limits(routeLimits []Limit, ignoreTriangular bool) Option {
	connect.Connect(con, &limitsFunc)
	return limitsFunc(routeLimits, ignoreTriangular)
}

// LimitDistances limits the distances of the routes by the given values. The
// values are indexed by vehicle and must be given in meters. To not limit a
// route to any value, use model.MaxInt.
func LimitDistances(maxDistances []float64, ignoreTriangular bool) Option {
	connect.Connect(con, &limitDistancesFunc)
	return limitDistancesFunc(maxDistances, ignoreTriangular)
}

// LimitDurations limits the the durations of the routes to the given values.
// The values are indexed by vehicle and must be given in seconds. To not limit
// a route to any value, use model.MaxInt.
func LimitDurations(maxDurations []float64, ignoreTriangular bool) Option {
	connect.Connect(con, &limitDurationsFunc)
	return limitDurationsFunc(maxDurations, ignoreTriangular)
}

// Grouper adds a custom location group to the list of location groups. When one
// or more groups of locations are defined, the router engine will ensure that
// all locations of a group will be assigned to the same route. If no groups are
// given, locations can be assigned together in the same routes without the need
// to assign any other locations to that route.
func Grouper(groups [][]string) Option {
	connect.Connect(con, &grouperFunc)
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
func ValueFunctionMeasures(
	valueFunctionMeasures []ByIndex,
) Option {
	connect.Connect(con, &valueFunctionMeasuresFunc)
	return valueFunctionMeasuresFunc(valueFunctionMeasures)
}

// ValueFunctionDependentMeasures sets custom dependent measures for every
// vehicle to calculate the overall solution costs, and should be indexed as
// such. If no custom measures are provided, a default haversine measure will be
// used to calculate costs (distances) between stops.
// Note that if your value function measures do represent time and you are using
// the window option with wait times, in order to see those wait times reflected
// in your value function, you will have to override the value function by using
// the `Update()` function. In the `Update()` function you can request a
// `TimeTracker` from the `state` and use it to get access to time information
// the route.
func ValueFunctionDependentMeasures(
	valueFunctionMeasures []DependentByIndex,
) Option {
	connect.Connect(con, &valueFunctionDependentMeasuresFunc)
	return valueFunctionDependentMeasuresFunc(valueFunctionMeasures)
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
	connect.Connect(con, &travelTimeMeasuresFunc)
	return travelTimeMeasuresFunc(timeMeasures)
}

/*
TravelTimeDependentMeasures sets custom dependent time measures for every
vehicle, and should be indexed as such. If no custom time measures are provided,
a default time measure will be used, based on haversine using a velocity of
10 m/s if no custom velocities are given using the Velocities option.
Time measures are used to:
  - calculate travel time in the Window option and check if time windows are
    met.
  - calculated route duration, the estimated time of arrival and departure at
    stops.

PLEASE NOTE: When defining a custom TravelTimeMeasure, this measure must not
account for any service times. To account for services times please use the
Services option.
*/
func TravelTimeDependentMeasures(timeMeasures []DependentByIndex) Option {
	connect.Connect(con, &travelTimeDependentMeasuresFunc)
	return travelTimeDependentMeasuresFunc(timeMeasures)
}

// TravelDistanceMeasures sets custom distance measures for every vehicle to
// calculate the travel distance, and should be indexed as such. If no
// distance measures are provided, a default haversine measure will be used to
// calculate distances between stops.
func TravelDistanceMeasures(
	travelDistanceMeasures []ByIndex,
) Option {
	connect.Connect(con, &travelDistanceMeasuresFunc)
	return travelDistanceMeasuresFunc(travelDistanceMeasures)
}

// Attribute sets a compatibility filter for stops and vehicles. It takes two
// arguments, vehicles and stops which define a slice of compatibility
// attributes for stops and vehicles. Stops that are not provided are compatible
// with any vehicle. Vehicles that are not provided are only compatible with
// stops without attributes.
func Attribute(vehicles []Attributes, stops []Attributes) Option {
	connect.Connect(con, &attributeFunc)
	return attributeFunc(vehicles, stops)
}

// Threads sets the number of threads that the internal solver uses. The router
// engine's solver is a hybrid solver that uses a Decision Diagram (DD) solver
// and various ALNS solvers with DD sub-solvers. If threads = 1, it means that
// only the first solver is used, which corresponds to pure DD. As a default,
// threads are calculated based on the number of machine processors, using
// the following expression: runtime.GOMAXPROCS(0) / 2.
func Threads(threads int) Option {
	connect.Connect(con, &threadsFunc)
	return threadsFunc(threads)
}

// Alternates sets a slice of alternate stops per vehicle. The vehicle will be
// assigned exactly one stop from the list of alternate stops, which are passed
// into this option, and any other stops from the list of stops that solve the
// TSP/VRP cost optimally.
func Alternates(alternates []Alternate) Option {
	connect.Connect(con, &alternatesFunc)
	return alternatesFunc(alternates)
}

// Velocities sets the speed for all vehicles to define a corresponding
// TravelTimeMeasure based on haversine distance and is indexed by vehicle. The
// length must match the vehicles' length.
func Velocities(velocities []float64) Option {
	connect.Connect(con, &velocitiesFunc)
	return velocitiesFunc(velocities)
}

// ServiceGroups adds an additional service time for a group of stops.
// Whenever a stop in the group is visited from another stop that is not part
// of it, the specified duration is added.
func ServiceGroups(serviceGroups []ServiceGroup) Option {
	connect.Connect(con, &serviceGroupsFunc)
	return serviceGroupsFunc(serviceGroups)
}

// Selector sets the given custom location selector. The location selector lets
// you define a function which selects the locations that will be inserted next
// into the solution. If no custom location selector is given, the location with
// the lowest index will be inserted next.
func Selector(selector func(PartialPlan) model.Domain) Option {
	connect.Connect(con, &selectorFunc)
	return selectorFunc(selector)
}

// VehicleUpdater defines an interface that is used to override the vehicle's
// default value function. The Update function takes a PartialVehicle as an
// input and returns three values, a VehicleUpdater with potentially updated
// bookkeeping data, a new solution value and a bool value to indicate wether
// the vehicle's solution value received an update.
// See the documentation of Update() for example usage.
type VehicleUpdater interface {
	Update(PartialVehicle) (VehicleUpdater, int, bool)
}

// PlanUpdater defines an interface that is used to override the router's
// default value function. The Update function takes a Plan and a slice of
// PartialVehicles as an input and returns three values, a PlanUpdater with
// potentially updated bookkeeping data, a new solution value and a bool value
// to indicate wether the vehicle's solution value received an update. The
// second parameter is a slice of PartialVehicles which may have been updated.
// All vehicles not part of that slice have definitely not changed. This
// knowledge can be used to more efficiently update the value of a plan. See the
// documentation of route.Update() for example usage.
type PlanUpdater interface {
	Update(PartialPlan, []PartialVehicle) (PlanUpdater, int, bool)
}

/*
Update sets the collection of functions that are called when transitioning from
one store to another in the router's Decision Diagram search for the best
solution in the time alloted. Updating information is useful for two purposes:
  - setting a custom value function (objective) that will be optimized.
  - bookkeeping of custom data.

The option takes the following arguments:
  - VehicleUpdater: replaces the value function of each vehicle. Can be nil
    if more than one vehicle is present.
  - PlanUpdater: replaces the value function of the full plan. Can be nil if
    only one vehicle is present.

User-defined custom types must implement the interfaces. When routing multiple
vehicles, the vehicleUpdater interface may be nil, if only information at the
fleet level is updated.

To customize the value function that will be optimized, the third parameter in
either of the interfaces must be true. If the last parameter is false, the
default value is used and it corresponds to the configured measure.

To achieve efficient customizations, always try to update the components of the
store that changed.
*/
func Update(v VehicleUpdater, f PlanUpdater) Option {
	connect.Connect(con, &updateFunc)
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
	connect.Connect(con, &filterWithRouteFunc)
	return filterWithRouteFunc(filter)
}

// Sorter sets the given custom vehicle sorter. The vehicle sorter lets you
// define a function which returns the vehicle indices in a specific order. The
// underlying engine will try to assign the locations to each vehicle in that
// returned order.
func Sorter(
	sorter func(
		p PartialPlan,
		locations model.Domain,
		vehicles model.Domain,
		random *rand.Rand,
	) []int,
) Option {
	connect.Connect(con, &sorterFunc)
	return sorterFunc(sorter)
}

// VehicleConstraint defines an interface that needs to be implemented when
// creating a custom vehicle constraint.
type VehicleConstraint interface {
	// Violated takes a PartialPlan and returns true if the vehicle with the
	// current route is not operationally valid. Often custom constraint hold
	// internal state. The first return value can be used to return a new
	// VehicleConstraint with updated state or nil in case the solution is not
	// operationally valid.
	Violated(PartialVehicle) (VehicleConstraint, bool)
}

// Constraint sets a custom constraint for specified vehicles. It takes two
// arguments, the constraint to be applied and a list of vehicles, indexed by
// ID, to which the constraint shall be applied.
func Constraint(constraint VehicleConstraint, ids []string) Option {
	connect.Connect(con, &constraintFunc)
	return constraintFunc(constraint, ids)
}

// Filter adds a custom vehicle filter to the list of filters. A filter checks
// which location is in general compatible with a vehicle. If no filter is given
// all locations are compatible with all vehicles and, thus, any location can be
// inserted into any route.
func Filter(compatible func(vehicle, location int) bool) Option {
	connect.Connect(con, &filterFunc)
	return filterFunc(compatible)
}

var (
	con                                = connect.NewConnector("sdk", "Route")
	startsFunc                         func([]Position) Option
	endsFunc                           func([]Position) Option
	capacityFunc                       func([]int, []int) Option
	initializationCostsFunc            func([]float64) Option
	precedenceFunc                     func([]Job) Option
	servicesFunc                       func([]Service) Option
	shiftsFunc                         func([]TimeWindow) Option
	windowsFunc                        func([]Window) Option
	multiWindowsFunc                   func([][]TimeWindow, []int) Option
	unassignedFunc                     func([]int) Option
	backlogsFunc                       func([]Backlog) Option
	minimizeFunc                       func() Option
	maximizeFunc                       func() Option
	limitsFunc                         func([]Limit, bool) Option
	limitDistancesFunc                 func([]float64, bool) Option
	limitDurationsFunc                 func([]float64, bool) Option
	grouperFunc                        func([][]string) Option
	valueFunctionMeasuresFunc          func([]ByIndex) Option
	valueFunctionDependentMeasuresFunc func([]DependentByIndex) Option
	travelTimeMeasuresFunc             func([]ByIndex) Option
	travelTimeDependentMeasuresFunc    func([]DependentByIndex) Option
	travelDistanceMeasuresFunc         func([]ByIndex) Option
	attributeFunc                      func([]Attributes, []Attributes) Option
	threadsFunc                        func(int) Option
	alternatesFunc                     func([]Alternate) Option
	velocitiesFunc                     func([]float64) Option
	serviceGroupsFunc                  func([]ServiceGroup) Option
	selectorFunc                       func(func(PartialPlan) model.Domain) Option
	updateFunc                         func(VehicleUpdater, PlanUpdater) Option
	filterWithRouteFunc                func(
		func(model.Domain, model.Domain, [][]int) model.Domain,
	) Option
	sorterFunc func(
		func(PartialPlan, model.Domain, model.Domain, *rand.Rand) []int,
	) Option
	constraintFunc func(VehicleConstraint, []string) Option
	filterFunc     func(func(int, int) bool) Option
)
