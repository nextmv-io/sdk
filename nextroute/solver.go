package nextroute

import (
	"time"

	"github.com/nextmv-io/sdk/connect"
)

type SolveOptions struct {
	Iterations      int           `json:"iterations"  usage:"number of iterations"`
	MaximumDuration time.Duration `json:"maximum_duration"  usage:"maximum duration of solver in seconds"`
}

// Solver is the interface for a solver.
type Solver interface {
	// Solve solves the problem usint the solve-options.
	Solve(solveOptions SolveOptions) (Solution, error)
	// SolverOptions returns the solver-options used to create the solver. The
	// returned options are a copy of the options used to create the solver.
	// They can be used to create a new solver and changes will have no effect
	// on this invoked solver.
	SolverOptions() SolverOptions

	SetStartSolution(solution Solution)
}

type IntParameterOptions struct {
	StartValue               int  `json:"start_value"  usage:"start value"`
	DeltaAfterIterations     int  `json:"delta_after_iterations"  usage:"delta after each iterations"`
	Delta                    int  `json:"delta"  usage:"delta"`
	MinValue                 int  `json:"min_value"  usage:"min value of parameter"`
	MaxValue                 int  `json:"max_value"  usage:"max value of parameter"`
	SnapBackAfterImprovement bool `json:"snap_back_after_improvement"  usage:"snap back to start value after improvement of best solution"`
	Zigzag                   bool `json:"zigzag"  usage:"zigzag between min and max value lik a jig saw"`
}

type SolverOptions struct {
	Unplan  IntParameterOptions `json:"unplan"  usage:"unplan parameter"`
	Plan    IntParameterOptions `json:"plan"  usage:"plan parameter"`
	Restart IntParameterOptions `json:"restart"  usage:"restart parameter"`
}

// SolverFactory is the interface for a solver-factory.
type SolverFactory interface {
	// NewSolver creates a new solver.
	NewSolver(model Model) (Solver, error)
}

func NewSolverFactory() SolverFactory {
	connect.Connect(con, &newSolverFactory)
	return newSolverFactory()
}

func NewSolver(
	solution Solution,
	options SolverOptions,
) (Solver, error) {
	connect.Connect(con, &newSolver)
	return newSolver(solution, options)
}
