package route

import (
	"errors"
	"sort"

	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/model"
)

// Point represents a point in space. It may have any dimension.
type Point []float64

// ByIndex estimates the cost of going from one index to another.
type ByIndex interface {
	// Cost estimates the cost of going from one index to another.
	Cost(from, to int) float64
}

// ByPoint estimates the cost of going from one point to another.
type ByPoint interface {
	// Cost estimates the cost of going from one point to another.
	Cost(from, to Point) float64
}

// Triangular indicates that the triangle inequality holds for every
// measure that implements it.
type Triangular interface {
	Triangular() bool
}

// DurationGroups represents a slice of duration groups. Each duration group is
// used to account for additional service costs whenever a stop of a group is
// approached first.
type DurationGroups []DurationGroup

// DurationGroup groups stops by index which have additional service costs
// attached to them.
type DurationGroup struct {
	Group    model.Domain
	Duration int
}

// Bin is a measure that selects from a slice of indexed measures. Logic
// defined in the selector function determines which measure is used in the
// cost calculation.
func Bin(
	measures []ByIndex,
	selector func(from, to int) int,
) ByIndex {
	connect.Connect(con, &binFunc)
	return binFunc(measures, selector)
}

// Override measure uses a default measure for all arcs that are not true for a
// condition. It uses an override measure for all arcs that are true for the
// condition.
func Override(
	defaultByIndex ByIndex,
	overrideByIndex ByIndex,
	condition func(from, to int) bool,
) ByIndex {
	connect.Connect(con, &overrideFunc)
	return overrideFunc(defaultByIndex, overrideByIndex, condition)
}

// DebugOverride returns an Override that when marshalled will include debugging
// information describing the number of queries for default and override
// elements.
func DebugOverride(
	defaultByIndex ByIndex,
	overrideByIndex ByIndex,
	condition func(from, to int) bool,
) ByIndex {
	connect.Connect(con, &debugOverrideFunc)
	return debugOverrideFunc(defaultByIndex, overrideByIndex, condition)
}

// Power raises the cost of some other measure to an exponent.
func Power(m ByIndex, exponent float64) ByIndex {
	connect.Connect(con, &powerFunc)
	return powerFunc(m, exponent)
}

// HaversineByPoint estimates meters connecting two points along the surface
// of the earth.
func HaversineByPoint() ByPoint {
	connect.Connect(con, &haversineByPointFunc)
	return haversineByPointFunc()
}

// ConstantByPoint measure always estimates the same cost.
func ConstantByPoint(c float64) ByPoint {
	connect.Connect(con, &constantByPointFunc)
	return constantByPointFunc(c)
}

// Constant measure always estimates the same cost.
func Constant(c float64) ByIndex {
	connect.Connect(con, &constantFunc)
	return constantFunc(c)
}

// EuclideanByPoint computes straight line distance connecting two indices.
func EuclideanByPoint() ByPoint {
	connect.Connect(con, &euclideanByPointFunc)
	return euclideanByPointFunc()
}

// Indexed creates a ByIndex measure from the given ByPoint measure
// and wrapping the provided points.
func Indexed(m ByPoint, points []Point) ByIndex {
	connect.Connect(con, &indexedFunc)
	return indexedFunc(m, points)
}

// Scale the cost of some other measure by a constant.
func Scale(m ByIndex, constant float64) ByIndex {
	connect.Connect(con, &scaleFunc)
	return scaleFunc(m, constant)
}

// ScaleByPoint scales the cost of some other measure by a constant.
func ScaleByPoint(m ByPoint, constant float64) ByPoint {
	connect.Connect(con, &scaleByPointFunc)
	return scaleByPointFunc(m, constant)
}

// ByClockwise implements sort.Interface for sorting points clockwise around a
// central point.
func ByClockwise(center Point, points []Point) sort.Interface {
	connect.Connect(con, &byClockwiseFunc)
	return byClockwiseFunc(center, points)
}

// LessClockwise returns true if a is closer to a central point than b, and
// false if it is not.
func LessClockwise(center, a, b Point) bool {
	connect.Connect(con, &lessClockwiseFunc)
	return lessClockwiseFunc(center, a, b)
}

// Sparse measure returns pre-computed costs between two locations without
// requiring a full data set. If two locations do not have an associated cost,
// then a backup measure is used.
func Sparse(m ByIndex, arcs map[int]map[int]float64) ByIndex {
	connect.Connect(con, &sparseFunc)
	return sparseFunc(m, arcs)
}

// Sum adds other measures together.
func Sum(m ...ByIndex) ByIndex {
	connect.Connect(con, &sumFunc)
	return sumFunc(m...)
}

// TaxicabByPoint adds absolute distances between two points in all dimensions.
func TaxicabByPoint() ByPoint {
	connect.Connect(con, &taxicabByPointFunc)
	return taxicabByPointFunc()
}

// Truncate the cost of some other measure.
func Truncate(m ByIndex, lower, upper float64) ByIndex {
	connect.Connect(con, &truncateFunc)
	return truncateFunc(m, lower, upper)
}

// Location measure returns the sum of the cost computed by the passed in
// measure and the specified cost of the 'to' location. This cost is read from
// the passed in costs slice.
func Location(
	m ByIndex,
	costs []float64,
	durationGroups DurationGroups,
) (ByIndex, error) {
	connect.Connect(con, &locationFunc)
	return locationFunc(m, costs, durationGroups)
}

// Matrix measure returns pre-computed cost between two locations. Cost is
// assumed to be asymmetric.
func Matrix(arcs [][]float64) ByIndex {
	connect.Connect(con, &matrixFunc)
	return matrixFunc(arcs)
}

// IsTriangular returns true if the triangle inequality holds for the provided
// measure. It returns false if the measure does not implement the Triangular
// interface or the triangle inequality does not hold.
func IsTriangular(m any) bool {
	connect.Connect(con, &isTriangularFunc)
	return isTriangularFunc(m)
}

// BuildMatrixRequestPoints builds a slice of points in the correct format to request a
// matrix from any of the supported platforms (e.g. OSRM, Routingkit, Google,
// HERE). It takes the stops to be routed, start and end stops of vehicles
// (optional) and the number of to be used.
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

var (
	con                  = connect.NewConnector("sdk", "Route")
	binFunc              func([]ByIndex, func(int, int) int) ByIndex
	indexedFunc          func(ByPoint, []Point) ByIndex
	constantByPointFunc  func(float64) ByPoint
	constantFunc         func(float64) ByIndex
	euclideanByPointFunc func() ByPoint
	haversineByPointFunc func() ByPoint
	locationFunc         func(ByIndex, []float64, DurationGroups) (ByIndex, error)
	matrixFunc           func([][]float64) ByIndex
	overrideFunc         func(ByIndex, ByIndex, func(int, int) bool) ByIndex
	debugOverrideFunc    func(ByIndex, ByIndex, func(int, int) bool) ByIndex
	powerFunc            func(ByIndex, float64) ByIndex
	scaleFunc            func(ByIndex, float64) ByIndex
	scaleByPointFunc     func(ByPoint, float64) ByPoint
	byClockwiseFunc      func(Point, []Point) sort.Interface
	lessClockwiseFunc    func(Point, Point, Point) bool
	sparseFunc           func(ByIndex, map[int]map[int]float64) ByIndex
	sumFunc              func(...ByIndex) ByIndex
	taxicabByPointFunc   func() ByPoint
	truncateFunc         func(ByIndex, float64, float64) ByIndex
	isTriangularFunc     func(m any) bool
)
