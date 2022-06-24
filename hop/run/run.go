// Package run provides tools for running solvers.
package run

import (
	"github.com/nextmv-io/sdk/hop/store/types"
)

/*
Run a solver via a handler.

	func main() {
		handler := func(v int, opt types.Options) (types.Solver, error) {
			s := store.New()
			x := store.Var(s, v)
			s = s.Value(...).Format(...).Generate(...) // Modify the Store.

			return s.Maximizer(opt), nil
			// return s.Minimizer(opt), nil
			// return s.Satisfier(opt), nil
		}
		run.Run(handler)
	}
*/
func Run[T any](handler func(T, types.Options) (types.Solver, error)) {
	connect()

	runFunc(handler)
}

var runFunc func(any)
