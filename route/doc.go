// Package route provides routing functionalities.
//
// Router
//
// The router API provides a convenient interface for solving vehicle routing
// problems. It employs a hybrid solver that relies on decision diagrams and
// ALNS. To use it, invoke the router by passing stops and vehicles.
//
//     router, err := NewRouter(stops, vehicles)
//
// The router can be configured through options. It is composable, meaning that
// several options (or none at all) could be used. Every option, unless
// otherwise noted, can be used independently of others. An example of this is
// setting vehicle start locations.
//
//     router, err := NewRouter(stops, vehicles, Starts(starts))
//
// For convenience, options can also be configured after the router is declared
// through the Options function. An example of this is setting vehicle end
// locations.
//
//     err := Options(Ends(ends))
//
// The Solver function is used to obtain the solver that searches for a
// solution.
//
//     solver, err := Solver(opt)
//
// Measures
//
// Routing models frequently need to determine the cost of connecting two things
// together. This may mean assigning one item to another, clustering two points
// together, or routing a vehicle from one location to another. These cost
// computations are generically referred to as "measures". The package provides
// a number of common patterns for constructing and using them inside models.
//
// Point-to-Point Measures
//
// When cost must be computed based on distance between two points, a model can
// use a ByPoint implementation. These is the case for models where points are
// determined dynamically within the model logic, such as in k-means clustering.
// Such measures map two points to a cost.
//
//     cost := m.Cost(fromPoint, toPoint)
//
// The following ByPoint implementations are available.
//
//     EuclideanByPoint: Euclidean distance between two points
//     HaversineByPoint: Haversine distance between two points
//     TaxicabByPoint:   Taxicab distance between two points
//
// Points may be of any dimension. If the points passed in to any of these
// measures have differing dimensionality, they will project the lower dimension
// point into the higher dimension by appending 0s.
//
// Indexed Measures
//
// Models that do not require points operate on indices. These indices may or
// may not refer to points. An ByIndex implementation provides the same
// functionality as a ByPoint implementation, except its cost method accepts two
// indices instead of two points.
//
//     cost := m.Cost(fromIndex, toIndex)
//
// Index measures are more common, and a number of them embed and operate on
// results from other index measures.
//
// The following ByIndex implementations are available.
//
//     Bin:      select from a slice of measure by some function
//     Location: adds fixed location costs to another measure
//     Constant: always returns the same cost
//     Matrix:   looks up cost from a row to a column index
//     Override: overrides some other measure given a condition
//     Power:    takes some other measure to a power
//     Scale:    scales some other measure by a constant
//     Sparse:   sparse matrix measure with a backup
//     Sum:      adds the costs of other measures together
//     Truncate: truncates cost values provided by another measure
//     Location: adds cost of visiting a location to another measure
//
// In addition, the package provides Indexed, which adapts any ByPoint into a
// ByIndex. In addition to the ByPoint to be converted, Indexed accepts a fixed
// slice of points that it will use to look up the positions of indexes passed
// to Cost.
//
package route
