// package main holds the implementation of the mip incentive allocation
// template.
package main

import (
	"context"
	"log"
	"math"

	"github.com/nextmv-io/sdk/mip"
	"github.com/nextmv-io/sdk/run"
	"github.com/nextmv-io/sdk/run/schema"
)

// This template demonstrates how to solve a Mixed Integer Programming problem.
// To solve a mixed integer problem is to optimize a linear objective function
// of many variables, subject to linear constraints. We demonstrate this by
// solving a incentive allocation problem.
func main() {
	err := run.CLI(solver).Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

// The options for the solver.
type options struct {
	Solve mip.SolveOptions `json:"solve,omitempty"`
}

// Input of the problem.
type input struct {
	// Users for the problem including name and a cost/effect pair per
	// available incentive.
	Users []struct {
		ID string `json:"id"`
		// Incentives that can be applied to the user.
		Incentives []struct {
			ID string `json:"id"`
			// Positive effect of the incentive.
			Effect float64 `json:"effect"`
			// Cost of the incentive that will be subtracted from the budget.
			Cost float64 `json:"cost"`
		} `json:"incentives"`
	} `json:"users"`
	Budget int `json:"budget"`
}

// assignments is used to print the solutio∑n and represents the
// combination of a user with the assigned incentive.
type assignments struct {
	User        string  `json:"user"`
	IncentiveID string  `json:"incentive_id"`
	Cost        float64 `json:"cost"`
	Effect      float64 `json:"effect"`
}

// solution represents the decisions made by the solver.
type solution struct {
	Assignments []assignments `json:"assignments,omitempty"`
}

// solver is the entrypoint of the program where a model is defined and solved.
func solver(_ context.Context, input input, options options) (schema.Output, error) {
	// Translate the input to a MIP model.
	model, variables := model(input)

	// Create a solver using a provider. Please see the documentation on
	// [mip.SolverProvider] for more information on the available providers.
	solver, err := mip.NewSolver(mip.Highs, model)
	if err != nil {
		return schema.Output{}, err
	}

	// Solve the model and get the solution.
	solution, err := solver.Solve(options.Solve)
	if err != nil {
		return schema.Output{}, err
	}

	// Format the solution into the desired output format and add custom
	// statistics.
	output := mip.Format(options, format(input, solution, variables), solution)
	output.Statistics.Result.Custom = mip.DefaultCustomResultStatistics(model, solution)

	return output, nil
}

// model creates a MIP model from the input. It also returns the decision
// variables.
func model(input input) (mip.Model, map[string][]mip.Var) {
	// We start by creating a MIP model.
	model := mip.NewModel()

	// We want to maximize the value of the problem.
	model.Objective().SetMaximize()

	// This constraint ensures the budget of the will not be exceeded.
	budgetConstraint := model.NewConstraint(
		mip.LessThanOrEqual,
		float64(input.Budget),
	)

	// Create a map of user ID to a slice of decision variables, one for each
	// incentive.
	userIncentiveVariables := make(map[string][]mip.Var, len(input.Users))
	for _, user := range input.Users {
		// For each user, create the slice of variables based on the number of
		// incentives.
		userIncentiveVariables[user.ID] = make([]mip.Var, len(user.Incentives))

		// This constraint ensures that each user is assigned at most one
		// incentive.
		oneIncentiveConstraint := model.NewConstraint(mip.LessThanOrEqual, 1.0)
		for i, incentive := range user.Incentives {
			// For each incentive, create a binary decision variable.
			userIncentiveVariables[user.ID][i] = model.NewBool()

			// Set the term of the variable on the objective, based on the
			// effect the incentive has on the user.
			model.Objective().NewTerm(
				incentive.Effect,
				userIncentiveVariables[user.ID][i],
			)

			// Set the term of the variable on the budget constraint, based on
			// the cost of the incentive for the user.
			budgetConstraint.NewTerm(
				incentive.Cost,
				userIncentiveVariables[user.ID][i],
			)

			// Set the term of the variable on the constraint that controls the
			// number of incentives per user.
			oneIncentiveConstraint.NewTerm(1, userIncentiveVariables[user.ID][i])
		}
	}

	return model, userIncentiveVariables
}

// format the solution from the solver into the desired output format.
func format(
	input input,
	solverSolution mip.Solution,
	userIncentiveVariables map[string][]mip.Var,
) solution {
	if !solverSolution.IsOptimal() && !solverSolution.IsSubOptimal() {
		return solution{}
	}

	assigned := []assignments{}
	for i, user := range input.Users {
		for j := range user.Incentives {
			// If the variable is not assigned, skip it.
			if int(math.Round(
				solverSolution.Value(userIncentiveVariables[user.ID][j])),
			) < 1 {
				continue
			}

			assigned = append(
				assigned,
				assignments{
					Cost:        input.Users[i].Incentives[j].Cost,
					Effect:      input.Users[i].Incentives[j].Effect,
					User:        user.ID,
					IncentiveID: input.Users[i].Incentives[j].ID,
				},
			)
		}
	}

	return solution{
		Assignments: assigned,
	}
}
