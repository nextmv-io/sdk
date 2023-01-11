package route

import (
	"errors"
	"sort"

	"github.com/nextmv-io/sdk/measure"
)

// Point represents a point in space. It may have any dimension.
type Point = measure.Point

// ByIndex estimates the cost of going from one index to another.
type ByIndex = measure.ByIndex

// ByPoint estimates the cost of going from one point to another.
type ByPoint = measure.ByPoint

// Triangular indicates that the triangle inequality holds for every
// measure that implements it.
type Triangular = measure.Triangular

// DurationGroups represents a slice of duration groups. Each duration group is
// used to account for additional service costs whenever a stop of a group is
// approached first.
type DurationGroups = measure.DurationGroups

// DurationGroup groups stops by index which have additional service costs
// attached to them.
type DurationGroup = measure.DurationGroup

// Bin is a measure that selects from a slice of indexed measures. Logic
// defined in the selector function determines which measure is used in the
// cost calculation.
func Bin(
	measures []measure.ByIndex,
	selector func(from, to int) int,
) ByIndex {
	return measure.Bin(measures, selector)
}

// Override measure uses a default measure for all arcs that are not true for a
// condition. It uses an override measure for all arcs that are true for the
// condition.
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
func DebugOverride(
	defaultByIndex ByIndex,
	overrideByIndex ByIndex,
	condition func(from, to int) bool,
) ByIndex {
	return measure.DebugOverride(defaultByIndex, overrideByIndex, condition)
}

// Power raises the cost of some other measure to an exponent.
func Power(m ByIndex, exponent float64) ByIndex {
	return measure.Power(m, exponent)
}

// HaversineByPoint estimates meters connecting two points along the surface
// of the earth.
func HaversineByPoint() ByPoint {
	return measure.HaversineByPoint()
}

// ConstantByPoint measure always estimates the same cost.
func ConstantByPoint(c float64) ByPoint {
	return measure.ConstantByPoint(c)
}

// Constant measure always estimates the same cost.
func Constant(c float64) ByIndex {
	return measure.Constant(c)
}

// EuclideanByPoint computes straight line distance connecting two indices.
func EuclideanByPoint() ByPoint {
	return measure.EuclideanByPoint()
}

// Indexed creates a ByIndex measure from the given ByPoint measure
// and wrapping the provided points.
func Indexed(m ByPoint, points []Point) ByIndex {
	return measure.Indexed(m, points)
}

// Scale the cost of some other measure by a constant.
func Scale(m ByIndex, constant float64) ByIndex {
	return measure.Scale(m, constant)
}

// ScaleByPoint scales the cost of some other measure by a constant.
func ScaleByPoint(m ByPoint, constant float64) ByPoint {
	return measure.ScaleByPoint(m, constant)
}

// ByClockwise implements sort.Interface for sorting points clockwise around a
// central point.
func ByClockwise(center Point, points []Point) sort.Interface {
	return measure.ByClockwise(center, points)
}

// LessClockwise returns true if a is closer to a central point than b, and
// false if it is not.
func LessClockwise(center, a, b Point) bool {
	return measure.LessClockwise(center, a, b)
}

// Sparse measure returns pre-computed costs between two locations without
// requiring a full data set. If two locations do not have an associated cost,
// then a backup measure is used.
func Sparse(m ByIndex, arcs map[int]map[int]float64) ByIndex {
	return measure.Sparse(m, arcs)
}

// Sum adds other measures together.
func Sum(m ...measure.ByIndex) ByIndex {
	return measure.Sum(m...)
}

// TaxicabByPoint adds absolute distances between two points in all dimensions.
func TaxicabByPoint() ByPoint {
	return measure.TaxicabByPoint()
}

// Truncate the cost of some other measure.
func Truncate(m ByIndex, lower, upper float64) ByIndex {
	return measure.Truncate(m, lower, upper)
}

// Location measure returns the sum of the cost computed by the passed in
// measure and the specified cost of the 'to' location. This cost is read from
// the passed in costs slice.
func Location(
	m ByIndex,
	costs []float64,
	durationGroups DurationGroups,
) (ByIndex, error) {
	return measure.Location(m, costs, durationGroups)
}

// Matrix measure returns pre-computed cost between two locations. Cost is
// assumed to be asymmetric.
func Matrix(arcs [][]float64) ByIndex {
	return measure.Matrix(arcs)
}

// IsTriangular returns true if the triangle inequality holds for the provided
// measure. It returns false if the measure does not implement the Triangular
// interface or the triangle inequality does not hold.
func IsTriangular(m any) bool {
	return measure.IsTriangular(m)
}

// BuildMatrixRequestPoints builds a slice of points in the correct format to
// request a matrix from any of the supported platforms (e.g. OSRM, Routingkit,
// Google, HERE). It takes the stops to be routed, start and end stops of
// vehicles (optional) and the number of to be used.
func BuildMatrixRequestPoints(
	stops, starts,
	ends []Point,
	vehiclesCount int,
) ([]Point, error) {
	if len(starts) > 0 && len(starts) != vehiclesCount {
		return nil, errors.New(
			"if starts are given, they must match the number of vehicles",
		)
	}
	if len(ends) > 0 && len(ends) != vehiclesCount {
		return nil, errors.New(
			"if ends are given, they must match the number of vehicles",
		)
	}
	count := len(stops)
	// Create points array of the expected size
	points := make([]Point, count+2*vehiclesCount)
	for i := range points {
		// Set default values
		points[i] = Point{0, 0}
	}
	copy(points, stops)

	if len(starts) > 0 {
		for v, start := range starts {
			points[count+v*2] = start
		}
	}

	if len(ends) > 0 {
		for v, end := range ends {
			points[count+v*2+1] = end
		}
	}
	return points, nil
}

// OverrideZeroPoints overrides points that have been passed as placeholders
// [0,0] to build the matrix with zero values.
func OverrideZeroPoints(
	points []Point,
	m ByIndex,
) ByIndex {
	m = Override(
		m,
		Constant(0),
		func(from, to int) bool {
			fromOverride := points[from][0] == 0 && points[from][1] == 0
			toOverride := points[to][0] == 0 && points[to][1] == 0
			if fromOverride || toOverride {
				return true
			}
			return false
		},
	)

	return m
}
