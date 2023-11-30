/*
Package route provides vehicle routing functionalities. Routing models direct
vehicles from some start location to an end location, servicing other locations
along the way. These locations often represent customer requests, though they
may logically correspond to any number of things. Such models attempt to
service all locations optimizing a value, such as minimizing cost.

# Router

The Router API provides a convenient interface for solving vehicle routing
problems. It employs a hybrid solver that relies on decision diagrams and
ALNS. To use it, invoke the Router by passing stops and vehicles.

	router, err := route.NewRouter(stops, vehicles)

The Router can be configured through options. It is composable, meaning that
several options (or none at all) could be used. Every option, unless
otherwise noted, can be used independently of others. An example of this is
setting vehicle start locations.

	router, err := route.NewRouter(stops, vehicles, route.Starts(starts))

For convenience, options can also be configured after the Router is declared
through the Options function. An example of this is setting vehicle end
locations.

	err := router.Options(route.Ends(ends))

The Router is built on top of the `store` package for solving decision
automation problems. As such, the Solver function is used to obtain the Solver
that searches for a Solution.

	solver, err := router.Solver(store.DefaultOptions())

Retrieve the routing plan -made up of the vehicle routes and any unassigned
stops- directly by calling the Last Solution on the Solver, for example.

	solution := solver.Last(context.Background())
	s := solution.Store
	plan := router.Plan()
	vehicles := plan.Get(s).Vehicles // On each vehicle call .Route
	unassigned := plan.Get(s).Unassigned

The Router can be paired with a Runner for convenience, to read data and
options and manage the call to the Solver directly. Please see the
documentation of the `store` and `run` packages for more information.

	package main

	import (
		"github.com/nextmv-io/sdk/route"
		"github.com/nextmv-io/sdk/run"
		"github.com/nextmv-io/sdk/store"
	)

	type input struct {
		Stops       []route.Stop `json:"stops"`
		Vehicles    []string     `json:"vehicles"`
		Capacities  []int        `json:"capacities"`
		Quantities  []int        `json:"quantities"`
		Precedences []route.Job  `json:"precedences"`
	}

	func main() {
		handler := func(i input, opt store.Options) (store.Solver, error) {
			router, err := route.NewRouter(
				i.Stops, i.Vehicles,
				// Add all options, or none.
				route.Capacity(i.Quantities, i.Capacities),
				route.Precedence(i.Precedences),
			)
			if err != nil {
				return nil, err
			}
			return router.Solver(opt) // Options are passed by the runner.
		}
		run.Run(handler)
	}

Given that the Router works with a Store (used for any type of decisions),
vehicle routing problems can be embedded into broader decision problems. E.g.:
determine the number of vehicles that minimize the cost of routing. The Store
defines a Variable x holding the number of vehicles. To estimate the cost of
adding an additional vehicle, the Router is used to estimate the total cost.
The problem is operationally valid if all stops are assigned.

	package main

	import (
		"context"
		"strconv"

		"github.com/nextmv-io/sdk/route"
		"github.com/nextmv-io/sdk/run"
		"github.com/nextmv-io/sdk/store"
	)

	type input struct {
		Stops []route.Stop `json:"stops"`
	}

	// Rough pseudo-code that shouldn't be run.
	func main() {
		handler := func(i input, opt store.Options) (store.Solver, error) {
			// Declare a new store of variables.
			s := store.New()

			// Outer decision: number of vehicles represented by an integer
			// variable.
			x := store.NewVar(s, 0)

			// Modify the main store.
			s = s.
				Value(func(s store.Store) int {
					// Each vehicle has a cost => the value of the store is the
					// total cost.
					return x.Get(s) * 99
				}).
				Generate(func(s store.Store) store.Generator {
					// Current number of vehicles.
					vehicles := make([]string, x.Get(s))
					for i := 0; i < x.Get(s); i++ {
						vehicles[i] = strconv.Itoa(i)
					}

					// Embedded decision: summon the router and its solver.
					router, err := route.NewRouter(
						i.Stops, vehicles,
					)
					if err != nil {
						return nil
					}
					solver, err := router.Solver(opt)
					if err != nil {
						return nil
					}

					// Get the solution to the routing problem. Estimate the
					// cost and obtain the unassigned stops.
					solution := solver.Last(context.Background())
					plan := router.Plan()
					unassigned := plan.Get(solution.Store).Unassigned
					routingCost := solution.Statistics.Value

					return store.Lazy(
						func() bool {
							// Use some generating condition, such as generating
							// children stores as long as there are unassigned
							// stops.
							return len(unassigned) > 0
						},
						func() store.Store {
							return s.
								// Add another vehicle.
								Apply(x.Set(x.Get(s) + 1)).
								// Modify the value of the Store.
								Value(func(s store.Store) int { return *routingCost }).
								// Operationally valid if all stops are assigned.
								Validate(func(s store.Store) bool {
									return len(unassigned) == 0
								})
						},
					)
				})

			// Minimize the total cost.
			return s.Minimizer(opt), nil
		}
		run.Run(handler)
	}

# Measures

Routing models frequently need to determine the cost of connecting two things
together. This may mean assigning one item to another, clustering two points
together, or routing a vehicle from one location to another. These cost
computations are generically referred to as "measures". The package provides
a number of common patterns for constructing and using them inside models.

# Point-to-Point Measures

When cost must be computed based on distance between two points, a model can
use a ByPoint implementation. These is the case for models where points are
determined dynamically within the model logic, such as in k-means clustering.
Such measures map two points to a cost.

	cost := m.Cost(fromPoint, toPoint)

The following ByPoint implementations are available.

	EuclideanByPoint: Euclidean distance between two points
	HaversineByPoint: Haversine distance between two points
	TaxicabByPoint:   Taxicab distance between two points

Points may be of any dimension. If the points passed in to any of these
measures have differing dimensionality, they will project the lower dimension
point into the higher dimension by appending 0s.

# Indexed Measures

Models that do not require points operate on indices. These indices may or
may not refer to points. An ByIndex implementation provides the same
functionality as a ByPoint implementation, except its cost method accepts two
indices instead of two points.

	cost := m.Cost(fromIndex, toIndex)

Index measures are more common, and a number of them embed and operate on
results from other index measures.

The following ByIndex implementations are available.

	Bin:      select from a slice of measure by some function
	Location: adds fixed location costs to another measure
	Constant: always returns the same cost
	Matrix:   looks up cost from a row to a column index
	Override: overrides some other measure given a condition
	Power:    takes some other measure to a power
	Scale:    scales some other measure by a constant
	Sparse:   sparse matrix measure with a backup
	Sum:      adds the costs of other measures together
	Truncate: truncates cost values provided by another measure
	Location: adds cost of visiting a location to another measure

In addition, the package provides Indexed, which adapts any ByPoint into a
ByIndex. In addition to the ByPoint to be converted, Indexed accepts a fixed
slice of points that it will use to look up the positions of indices passed
to Cost.

Deprecated: This package is deprecated and will be removed in a future.
Use [github.com/nextmv-io/sdk/nextroute] instead.
*/
package route
