package route

// Override measure uses a default measure for all arcs that are not true for a
// condition. It uses an override measure for all arcs that are true for the
// condition.
func Override(
	defaultByIndex ByIndex,
	overrideByIndex ByIndex,
	condition func(from, to int) bool,
) ByIndex {
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

var (
	overrideFunc         func(ByIndex, ByIndex, func(int, int) bool) ByIndex
	haversineByPointFunc func() ByPoint
	constantByPointFunc  func(float64) ByPoint
	constantFunc         func(float64) ByIndex
	indexedFunc          func(ByPoint, []Point) ByIndex
)
