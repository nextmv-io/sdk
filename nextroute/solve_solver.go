package nextroute

import (
	"context"
	"github.com/nextmv-io/sdk/connect"
	"math/rand"
	"time"
)

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

// NewSkeletonSolver creates a new solver for nextroute.
func NewSkeletonSolver(model Model) (Solver, error) {
	connect.Connect(con, &newSkeletonSolver)
	return newSkeletonSolver(model)
}

// SolveOptions holds the options for the solve process.
type SolveOptions struct {
	Iterations int           `json:"iterations"  usage:"maximum number of iterations, -1 assumes no limit" default:"-1"`
	Duration   time.Duration `json:"duration"  usage:"maximum duration of solver in seconds" default:"30s"`
}

// Solver is the interface for the Adaptive Local Neighborhood Search algorithm
// (ALNS) solver.
type Solver interface {
	Progressioner
	// AddSolveOperators adds a number of solve-operators to the solver.
	AddSolveOperators(...SolveOperator)

	// BestSolution returns the best solution found so far.
	BestSolution() Solution

	// HasBestSolution returns true if the solver has a best solution.
	HasBestSolution() bool
	// HasWorkSolution returns true if the solver has a work solution.
	HasWorkSolution() bool

	// Model returns the model used by the solver.
	Model() Model

	// Random returns the random number generator used by the solver.
	Random() *rand.Rand
	// Reset will reset the solver to use solution as work solution.
	Reset(solution Solution, solveInformation SolveInformation)

	// Solve starts the solving process using the given options. It returns the
	// solutions as a channel.
	Solve(context.Context, SolveOptions, ...Solution) (SolutionChannel, error)
	// SolveEvents returns the solve-events used by the solver.
	SolveEvents() SolveEvents
	// SolveOperators returns the solve-operators used by the solver.
	SolveOperators() SolveOperators

	// WorkSolution returns the current work solution.
	WorkSolution() Solution
}
