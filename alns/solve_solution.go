// Package alns defines interfaces for Adaptive Large Neighborhood Search.
package alns

// BaseSolution is the base interface that a solution must implement.
type BaseSolution[T any] interface {
	// Copy returns a copy of the solution which must be of type T the
	// `derived` type. This copy must be a deep copy.
	Copy() T

	// Score returns the score of the solution.
	Score() float64
}

// Solution is a solution to a problem. It defines the minimum interface
// that a solution must implement to be used in Adaptive Local Neighborhood
// Search.
type Solution[S BaseSolution[S]] interface {
	BaseSolution[S]
}
