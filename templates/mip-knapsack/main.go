// package main holds the implementation of the mip-knapsack template.
package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/nextmv-io/sdk/mip"
	"github.com/nextmv-io/sdk/model"
	"github.com/nextmv-io/sdk/run"
	"github.com/nextmv-io/sdk/run/schema"
)

// This template demonstrates how to solve a Mixed Integer Programming problem.
// To solve a mixed integer problem is to optimize a linear objective function
// of many variables, subject to linear constraints. We demonstrate this by
// solving the well known knapsack problem.
func main() {
	err := run.CLI(solver).Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
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

// The Option for the solver.
type Option struct {
	// A duration limit of 0 is treated as infinity. For cloud runs you need to
	// set an explicit duration limit which is why it is currently set to 10s
	// here in case no duration limit is set. For local runs there is no time
	// limitation. If you want to make cloud runs for longer than 5 minutes,
	// please contact: support@nextmv.io
	Limits struct {
		Duration time.Duration `json:"duration" default:"10s"`
	} `json:"limits"`
}

// Output is the output of the solver.
type Output struct {
	Status  string  `json:"status,omitempty"`
	Runtime string  `json:"runtime,omitempty"`
	Items   []item  `json:"items,omitempty"`
	Value   float64 `json:"value,omitempty"`
}

func solver(_ context.Context, input input, opts Option) (schema.Output, error) {
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
		return schema.Output{}, err
	}

	// We create the solve options we will use
	solveOptions := mip.NewSolveOptions()

	// Limit the solve to a maximum duration
	if err = solveOptions.SetMaximumDuration(opts.Limits.Duration); err != nil {
		return schema.Output{}, err
	}

	// Set the relative gap to 0% (highs' default is 5%)
	if err = solveOptions.SetMIPGapRelative(0); err != nil {
		return schema.Output{}, err
	}

	// Set verbose level to see a more detailed output
	solveOptions.SetVerbosity(mip.Off)

	solution, err := solver.Solve(solveOptions)
	if err != nil {
		return schema.Output{}, err
	}

	output, err := format(solution, input, x, opts)
	if err != nil {
		return schema.Output{}, err
	}

	return output, nil
}

func format(
	solution mip.Solution,
	input input,
	x model.MultiMap[mip.Bool, item],
	opts Option,
) (output schema.Output, err error) {
	o := Output{}
	o.Status = "infeasible"
	o.Runtime = solution.RunTime().String()

	if solution != nil && solution.HasValues() {
		if solution.IsOptimal() {
			o.Status = "optimal"
		} else {
			o.Status = "suboptimal"
		}

		o.Value = solution.ObjectiveValue()

		items := make([]item, 0)

		for _, item := range input.Items {
			// if the value of x for an item is 1 it means it is in the
			// knapsack
			if solution.Value(x.Get(item)) == 1 {
				items = append(items, item)
			}
		}
		o.Items = items
	} else {
		return schema.Output{}, errors.New("no solution found")
	}

	output = schema.NewOutput(opts, o)
	return output, nil
}
