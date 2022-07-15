package route

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

// Override measure uses a default measure for all arcs that are not true for a
// condition. It uses an override measure for all arcs that are true for the
// condition.
func Override(
	defaultByIndex ByIndex,
	overrideByIndex ByIndex,
	condition func(from, to int) bool,
) ByIndex {
	connect()
	return overrideFunc(defaultByIndex, overrideByIndex, condition)
}

// HaversineByPoint estimates meters connecting two points along the surface
// of the earth.
func HaversineByPoint() ByPoint {
	connect()
	return haversineByPointFunc()
}

// ConstantByPoint measure always estimates the same cost.
func ConstantByPoint(c float64) ByPoint {
	connect()
	return constantByPointFunc(c)
}

// Constant measure always estimates the same cost.
func Constant(c float64) ByIndex {
	connect()
	return constantFunc(c)
}

// Indexed creates a ByIndex measure from the given ByPoint measure
// and wrapping the provided points.
func Indexed(m ByPoint, points []Point) ByIndex {
	connect()
	return indexedFunc(m, points)
}

// Scale the cost of some other measure by a constant.
func Scale(m ByIndex, constant float64) ByIndex {
	connect()
	return scaleFunc(m, constant)
}

var (
	overrideFunc         func(ByIndex, ByIndex, func(int, int) bool) ByIndex
	haversineByPointFunc func() ByPoint
	constantByPointFunc  func(float64) ByPoint
	constantFunc         func(float64) ByIndex
	indexedFunc          func(ByPoint, []Point) ByIndex
	scaleFunc            func(ByIndex, float64) ByIndex
)
