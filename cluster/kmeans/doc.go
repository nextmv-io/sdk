/*
Package kmeans provides a general interface for solving kmeans clustering
problems. The base interface is the Model which is a collection of points,
cluster models and constraints. The interface Solver is constructed by
kmeans.NewSolver. The solver can be invoked using Solver.Solve and returns
a Solution.

A new Model is created:

	points := []measure.Point{
			{2.5, 2.5},
			{7.5, 7.5},
			{5.0, 7.5},
		}

	numberOfClusters := 2

	model, err := kmeans.NewModel(points, numberOfClusters)

A Solver is created and invoked to produce a Solution:

	solver, err := kmeans.NewSolver(model)

	solution, err  := solver.Solve(kmeans.NewSolveOptions())
*/
package kmeans
