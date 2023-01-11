package kmeans_test

import (
	"fmt"

	"github.com/nextmv-io/sdk/cluster/kmeans"
	"github.com/nextmv-io/sdk/measure"
)

func ExampleSolver() {
	points := []measure.Point{
		{2.5, 2.5},
		{7.5, 7.5},
		{5.0, 7.5},
	}
	// Create a model.
	model, err := kmeans.NewModel(points, 2)
	if err != nil {
		panic(err)
	}

	// Set maximum points in first cluster to one to make the
	// solution predictable.
	model.ClusterModels()[0].SetMaximumPoints(1)

	// Create a solver using the model.
	solver, err := kmeans.NewSolver(model)
	if err != nil {
		panic(err)
	}

	// Create solve options to configure the solver.
	solveOptions := kmeans.NewSolveOptions()

	// Solve the model using the solve options.
	solution, err := solver.Solve(solveOptions)
	if err != nil {
		panic(err)
	}

	// Print the number of clusters in the solution.
	fmt.Println(len(solution.Clusters()))
	// Print the number of unassigned points in the solution.
	fmt.Println(len(solution.Unassigned()))
	// Print the number of points in the first cluster.
	fmt.Println(len(solution.Clusters()[0].Points()))
	// Print the number of points in the second cluster.
	fmt.Println(len(solution.Clusters()[1].Points()))
	// Print the centroid of the second cluster.

	// Output:
	// 2
	// 0
	// 1
	// 2
}
