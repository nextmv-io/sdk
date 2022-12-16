// package main holds the implementation of the mip-meal-allocation template.
package main

import (
	"fmt"
	"log"
	"math"

	"github.com/nextmv-io/sdk/mip"
	"github.com/nextmv-io/sdk/run"
	"github.com/nextmv-io/sdk/store"
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
	err := run.Run(solver)
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

func solver(input input, opts store.Options) (store.Solver, error) {
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
				return nil,
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
		return nil, err
	}

	// We create the solve options we will use
	solveOptions := mip.NewSolveOptions()

	// Limit the solve to a maximum duration
	if err := solveOptions.SetMaximumDuration(opts.Limits.Duration); err != nil {
		return nil, err
	}

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
	})

	root = root.Format(format(so, nrMealsVars))

	// We invoke Maximizer which will result in invoking Format and
	// report the solution
	opts.Sense = store.Maximize
	return root.Maximizer(opts), nil
}

// format returns a function to format the solution output.
func format(
	so store.Var[mip.Solution],
	nrMealsVars map[string]mip.Int,
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

			report["binkies"] = solution.ObjectiveValue()

			type meal struct {
				Name     string `json:"name"`
				Quantity int    `json:"quantity"`
			}

			meals := make([]meal, 0)

			for name, nrMealsVar := range nrMealsVars {
				meals = append(meals, meal{
					Name:     name,
					Quantity: int(math.Round(solution.Value(nrMealsVar))),
				})
			}

			report["meals"] = meals
		}

		return report
	}
}
