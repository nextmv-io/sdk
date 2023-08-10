package mip

import (
	"github.com/nextmv-io/sdk/connect"
)

// NewSolveOptions returns default solver options.
func NewSolveOptions() SolveOptions {
	connect.Connect(con, &newSolveOptions)
	return newSolveOptions()
}

// NewModel creates an empty MIP model.
func NewModel() Model {
	connect.Connect(con, &newModel)
	return newModel()
}

// NewSolver returns a new Solver implemented by the given provider.
func NewSolver(provider SolverProvider, model Model) (Solver, error) {
	connect.Connect(con, &newSolver)
	return newSolver(provider, model)
}

var (
	con             = connect.NewConnector("sdk", "MIP")
	newSolveOptions func() SolveOptions
	newSolver       func(SolverProvider, Model) (Solver, error)
	newModel        func() Model
)
