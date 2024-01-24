package nextroute

import "time"

// SolveInformation contains information about the current solve.
type SolveInformation interface {
	// DeltaScore returns the delta score of the last executed solve operator.
	DeltaScore() float64

	// Iteration returns the current iteration.
	Iteration() int

	// Solver returns the solver.
	Solver() Solver
	// SolveOperators returns the solve-operators that has been executed in
	// the current iteration.
	SolveOperators() SolveOperators
	// Start returns the start time of the solver.
	Start() time.Time
}
