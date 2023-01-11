// package main holds the implementation of the knapsack template.
package main

import (
	"log"
	"time"

	"github.com/nextmv-io/sdk/cluster/kmeans"
	"github.com/nextmv-io/sdk/measure"
	"github.com/nextmv-io/sdk/run"
	"github.com/nextmv-io/sdk/store"
)

type input struct {
	Points        []measure.Point `json:"points"`
	Clusters      int             `json:"clusters"`
	Weights       []int           `json:"weights"`
	MaximumWeight int             `json:"maximum_weight"`
	MaximumPoints int             `json:"maximum_points"`
}

type cluster struct {
	Index    int             `json:"index"`
	Centroid measure.Point   `json:"centroid"`
	Points   []measure.Point `json:"points"`
	Indices  []int           `json:"indices"`
}
type output struct {
	Clusters          []cluster       `json:"clusters"`
	Feasible          bool            `json:"feasible"`
	Unassigned        []measure.Point `json:"unassigned"`
	UnassignedIndices []int           `json:"unassigned_indices"`
}

func main() {
	err := run.Run(solver)
	if err != nil {
		log.Fatal(err)
	}
}

func solver(input input, opts store.Options) (store.Solver, error) {
	// We start by creating a kmeans model.
	model, err := kmeans.NewModel(input.Points, input.Clusters)
	if err != nil {
		return nil, err
	}

	// We set the maximum weight and maximum points for each cluster.
	for _, c := range model.ClusterModels() {
		c.SetMaximumPoints(input.MaximumPoints)

		_, err = c.SetMaximumSumValue(input.MaximumWeight, input.Weights)
		if err != nil {
			return nil, err
		}
	}

	// We create a solver with the model.
	solver, err := kmeans.NewSolver(model)
	if err != nil {
		return nil, err
	}

	// We create the solve options we will use and set the time limit
	// and th measure to use.
	solveOptions := kmeans.NewSolveOptions().
		SetMeasure(measure.EuclideanByPoint()).
		SetMaximumDuration(opts.Limits.Duration)

	// We use a store, and it's corresponding Format to report the solution
	// Doing this allows us to use the CLI runner in the main function.
	root := store.New()

	// Add initial solution as nil
	so := store.NewVar[kmeans.Solution](root, nil)

	i := 0
	root = root.Generate(func(s store.Store) store.Generator {
		return store.Lazy(func() bool {
			// only run one state transition in which we solve the cluster model
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
		return solution != nil
	}).Format(format(so))

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
	so store.Var[kmeans.Solution],
) func(s store.Store) any {
	return func(s store.Store) any {
		// get solution from store
		solution := so.Get(s)

		o := output{
			Clusters:          make([]cluster, len(solution.Clusters())),
			Feasible:          solution.Feasible(),
			Unassigned:        solution.Unassigned(),
			UnassignedIndices: solution.UnassignedIndices(),
		}

		for idx, c := range solution.Clusters() {
			o.Clusters[idx] = cluster{
				Index:    idx,
				Centroid: c.Centroid(),
				Points:   c.Points(),
				Indices:  c.Indices(),
			}
		}

		return o
	}
}
