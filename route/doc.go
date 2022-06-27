/*
Package route provides routing engines.

Routing models direct vehicles from some start location to an end location,
servicing other locations along the way. These locations often represent
customer requests, though they may logically correspond to any number of things.
Such models attempt to service all locations optimizing a value, such as
minimizing cost.

The following routing engines are provided.
	- engines/route/vehicle: single-vehicle routing.
	- engines/route/fleet:   multi-vehicle fleet routing.
	- engines/route:         Router API for multi- or single-vehicle routing.

Vehicle

Single-vehicle models route one vehicle from a start location through a known
domain of locations to an end location. The cost of a solution is the sum of its
arc costs. Routes may be constrained to enforce precedence relationships,
capacity limits, time windows, and user-defined constraints.

Fleet

Fleet routing models generalize single-vehicle routing to multiple vehicles.
Locations are assigned to drivers and routed using the Vehicle engine. The cost
of a solution is the sum of its collective routes.

Router

The router API uses the vehicle and fleet engines to provide a convenient
interface for solving vehicle routing problems. It employs a hybrid solver that
relies on decision diagrams and ALNS. To use it, invoke the router by passing
stops and vehicles.

    router, err := route.NewRouter(stops, vehicles)

The router can be configured through options. It is composable, meaning that
several options (or none at all) could be used. Every option, unless otherwise
noted, can be used independently of others. An example of this is setting
vehicle start locations.

	router, err := route.NewRouter(stops, vehicles, route.Starts(starts))

For convenience, options can also be configured after the router is declared
through the Options function. An example of this is setting vehicle end
locations.

	err := router.Options(route.Ends(ends))

The Solver function is used to obtain the solver that searches for a solution.

	solver, err := router.Solver(opt)
*/
package route
