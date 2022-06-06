package run

import (
	"github.com/nextmv-io/sdk/hop/solve"
)

// Run a solver by a handler.
func Run[T any](handler func(T, solve.Options) (solve.Solver, error)) {
	connect()

	runFunc(handler)
}

var runFunc func(any)
