// package main holds the implementation of the new-app template.
package main

import (
	"log"

	"github.com/nextmv-io/sdk/run"
	"github.com/nextmv-io/sdk/store"
)

// run.Run reads input data and solver options to run the solver.
func main() {
	err := run.Run(solver)
	if err != nil {
		log.Fatal(err)
	}
}

// input data is parsed to this struct before it is used in the solver func.
type input struct {
	Number int `json:"number"`
}

// solver returns a Solver for a problem.
func solver(i input, opt store.Options) (store.Solver, error) {
	// A new store is created.
	newStore := store.New()

	// Add variables that you want to track and make decisions on.
	aVariable := store.NewVar(newStore, i.Number)
	newStore = newStore.Generate(func(s store.Store) store.Generator {
		i := aVariable.Get(s)
		return store.Lazy(
			func() bool {
				// While this condition returns true new stores will be
				// generated. Make sure that you reach a point where false is
				// returned to stop generating new stores.
				return i > 0
			},
			func() store.Store {
				i--
				// Create new child stores by creating and applying changes to
				// the current store.
				changes := make([]store.Change, 0)
				changes = append(
					changes,
					aVariable.Set(aVariable.Get(s)-1),
				)
				return s.Apply(changes...)
			},
		)
	}).Value(func(s store.Store) int {
		// Write a term to express the solution value and return it.
		return aVariable.Get(s)
	}).Validate(func(s store.Store) bool {
		// Write a condition that checks whether the store is operationally
		// valid.
		return true
	}).Format(format(aVariable))
	// Return the solver using the Minimizer, Maximizer or Satisfier func.
	return newStore.Minimizer(opt), nil
}

// format returns a function to format the solution output.
func format(
	aVariable store.Var[int],
) func(s store.Store) any {
	return func(s store.Store) any {
		// Define the output that you need here. E.g., you can use a map
		// like it is shown below.
		output := map[string]any{
			"new_app": "This is a skeleton example.",
			"value":   aVariable.Get(s),
		}
		return output
	}
}
