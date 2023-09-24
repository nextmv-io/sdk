package mip

// Solver for a MIP problem.
type Solver interface {
	// Solve is the entrypoint to solve the model associated with
	// the invoking solver. Returns a solution when the invoking solver
	// reaches a conclusion.
	Solve(options SolveOptions) (Solution, error)
}

// SolverProvider identifier for a back-end solver.
type SolverProvider string

const (
	// Highs is the [HiGHS] open-source solver.
	//
	// [HiGHS]: https://highs.dev
	Highs SolverProvider = "highs"
	// Xpress is the [FICO Xpress] commercial solver.
	//
	// [FICO Xpress]: https://www.fico.com/en/products/fico-xpress-optimization
	Xpress SolverProvider = "xpress"
)
