// Package run provides tools for running solvers.
package run

import (
	"github.com/nextmv-io/sdk/hop/store/types"
)

// Run a solver via a handler.
func Run[T any](handler func(T, types.Options) (types.Solver, error)) {
	connect()

	runFunc(handler)
}

var runFunc func(any)
