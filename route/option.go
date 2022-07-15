package route

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
func Limits(routeLimits []Limit, ignoreTriangular bool) Option {
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
)
