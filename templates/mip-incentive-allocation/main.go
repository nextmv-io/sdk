// package main holds the implementation of the mip incentive allocation
// template.
package main

import (
	"math"
	"time"

	"github.com/nextmv-io/sdk/mip"
	"github.com/nextmv-io/sdk/run"
	"github.com/nextmv-io/sdk/store"
)

// This template demonstrates how to solve a Mixed Integer Programming problem.
// To solve a mixed integer problem is to optimize a linear objective function
// of many variables, subject to linear constraints. We demonstrate this by
// solving a incentive allocation problem.
func main() {
	run.Run(solver)
}

// incentiveAllocationProblem describes the needed input data to run the model.
// There is a fixed budget that must not be exceeded, a fixed number of users
// that may receive incentives, and a description of the available incentives.
type incentiveAllocationProblem struct {
	Users  []user `json:"users"`
	Budget int    `json:"budget"`
}

// incentive contains a Cost, Effect, and Availability.
type incentive struct {
	Effect int    `json:"effect"`
	Cost   int    `json:"cost"`
	ID     string `json:"id"`
}

// user represents the user for the problem including its name and a
// cost/effect pair per available incentive.
type user struct {
	ID         string      `json:"id"`
	Incentives []incentive `json:"incentives"`
}

// assignments is used to print the solution and represents the
// combination of a user with the assigned incentive.
type assignments struct {
	User        string `json:"user"`
	IncentiveID string `json:"incentive_id"`
	Cost        int    `json:"cost"`
	Effect      int    `json:"effect"`
}

func solver(
	input incentiveAllocationProblem,
	opts store.Options,
) (store.Solver, error) {
	// We start by creating a MIP model.
	m := mip.NewModel()

	// We want to maximize the value of the problem.
	m.Objective().SetMaximize()

	// This constraint ensures the budget of the will not be exceeded.
	budgetConstraint := m.NewConstraint(
		mip.LessThanOrEqual, float64(input.Budget),
	)

	userIncentive := make(map[string][]mip.Var, len(input.Users))
	for _, user := range input.Users {
		userIncentive[user.ID] = make([]mip.Var, len(user.Incentives))
		atMostOne := m.NewConstraint(mip.LessThanOrEqual, 1.0)
		for i, incentive := range user.Incentives {
			x := m.NewBool()
			userIncentive[user.ID][i] = x
			m.Objective().NewTerm(float64(incentive.Effect), x)
			budgetConstraint.NewTerm(float64(incentive.Cost), x)
			atMostOne.NewTerm(1, x)
		}
	}

	// We create a solver using the 'highs' provider.
	solver, err := mip.NewSolver("highs", m)
	if err != nil {
		return nil, err
	}

	// We create the solve options we will use.
	solveOptions := mip.NewSolveOptions()

	// Limit the solve to a maximum duration.
	err = solveOptions.SetMaximumDuration(opts.Limits.Duration)
	if err != nil {
		return nil, err
	}

	// Set the relative gap to 0% (highs' default is 5%).
	err = solveOptions.SetMIPGapRelative(0)
	if err != nil {
		return nil, err
	}

	// Set verbose level to see a more detailed output.
	solveOptions.SetVerbosity(mip.Low)

	// We use a store, and it's corresponding Format to report the solution
	// Doing this allows us to use the CLI runner in the main function.
	root := store.New()

	// Add initial solution as nil.
	so := store.NewVar[mip.Solution](root, nil)

	i := 0
	root = root.Generate(func(s store.Store) store.Generator {
		return store.Lazy(func() bool {
			// Only run one state transition in which we solve the mip model.
			return i == 0
		}, func() store.Store {
			i++
			// Invoke the solver.
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
		// If the solution has values, accept it as being valid, optionally
		// write a check to test for actual validity.
		b := solution.HasValues()
		return b
	}).Format(format(so, userIncentive, input))

	// If the duration limit is unset, we set it to 10s. You can configure
	// longer solver run times here. For local runs there is no time limitation.
	// If you want to make cloud runs for longer than 5 minutes, please contact:
	// sales@nextmv.io
	if opts.Limits.Duration == 0 {
		opts.Limits.Duration = 10 * time.Second
	}

	// We invoke Satisfier which will result in invoking Format and
	// report the solution.
	return root.Satisfier(opts), nil
}

// format returns a function to format the solution output.
func format(
	so store.Var[mip.Solution],
	userIncentive map[string][]mip.Var,
	input incentiveAllocationProblem,
) func(s store.Store) any {
	return func(s store.Store) any {
		// Get solution from store.
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
			// We change the output so that it is easier to read.
			assigned := []assignments{}
			for i, user := range input.Users {
				for j := range user.Incentives {
					if int(math.Round(
						solution.Value(userIncentive[user.ID][j])),
					) < 1 {
						continue
					}
					oc := assignments{
						Cost:        input.Users[i].Incentives[j].Cost,
						Effect:      input.Users[i].Incentives[j].Effect,
						User:        user.ID,
						IncentiveID: input.Users[i].Incentives[j].ID,
					}
					assigned = append(assigned, oc)
				}
			}
			report["assignments"] = assigned
			report["value"] = solution.ObjectiveValue()
		}
		return report
	}
}
