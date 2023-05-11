// package main holds the implementation of the mip-meal-allocation template.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/nextmv-io/sdk/mip"
	"github.com/nextmv-io/sdk/run"
	"github.com/nextmv-io/sdk/run/schema"
)

// This template demonstrates how to solve a Mixed Integer Programming problem.
// To solve a mixed integer problem is to optimize a linear objective function
// of many variables, subject to linear constraints. We demonstrate this by
// solving a made up problem we named MIP meal allocation.
//
// MIP meal allocation is a demo program in which we maximize the number of
// binkies our bunnies will execute by selecting their meals.
//
// A binky is when a bunny jumps straight up and quickly twists its hind end,
// head, or both. A bunny may binky because it is feeling happy or safe in its
// environment.
func main() {
	err := run.CLI(solver).Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

// The input defines a number of meals we can use to maximize binkies. Each
// meal consists out of one or more items and in total we can only use the
// number of items we have in stock.
type input struct {
	Items []item `json:"items"`
	Meals []meal `json:"meals"`
}

type item struct {
	Name  string `json:"name"`
	Stock int    `json:"stock"`
}

type ingredient struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type meal struct {
	Name        string       `json:"name"`
	Ingredients []ingredient `json:"ingredients"`
	Binkies     int          `json:"binkies"`
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
	Status  string         `json:"status,omitempty"`
	Runtime string         `json:"runtime,omitempty"`
	Binkies float64        `json:"binkies,omitempty"`
	Meals   []MealQuantity `json:"meals,omitempty"`
}

// MealQuantity is the number of meals of a specific type.
type MealQuantity struct {
	Name     string `json:"name,omitempty"`
	Quantity int    `json:"quantity,omitempty"`
}

func solver(input input, opts Option) (schema.Output, error) {
	// We start by creating a MIP model.
	m := mip.NewModel()

	// We want to maximize the number of binkies.
	m.Objective().SetMaximize()

	// Map the name of the item to a constraint so we can retrieve it to
	// add a term for the consumption by each meal
	itemInStockConstraints := make(map[string]mip.Constraint)

	for _, item := range input.Items {
		// We create a constraint for each item which will constrain the
		// number of items we use to what we have in stock.
		itemInStockConstraint := m.NewConstraint(
			mip.LessThanOrEqual,
			float64(item.Stock),
		)

		itemInStockConstraints[item.Name] = itemInStockConstraint
	}

	// Map the name of the meal to the variable so we can retrieve it for
	// reporting purposes
	nrMealsVars := make(map[string]mip.Int)

	for _, meal := range input.Meals {
		// We create an integer variable for each meal representing how
		// many instances of meal we will serve
		nrMealsVar := m.NewInt(
			0,
			math.MaxInt64,
		)

		// We add the number of binkies we generate by serving nrMealsVar meals
		m.Objective().NewTerm(float64(meal.Binkies), nrMealsVar)

		nrMealsVars[meal.Name] = nrMealsVar

		for _, ingredient := range meal.Ingredients {
			if _, present := itemInStockConstraints[ingredient.Name]; !present {
				return schema.Output{},
					fmt.Errorf("meal %v, uses undefined item %v",
						meal.Name,
						ingredient.Name,
					)
			}

			// We add the number of items we consume of item in the ingredient
			// to the constraint that limits how many we can serve
			itemInStockConstraints[ingredient.Name].NewTerm(
				float64(ingredient.Quantity),
				nrMealsVar,
			)
		}
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

	solution, err := solver.Solve(solveOptions)
	if err != nil {
		return schema.Output{}, err
	}

	output, err := format(solution, nrMealsVars, opts)
	if err != nil {
		return schema.Output{}, err
	}

	return output, nil
}

func format(
	solution mip.Solution,
	nrMealsVars map[string]mip.Int,
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
		o.Binkies = solution.ObjectiveValue()
		meals := make([]MealQuantity, 0)
		for name, nrMealsVar := range nrMealsVars {
			meals = append(meals, MealQuantity{
				Name:     name,
				Quantity: int(math.Round(solution.Value(nrMealsVar))),
			})
		}
		o.Meals = meals
	} else {
		return schema.Output{}, errors.New("no solution found")
	}
	output = schema.NewOutput(o, opts)
	return output, nil
}
