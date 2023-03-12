package nextroute

import (
	"context"
	"time"

	"github.com/nextmv-io/sdk/connect"
)

type ParallelSolveOptions struct {
	MaximumDuration     time.Duration `json:"maximum_duration"  usage:"maximum duration of solver in seconds"`
	MaximumParallelRuns int           `json:"maximum_parallel_runs"  usage:"maximum number of parallel runs, -1 implies using all available resources"`
}

// ParallelSolver is the interface for parallel solver. The parallel solver will
// run multiple solver in parallel and return the best solution. The parallel
// solver will stop when the maximum duration is reached.
type ParallelSolver interface {
	Progressioner
	Solve(ctx context.Context, solveOptions ParallelSolveOptions) (Solution, error)
}

// NewParallelSolver creates a new parallel solver for the given work solutions.
func NewParallelSolver(
	model Model,
) (ParallelSolver, error) {
	connect.Connect(con, &newParallelSolver)
	return newParallelSolver(model)
}
