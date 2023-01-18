// package main holds the implementation of the mip incentive allocation
// template.
package main

import (
	"context"
	"errors"
	"log"
	"math"
	"time"

	"github.com/nextmv-io/sdk/mip"
	"github.com/nextmv-io/sdk/run"
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

// incentiveAllocationProblem describes the needed input data to run the model.
// There is a fixed budget that must not be exceeded, a fixed number of users
// that may receive incentives, and a description of the available incentives.
type incentiveAllocationProblem struct {
	Users  []user `json:"users"`
	Budget int    `json:"budget"`
}

// incentive contains a Cost, Effect, and Availability.
type incentive struct {
	ID     string `json:"id"`
	Effect int    `json:"effect"`
	Cost   int    `json:"cost"`
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
	Status      string        `json:"status,omitempty"`
	Runtime     string        `json:"runtime,omitempty"`
	Assignments []assignments `json:"assignments,omitempty"`
	Value       float64       `json:"value,omitempty"`
}

func solver(
	input incentiveAllocationProblem,
	opts Option,
) ([]Output, error) {
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
	solveOptions.SetVerbosity(mip.Off)

	solution, err := solver.Solve(solveOptions)
	if err != nil {
		return nil, err
	}

	output, err := format(solution, input, userIncentive)
	if err != nil {
		return nil, err
	}

	return []Output{output}, nil
}

func format(
	solution mip.Solution,
	input incentiveAllocationProblem,
	userIncentive map[string][]mip.Var,
) (output Output, err error) {
	output.Status = "infeasible"
	output.Runtime = solution.RunTime().String()

	if solution != nil && solution.HasValues() {
		if solution.IsOptimal() {
			output.Status = "optimal"
		} else {
			output.Status = "suboptimal"
		}

		output.Value = solution.ObjectiveValue()

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
		output.Assignments = assigned
	} else {
		return output, errors.New("no solution found")
	}

	return output, nil
}
