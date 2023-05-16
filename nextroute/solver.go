package nextroute

import (
	"time"

	"github.com/nextmv-io/sdk/alns"
	"github.com/nextmv-io/sdk/connect"
)

// NewSolverFactory creates a new SolverFactory. A SolverFactory can be used to
// create a new solver.
func NewSolverFactory() SolverFactory {
	connect.Connect(con, &newSolverFactory)
	return newSolverFactory()
}

// NewSolver creates a new solver. The solver can be used to solve a solution.
// The solution passed to the solver is the starting point of the solver. The
// solver will try to improve the solution.
func NewSolver(
	model Model,
	options SolverOptions,
) (Solver, error) {
	connect.Connect(con, &newSolver)
	return newSolver(model, options)
}

// SolveOptions are the options for the solver.
type SolveOptions struct {
	Iterations int           `json:"iterations"  usage:"maximum number of iterations, -1 assumes no limit" default:"-1"`
	Duration   time.Duration `json:"duration"  usage:"maximum duration of solver in seconds" default:"30s"`
}

// Solver is the interface for a solver.
type Solver interface {
	alns.Progressioner
	alns.Solver[Solution, SolveOptions]
	// SolverOptions returns the solver-options used to create the solver. The
	// returned options are a copy of the options used to create the solver.
	// They can be used to create a new solver and changes will have no effect
	// on this invoked solver.
	SolverOptions() SolverOptions
}

// IntParameterOptions are the options for an integer parameter.
type IntParameterOptions struct {
	StartValue               int  `json:"start_value"  usage:"start value"`
	DeltaAfterIterations     int  `json:"delta_after_iterations"  usage:"delta after each iterations"`
	Delta                    int  `json:"delta"  usage:"delta"`
	MinValue                 int  `json:"min_value"  usage:"min value of parameter"`
	MaxValue                 int  `json:"max_value"  usage:"max value of parameter"`
	SnapBackAfterImprovement bool `json:"snap_back_after_improvement"  usage:"snap back to start value after improvement of best solution"`
	Zigzag                   bool `json:"zigzag"  usage:"zigzag between min and max value lik a jig saw"`
}

// SolverOptions are the options for the solver and it's operators.
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
