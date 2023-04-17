package nextroute

import (
	"context"
	"time"

	"github.com/nextmv-io/sdk/alns"
	"github.com/nextmv-io/sdk/connect"
)

// ParallelSolveOptions are the options for the parallel solver.
type ParallelSolveOptions struct {
	MaximumDuration      time.Duration `json:"maximum_duration" usage:"maximum duration of the solver" default:"30s"`
	MaximumParallelRuns  int           `json:"maximum_parallel_runs" usage:"maximum number of parallel runs, -1 results in using all available resources" default:"100"`
	StartSolutions       int           `json:"start_solutions" usage:"number of solutions to start with; one solution generated with sweep algorithm, the rest generated randomly" default:"0"`
	UseNVehicles         int           `json:"use_n_vehicles"  usage:"use n vehicles to generate one of the initial solutions"`
	IterationMultiplier  int           `json:"iteration_multiplier"  usage:"number of iterations per solver; this gets multiplied by a random value between 1 and 10" default:"200"`
	RunDeterministically bool          `json:"run_deterministically"  usage:"run the parallel solver deterministically"`
}

// ParallelSolver is the interface for parallel solver. The parallel solver will
// run multiple solver in parallel and return the best solution. The parallel
// solver will stop when the maximum duration is reached.
type ParallelSolver interface {
	alns.Progressioner
	alns.Solver[Solution, ParallelSolveOptions]
}

// NewParallelSolver creates a new parallel solver for the given work solutions.
func NewParallelSolver(
	model Model,
) (ParallelSolver, error) {
	connect.Connect(con, &newParallelSolver)
	return newParallelSolver(model)
}

// Solve a nextroute vehicle routing problem.
func Solve(ctx context.Context, solver ParallelSolver, options ParallelSolveOptions, solutions chan<- any) error {
	connect.Connect(con, &solve)
	return solve(ctx, solver, options, solutions)
}
