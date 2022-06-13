// Package model provides modeling tools.
package model

// Bounds on an objective value at some node in the search tree consist of a
// lower value and an upper value. If the lower and upper value are the same,
// the bounds have converged.
type Bounds struct {
	Lower int
	Upper int
}
