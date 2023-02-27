package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
)

// SolveOptions are the options for the solver to be used when solving.
type SolveOptions struct{}

// Solver is the interface that a solver must implement.
type Solver interface {
	Solve(solveOptions SolveOptions) (Solution, error)
}

// SolverOptions are the options for the solver.
type SolverOptions struct{}

// NewSolver creates a new solver starting from the given solution.
func NewSolver(
	solution Solution,
	options SolverOptions,
) (Solver, error) {
	connect.Connect(con, &newSolver)
	return newSolver(solution, options)
}
