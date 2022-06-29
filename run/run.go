// Package run provides tools for running solvers.
package run

import (
	"github.com/nextmv-io/sdk/store"
)

/*
Run a solver via a handler.

	func main() {
		handler := func(v int, opt store.Options) (store.Solver, error) {
			s := store.New()
			x := store.NewVar(s, v)
			s = s.Value(...).Format(...).Generate(...) // Modify the Store.

			return s.Maximizer(opt), nil
			// return s.Minimizer(opt), nil
			// return s.Satisfier(opt), nil
		}
		run.Run(handler)
	}
*/
func Run[T any](handler func(T, store.Options) (store.Solver, error)) {
	connect()

	runFunc(handler)
}

var runFunc func(any)
