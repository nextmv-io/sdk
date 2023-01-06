package kmeans

import (
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/measure"
)

// NewSolveOptions returns default solver options.
func NewSolveOptions() SolveOptions {
	connect.Connect(con, &newSolveOptions)
	return newSolveOptions()
}

// NewModel creates a new Model with the given points and number of clusters.
func NewModel(
	points []measure.Point,
	clusters int,
) (Model, error) {
	connect.Connect(con, &newModel)
	return newModel(points, clusters)
}

// NewSolver returns a new Solver implemented by the given provider.
func NewSolver(model Model) (Solver, error) {
	connect.Connect(con, &newSolver)
	return newSolver(model)
}

var (
	con             = connect.NewConnector("sdk", "KMeans")
	newSolveOptions func() SolveOptions
	newSolver       func(Model) (Solver, error)
	newModel        func([]measure.Point, int) (Model, error)
)
