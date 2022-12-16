// package main holds the implementation of the mip-knapsack template.
package main

import (
	"time"

	"github.com/nextmv-io/sdk/mip"
	"github.com/nextmv-io/sdk/model"
	"github.com/nextmv-io/sdk/run"
	"github.com/nextmv-io/sdk/store"
)

// This template demonstrates how to solve a Mixed Integer Programming problem.
// To solve a mixed integer problem is to optimize a linear objective function
// of many variables, subject to linear constraints. We demonstrate this by
// solving the well known knapsack problem.
func main() {
	run.Run(solver)
}

// An item has a Value, Weight and Volume. ItemID is optional and can be any
// type.
type item struct {
	ItemID string  `json:"item_id,omitempty"`
	Value  float64 `json:"value"`
	Weight float64 `json:"weight"`
}

// ID is implemented to fulfill the model.Identifier interface.
func (i item) ID() string {
	return i.ItemID
}

// A knapsack holds the most valuable set of items possible while not exceeding
// its carrying capacity.
type input struct {
	Items          []item `json:"items"`
	WeightCapacity int    `json:"weight_capacity"`
}

func solver(input input, opts store.Options) (store.Solver, error) {
	// We start by creating a MIP model.
	m := mip.NewModel()

	// x is a multimap representing a set of variables. It is initialized with a
	// create function and, in this case one set of elements. The elements can
	// be used as an index to the multimap. To retrieve a variable, call
	// x.Get(element) where element is an element from the index set
	// (input.Items).
	x := model.NewMultiMap(
		func(...item) mip.Bool {
			return m.NewBool()
		}, input.Items)

	// We want to maximize the value of the knapsack.
	m.Objective().SetMaximize()

	// This constraint ensures the weightCapacity of the knapsack will not be
	// exceeded.
	weightCapacity := m.NewConstraint(
		mip.LessThanOrEqual,
		float64(input.WeightCapacity),
	)

	for _, item := range input.Items {
		m.Objective().NewTerm(item.Value, x.Get(item))
		weightCapacity.NewTerm(item.Weight, x.Get(item))
	}

	// We create a solver using the 'highs' provider
	solver, err := mip.NewSolver("highs", m)
	if err != nil {
		return nil, err
	}

	// We create the solve options we will use
	solveOptions := mip.NewSolveOptions()

	// Limit the solve to a maximum duration
	if err := solveOptions.SetMaximumDuration(opts.Limits.Duration); err != nil {
		return nil, err
	}

	// Set the relative gap to 0% (highs' default is 5%)
	if err := solveOptions.SetMIPGapRelative(0); err != nil {
		return nil, err
	}

	// Set verbose level to see a more detailed output
	solveOptions.SetVerbosity(mip.Off)

	// We use a store, and it's corresponding Format to report the solution
	// Doing this allows us to use the CLI runner in the main function.
	root := store.New()

	// Add initial solution as nil
	so := store.NewVar[mip.Solution](root, nil)

	i := 0
	root = root.Generate(func(s store.Store) store.Generator {
		return store.Lazy(func() bool {
			// only run one state transition in which we solve the mip model
			return i == 0
		}, func() store.Store {
			i++
			// Invoke the solver
			solution, err := solver.Solve(solveOptions)
			if err != nil {
				panic(err)
			}
			return s.Apply(so.Set(solution))
		})
	}).Validate(func(s store.Store) bool {
		solution := so.Get(s)
		if solution == nil {
			return false
		}
		// if the solution has values, accept it as being valid, optionally
		// write a check to test for actual validity
		b := solution.HasValues()
		return b
	}).Format(format(so, x, input))
	// A duration limit of 0 is treated as infinity. For cloud runs you need to
	// set an explicit duration limit which is why it is currently set to 10s
	// here in case no duration limit is set. For local runs there is no time
	// limitation. If you want to make cloud runs for longer than 5 minutes,
	// please contact: support@nextmv.io
	if opts.Limits.Duration == 0 {
		opts.Limits.Duration = 10 * time.Second
	}

	// We invoke Satisfier which will result in invoking Format and
	// report the solution
	return root.Satisfier(opts), nil
}

// format returns a function to format the solution output.
func format(
	so store.Var[mip.Solution],
	x model.MultiMap[mip.Bool, item],
	input input,
) func(s store.Store) any {
	return func(s store.Store) any {
		// get solution from store
		solution := so.Get(s)

		report := make(map[string]any)

		report["status"] = "infeasible"
		report["runtime"] = solution.RunTime().String()

		if solution.HasValues() {
			if solution.IsOptimal() {
				report["status"] = "optimal"
			} else {
				report["status"] = "suboptimal"
			}

			report["value"] = solution.ObjectiveValue()

			items := make([]item, 0)

			for _, item := range input.Items {
				// if the value of x for an item is 1 it means it is in the
				// knapsack
				if solution.Value(x.Get(item)) == 1 {
					items = append(items, item)
				}
			}
			report["items"] = items
		}
		return report
	}
}
