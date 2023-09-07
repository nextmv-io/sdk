// package main holds the implementation of the mip-knapsack template.
package main

import (
	"context"
	"log"

	"github.com/nextmv-io/sdk/mip"
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

// The options for the solver.
type options struct {
	Limits mip.Limits `json:"limits,omitempty"`
}

// Input of the problem.
type input struct {
	Items          []item  `json:"items"`
	WeightCapacity float64 `json:"weight_capacity"`
}

// An item has a Value and Weight. ItemID is used to identify the item.
type item struct {
	ItemID string  `json:"item_id,omitempty"`
	Value  float64 `json:"value"`
	Weight float64 `json:"weight"`
}

// ID is implemented to fulfill the model.Identifier interface.
func (i item) ID() string {
	return i.ItemID
}

// solution represents the decisions made by the solver.
type solution struct {
	Items []item `json:"items,omitempty"`
}

func solver(_ context.Context, input input, options options) (schema.Output, error) {
	// We start by creating a MIP model.
	m := mip.NewModel()

	// Create a map of ID to decision variables for each item in the knapsack.
	itemVariables := make(map[string]mip.Bool, len(input.Items))
	for _, item := range input.Items {
		// Create a new binary decision variable for each item in the knapsack.
		itemVariables[item.ItemID] = m.NewBool()
	}

	// We want to maximize the value of the knapsack.
	m.Objective().SetMaximize()

	// This constraint ensures the weight capacity of the knapsack will not be
	// exceeded.
	capacityConstraint := m.NewConstraint(
		mip.LessThanOrEqual,
		input.WeightCapacity,
	)

	// For each item, set the term in the objective function and in the
	// constraint.
	for _, item := range input.Items {
		// Sets the value of the item in the objective function.
		m.Objective().NewTerm(item.Value, itemVariables[item.ItemID])

		// Sets the weight of the item in the constraint.
		capacityConstraint.NewTerm(item.Weight, itemVariables[item.ItemID])
	}

	// Create a solver using a provider. Please see the documentation on
	// [mip.SolverProvider] for more information on the available providers.
	solver, err := mip.NewSolver(mip.Highs, m)
	if err != nil {
		return schema.Output{}, err
	}

	// We create the solve options we will use.
	solveOptions := mip.NewSolveOptions()

	// Limit the solve to a maximum duration.
	if err = solveOptions.SetMaximumDuration(options.Limits.Duration); err != nil {
		return schema.Output{}, err
	}

	// Set the relative gap to 0% (highs' default is 5%)
	if err = solveOptions.SetMIPGapRelative(0); err != nil {
		return schema.Output{}, err
	}

	// Set verbose level to see a more detailed output
	solveOptions.SetVerbosity(mip.Off)

	// Solve the model and get the solution.
	solution, err := solver.Solve(solveOptions)
	if err != nil {
		return schema.Output{}, err
	}

	// Format the solution into the desired output format and add custom
	// statistics.
	output := mip.Format(options, format(input, solution, itemVariables), solution)
	output.Statistics.Result.Custom = mip.DefaultCustomResultStatistics(m, solution)

	return output, nil
}

// format the solution from the solver into the desired output format.
func format(input input, solverSolution mip.Solution, itemVariables map[string]mip.Bool) solution {
	if !solverSolution.IsOptimal() && !solverSolution.IsSubOptimal() {
		return solution{}
	}

	items := make([]item, 0)
	for _, item := range input.Items {
		if solverSolution.Value(itemVariables[item.ItemID]) < 1 {
			continue
		}
		items = append(items, item)
	}

	return solution{
		Items: items,
	}
}
