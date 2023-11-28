// Package alns defines interfaces for Adaptive Large Neighborhood Search.
package alns

import (
	"context"
	"math/rand"
)

// Solver is the interface for the Adaptive Local Neighborhood Search algorithm.
type Solver[T Solution[T], Options any] interface {
	BaseSolver[T, Options]

	// WorkSolution returns the current work solution.
	WorkSolution() T

	// BestSolution returns the best solution found so far.
	BestSolution() T

	// Random returns the random number generator used by the solver.
	Random() *rand.Rand

	// Reset will reset the solver to use solution as work solution.
	Reset(solution Solution[T])

	// SolveOperators returns the solve-operators used by the solver.
	SolveOperators() SolveOperators[T]

	// AddSolveOperators adds a number of solve-operators to the solver.
	AddSolveOperators(...SolveOperator[T])

	// SolveEvents returns the solve-events used by the solver.
	SolveEvents() SolveEvents[T]
}

// Solution is a solution to a problem. It defines the minimum interface
// that a solution must implement to be used in Adaptive Local Neighborhood
// Search.
type Solution[T any] interface {
	// Copy returns a copy of the solution which must be of type T the
	// `derived` type. This copy must be a deep copy.
	Copy() T

	// Score returns the score of the solution.
	Score() float64

	// Random returns the random number generator used by the solution.
	Random() *rand.Rand
}

// Solutions is a channel of solutions.
type Solutions[T any] <-chan T

// All returns all solutions in the channel.
func (solutions Solutions[T]) All() []T {
	solutionArray := make([]T, 0)
	for s := range solutions {
		solutionArray = append(solutionArray, s)
	}
	return solutionArray
}

// Last returns the last solution in the channel.
func (solutions Solutions[T]) Last() T {
	var solution T
	for s := range solutions {
		solution = s
	}
	return solution
}

// BaseSolver is an interface that can be implemented by a struct to indicate that
// it can solve a problem.
type BaseSolver[T Solution[T], Options any] interface {
	// Solve starts the solving process using the given options. It returns the
	// solutions as a channel.
	Solve(context.Context, Options, ...T) (Solutions[T], error)
}
