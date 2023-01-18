// package main holds the implementation of the cluster template.
package main

import (
	"context"
	"log"
	"time"

	"github.com/nextmv-io/sdk/cluster/kmeans"
	"github.com/nextmv-io/sdk/measure"
	"github.com/nextmv-io/sdk/run"
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

// Output is the output of the solver.
type Output struct {
	Clusters          []cluster       `json:"clusters"`
	Feasible          bool            `json:"feasible"`
	Unassigned        []measure.Point `json:"unassigned"`
	UnassignedIndices []int           `json:"unassigned_indices"`
}

// ClusterOptions holds the options for the solver.
type ClusterOptions struct {
	Limits struct {
		Duration time.Duration `json:"duration" default:"10s"`
	} `json:"limits"`
}

func main() {
	err := run.CLI(solver).Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func solver(input input, opts ClusterOptions) ([]Output, error) {
	// We create a new model with the given points and number of clusters.
	// We also pass the options to the model to set the maximum weight and
	// maximum number of points per cluster.
	maximumPoints := make([]int, input.Clusters)
	maximumValues := make([]int, input.Clusters)
	values := make([][]int, input.Clusters)

	for idx := 0; idx < input.Clusters; idx++ {
		maximumPoints[idx] = input.MaximumPoints
		maximumValues[idx] = input.MaximumWeight
		values[idx] = input.Weights
	}

	// We create a kmeans model using options.
	model, err := kmeans.NewModel(
		input.Points,
		input.Clusters,
		kmeans.MaximumPoints(maximumPoints),
		kmeans.MaximumSumValue(maximumValues, values),
	)
	if err != nil {
		return nil, err
	}

	// We create a solver with the model.
	solver, err := kmeans.NewSolver(model)
	if err != nil {
		return nil, err
	}

	// We create the solve options we will use and set the time limit
	// and the measure to use.
	solveOptions := kmeans.NewSolveOptions().
		SetMeasure(measure.EuclideanByPoint()).
		SetMaximumDuration(opts.Limits.Duration)

	solution, err := solver.Solve(solveOptions)
	if err != nil {
		panic(err)
	}

	return []Output{format(solution)}, nil
}

// format returns a function to format the solution output.
func format(solution kmeans.Solution) Output {
	o := Output{
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
