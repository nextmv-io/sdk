package mip

// NewSolveOptions returns default solver options.
func NewSolveOptions() SolveOptions {
	connect()
	return newSolveOptions()
}

// NewModel creates an empty MIP model.
func NewModel() Model {
	connect()
	return newModel()
}

// NewSolver returns a new Solver which will use a solver
// implemented by provider.
func NewSolver(provider SolverProvider, model Model) (Solver, error) {
	connect()
	return newSolver(provider, model)
}

var (
	newSolveOptions func() SolveOptions
	newSolver       func(SolverProvider, Model) (Solver, error)
	newModel        func() Model
)
