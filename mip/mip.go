package mip

// DefaultSolverOptions returns default solver options.
func DefaultSolverOptions() SolverOptions {
	connect()
	return defaultSolverOptions()
}

// NewDefinition creates an empty MIP definition.
func NewDefinition() Definition {
	connect()
	return newDefinition()
}

// NewSolver returns a new Solver which will use a solver
// implemented by provider.
func NewSolver(provider SolverProvider, definition Definition) (Solver, error) {
	connect()
	return newSolver(provider, definition)
}

var (
	defaultSolverOptions func() SolverOptions
	newSolver            func(SolverProvider, Definition) (Solver, error)
	newDefinition        func() Definition
)
