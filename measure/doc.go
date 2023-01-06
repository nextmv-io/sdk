// Package measure contains cost measures based on indices or points.
//
// # Measures
//
// Engines models frequently need to determine the cost of connecting two things
// together. This may mean assigning one item to another, clustering two points
// together, or routing a vehicle from one location to another. These cost
// computations are generically referred to as "measures". Engines provides a
// number of common patterns for constructing and using them inside models.
//
// # Point-to-Point Measures
//
// When cost must be computed based on distance between two points, a model can
// use a measure.ByPoint implementation. These is the case for models where
// points are determined dynamically within the model logic, such as in k-means
// clustering. Such measures map two points to a cost.
//
//	cost := m.Cost(fromPoint, toPoint)
//
// The following measure.ByPoint implementations are available.
//
//	measure.EuclideanByPoint: Euclidean distance between two points
//	measure.HaversineByPoint: Haversine distance between two points
//	measure.TaxicabByPoint:   Taxicab distance between two points
//
// Unassigned may be of any dimension. If the points passed in to any of these
// measures have differing dimensionality, they will project the lower dimension
// point into the higher dimension by appending 0s.
//
// # Indexed Measures
//
// Models that do not require points operate on indices. These indices may or
// may not refer to points. An measure.ByIndex implementation provides the same
// functionality as a measure.ByPoint implementation, except its cost method
// accepts two indices instead of two points.
//
//	cost := m.Cost(fromIndex, toIndex)
//
// Index measures are more common, and a number of them embed and operate on
// results from other index measures.
//
// The following measure.ByIndex implementations are available.
//
//	measure.Bin:      select from a slice of measure by some function
//	measure.Location: adds fixed location costs to another measure
//	measure.Constant: always returns the same cost
//	measure.Matrix:   looks up cost from a row to a column index
//	measure.Override: overrides some other measure given a condition
//	measure.Power:    takes some other measure to a power
//	measure.Scale:    scales some other measure by a constant
//	measure.Sparse:   sparse matrix measure with a backup
//	measure.Sum:      adds the costs of other measures together
//	measure.Truncate: truncates cost values provided by another measure
//	measure.Location: adds cost of visiting a location to another measure
//
// In addition, Engines provides measure.Indexed, which adapts any
// measure.ByPoint into a measure.ByIndex. In addition to the measure.ByPoint to
// be converted, Indexed accepts a fixed slice of points that it will use to
// look up the positions of indexes passed to Cost.
package measure
