package measure

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

// IsTriangular returns true if the triangle inequality holds for the provided
// measure.
func IsTriangular(m any) bool {
	if t, ok := m.(Triangular); ok {
		return t.Triangular()
	}
	return false
}
