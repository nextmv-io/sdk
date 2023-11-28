package alns

import (
	"time"
)

// SolveInformation defines the information available to the solver.
type SolveInformation[T Solution[T]] interface {
	// DeltaScore returns the delta score of the last executed solve operator.
	DeltaScore() float64

	// Iteration returns the current iteration.
	Iteration() int

	// Solver returns the solver.
	Solver() Solver[T, SolveOptions]

	// SolveOperators returns the solve-operators that has been executed in
	// the current iteration.
	SolveOperators() SolveOperators[T]

	// Start returns the start time of the solver.
	Start() time.Time
}
