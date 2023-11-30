package route

import (
	"sort"

	"github.com/nextmv-io/sdk/measure"
)

// Point represents a point in space. It may have any dimension.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
type Point = measure.Point

// ByIndex estimates the cost of going from one index to another.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
type ByIndex = measure.ByIndex

// DependentByIndex is a measure uses a custom cost func to calculate parameter
// dependent costs for connecting to points by index.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
type DependentByIndex = measure.DependentByIndex

// ByPoint estimates the cost of going from one point to another.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
type ByPoint = measure.ByPoint

// Triangular indicates that the triangle inequality holds for every
// measure that implements it.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
type Triangular = measure.Triangular

// DurationGroups represents a slice of duration groups. Each duration group is
// used to account for additional service costs whenever a stop of a group is
// approached first.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
type DurationGroups = measure.DurationGroups

// DurationGroup groups stops by index which have additional service costs
// attached to them.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
type DurationGroup = measure.DurationGroup

// Times holds the estimated time of arrival (ETA), the estimated time of when
// service starts (ETS) and estimated time of departure (ETD).
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
type Times = measure.Times

// VehicleData holds vehicle specific data, including times by index (ETA, ETD
// and ETS), a vehicle id, the vehicle's route and the solution value for that
// vehicle.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
type VehicleData = measure.VehicleData

// Bin is a measure that selects from a slice of indexed measures. Logic
// defined in the selector function determines which measure is used in the
// cost calculation.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
func Bin(
	measures []measure.ByIndex,
	selector func(from, to int) int,
) ByIndex {
	return measure.Bin(measures, selector)
}

// Override measure uses a default measure for all arcs that are not true for a
// condition. It uses an override measure for all arcs that are true for the
// condition.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
func Override(
	defaultByIndex ByIndex,
	overrideByIndex ByIndex,
	condition func(from, to int) bool,
) ByIndex {
	return measure.Override(defaultByIndex, overrideByIndex, condition)
}

// DebugOverride returns an Override that when marshalled will include debugging
// information describing the number of queries for default and override
// elements.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
func DebugOverride(
	defaultByIndex ByIndex,
	overrideByIndex ByIndex,
	condition func(from, to int) bool,
) ByIndex {
	return measure.DebugOverride(defaultByIndex, overrideByIndex, condition)
}

// Power raises the cost of some other measure to an exponent.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
func Power(m ByIndex, exponent float64) ByIndex {
	return measure.Power(m, exponent)
}

// HaversineByPoint estimates meters connecting two points along the surface
// of the earth.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
func HaversineByPoint() ByPoint {
	return measure.HaversineByPoint()
}

// ConstantByPoint measure always estimates the same cost.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
func ConstantByPoint(c float64) ByPoint {
	return measure.ConstantByPoint(c)
}

// Constant measure always estimates the same cost.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
func Constant(c float64) ByIndex {
	return measure.Constant(c)
}

// EuclideanByPoint computes straight line distance connecting two indices.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
func EuclideanByPoint() ByPoint {
	return measure.EuclideanByPoint()
}

// Indexed creates a ByIndex measure from the given ByPoint measure
// and wrapping the provided points.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
func Indexed(m ByPoint, points []Point) ByIndex {
	return measure.Indexed(m, points)
}

// DependentIndexed is a measure that uses a custom cost function to calculate
// parameter dependent costs for connecting two points by index. If the measures
// are time dependent all future stops in the sequence will be fully
// recalculated. Otherwise there will be a constant shift to achieve better
// performance.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
func DependentIndexed(
	timeDependent bool,
	cost func(
		from,
		to int,
		data *measure.VehicleData,
	) float64,
) DependentByIndex {
	return measure.DependentIndexed(timeDependent, cost)
}

// Scale the cost of some other measure by a constant.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
func Scale(m ByIndex, constant float64) ByIndex {
	return measure.Scale(m, constant)
}

// ScaleByPoint scales the cost of some other measure by a constant.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
func ScaleByPoint(m ByPoint, constant float64) ByPoint {
	return measure.ScaleByPoint(m, constant)
}

// ByClockwise implements sort.Interface for sorting points clockwise around a
// central point.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
func ByClockwise(center Point, points []Point) sort.Interface {
	return measure.ByClockwise(center, points)
}

// LessClockwise returns true if a is closer to a central point than b, and
// false if it is not.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
func LessClockwise(center, a, b Point) bool {
	return measure.LessClockwise(center, a, b)
}

// Sparse measure returns pre-computed costs between two locations without
// requiring a full data set. If two locations do not have an associated cost,
// then a backup measure is used.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
func Sparse(m ByIndex, arcs map[int]map[int]float64) ByIndex {
	return measure.Sparse(m, arcs)
}

// Sum adds other measures together.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
func Sum(m ...measure.ByIndex) ByIndex {
	return measure.Sum(m...)
}

// TaxicabByPoint adds absolute distances between two points in all dimensions.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
func TaxicabByPoint() ByPoint {
	return measure.TaxicabByPoint()
}

// Truncate the cost of some other measure.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
func Truncate(m ByIndex, lower, upper float64) ByIndex {
	return measure.Truncate(m, lower, upper)
}

// Unique returns a ByIndex that uses a reference slice to map the indices of a
// point to the index of the measure.
// m represents a matrix of unique points.
// references maps a stop (by index) to an index in m.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
func Unique(m ByIndex, references []int) ByIndex {
	return measure.Unique(m, references)
}

// Location measure returns the sum of the cost computed by the passed in
// measure and the specified cost of the 'to' location. This cost is read from
// the passed in costs slice.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
func Location(
	m ByIndex,
	costs []float64,
	durationGroups DurationGroups,
) (ByIndex, error) {
	return measure.Location(m, costs, durationGroups)
}

// Matrix measure returns pre-computed cost between two locations. Cost is
// assumed to be asymmetric.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
func Matrix(arcs [][]float64) ByIndex {
	return measure.Matrix(arcs)
}

// IsTriangular returns true if the triangle inequality holds for the provided
// measure. It returns false if the measure does not implement the Triangular
// interface or the triangle inequality does not hold.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
func IsTriangular(m any) bool {
	return measure.IsTriangular(m)
}

// BuildMatrixRequestPoints builds a slice of points in the correct format to
// request a matrix from any of the supported platforms (e.g. OSRM, Routingkit,
// Google, HERE). It takes the stops to be routed, start and end stops of
// vehicles (optional) and the number of to be used.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
func BuildMatrixRequestPoints(
	stops, starts,
	ends []Point,
	vehiclesCount int,
) ([]Point, error) {
	return measure.BuildMatrixRequestPoints(stops, starts, ends, vehiclesCount)
}

// OverrideZeroPoints overrides points that have been passed as placeholders
// [0,0] to build the matrix with zero values.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
func OverrideZeroPoints(
	points []Point,
	m ByIndex,
) ByIndex {
	return measure.OverrideZeroPoints(points, m)
}

// ByPointLoader can be embedded in schema structs and unmarshals a ByPoint JSON
// object into the appropriate implementation.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
type ByPointLoader measure.ByPointLoader

// ByIndexLoader can be embedded in schema structs and unmarshals a ByIndex JSON
// object into the appropriate implementation.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/nextroute].
type ByIndexLoader measure.ByIndexLoader
