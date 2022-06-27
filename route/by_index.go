package route

// Indexed creates a ByIndex measure from the given ByPoint measure
// and wrapping the provided points.
func Indexed(m ByPoint, points []Point) ByIndex {
	return indexedFunc(m, points)
}

var indexedFunc func(ByPoint, []Point) ByIndex
