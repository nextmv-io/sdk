package alns

import "time"

// SolveInformation defines the information available to the solver.
type SolveInformation[T baseSolution[T]] interface {
	// DeltaScore returns the delta score of the last executed solve operator.
	DeltaScore() float64

	// Iteration returns the current iteration.
	Iteration() int

	// Solver returns the solver.
	Solver() Solver[T]

	// SolveOperators returns the solve-operators that has been executed in
	// the current iteration.
	SolveOperators() SolveOperators[T]

	// Start returns the start time of the solver.
	Start() time.Time
}

type solveInformationImpl[T baseSolution[T]] struct {
	deltaScore     float64
	iteration      int
	solver         Solver[T]
	solveOperators SolveOperators[T]
	start          time.Time
}

func (s *solveInformationImpl[T]) Iteration() int {
	return s.iteration
}

func (s *solveInformationImpl[T]) Solver() Solver[T] {
	return s.solver
}

func (s *solveInformationImpl[T]) SolveOperators() SolveOperators[T] {
	return s.solveOperators
}

func (s *solveInformationImpl[T]) Start() time.Time {
	return s.start
}

func (s *solveInformationImpl[T]) DeltaScore() float64 {
	return s.deltaScore
}
