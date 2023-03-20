// Package alns defines interfaces for Adaptive Large Neighborhood Search.
package alns

import (
	"context"
	"math/rand"
)

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

type Solutions[T any] <-chan T

func (solutions Solutions[T]) Last() T {
	var solution T
	for s := range solutions {
		solution = s
	}
	return solution
}

type Solver[T Solution[T], Options any] interface {
	// Solve starts the solving process using the given options. It returns the
	// solutions as a channel.
	Solve(context.Context, Options, ...T) (Solutions[T], error)
}

type Progressioner interface {
	// Progression returns the progression of the solver.
	Progression() []ProgressionEntry
}

type ProgressionEntry struct {
	ElapsedSeconds float64 `json:"elapsed_seconds"`
	Value          float64 `json:"value"`
}
