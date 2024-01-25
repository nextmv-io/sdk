package nextroute

import (
	"context"
	"math/rand"
	"time"

	"github.com/nextmv-io/sdk/connect"
)

// NewParallelSolver creates a new parallel solver for the given work solutions.
func NewParallelSolver(
	model Model,
) (ParallelSolver, error) {
	connect.Connect(con, &newParallelSolver)
	return newParallelSolver(model)
}

// NewSkeletonParallelSolver creates a new parallel solver for nextroute.
func NewSkeletonParallelSolver(model Model) (ParallelSolver, error) {
	connect.Connect(con, &newSkeletonParallelSolver)
	return newSkeletonParallelSolver(model)
}

// ParallelSolveOptions holds the options for the parallel solver.
type ParallelSolveOptions struct {
	Iterations           int           `json:"iterations"  usage:"maximum number of iterations, -1 assumes no limit; iterations are counted after start solutions are generated" default:"-1"`
	Duration             time.Duration `json:"duration" usage:"maximum duration of the solver" default:"30s"`
	ParallelRuns         int           `json:"parallel_runs" usage:"maximum number of parallel runs, -1 results in using all available resources" default:"-1"`
	StartSolutions       int           `json:"start_solutions" usage:"number of solutions to generate on top of those passed in; one solution generated with sweep algorithm, the rest generated randomly" default:"-1"`
	RunDeterministically bool          `json:"run_deterministically"  usage:"run the parallel solver deterministically"`
}

// DefaultSolveOptionsFactory creates a new SolveOptionsFactory.
// This factory is used by the parallel solver to create new solve options for
// a new run of a solver within the parallel solver. The default will
// limit the run to a random draw in the range [200, 2000] iterations and 30
// seconds.
func DefaultSolveOptionsFactory() SolveOptionsFactory {
	connect.Connect(con, &defaultSolveOptionsFactory)
	return defaultSolveOptionsFactory()
}

// NewSolveOptionsFactory creates a new SolveOptionsFactory which is independent
// of the state of the parallel solver. This factory will produce SolveOptions
// which limit a search run in a cycle to the given number of iterations and
// duration.
func NewSolveOptionsFactory(iterations int, duration time.Duration) SolveOptionsFactory {
	return func(_ ParallelSolveInformation) (SolveOptions, error) {
		solveOptions := SolveOptions{
			Iterations: iterations,
			Duration:   duration,
		}
		return solveOptions, nil
	}
}

// DefaultSolverFactory creates a new SolverFactory.
func DefaultSolverFactory() SolverFactory {
	connect.Connect(con, &defaultSolverFactory)
	return defaultSolverFactory()
}

// ParallelSolver is the interface for parallel solver. The parallel solver will
// run multiple solver in parallel and return the best solution. The parallel
// solver will stop when the maximum duration is reached.
type ParallelSolver interface {
	Progressioner

	// Model returns the model of the solver.
	Model() Model

	// SetSolverFactory sets the factory for creating new solver.
	SetSolverFactory(SolverFactory)
	// SetSolveOptionsFactory sets the factory for creating new solve options.
	SetSolveOptionsFactory(SolveOptionsFactory)
	// Solve starts the solving process using the given options. It returns the
	// solutions as a channel.
	Solve(context.Context, ParallelSolveOptions, ...Solution) (SolutionChannel, error)
	// SolveEvents returns the solve-events used by the solver.
	SolveEvents() SolveEvents
}

// SolveOptionsFactory is a factory type for creating new solve options.
// This factory is used by the parallel solver to create new solve options for
// a new run of a solver.
type SolveOptionsFactory func(
	information ParallelSolveInformation,
) (SolveOptions, error)

// SolverFactory is a factory type for creating new solver. This factory is
// used by the parallel solver to create new solver for a new run.
type SolverFactory func(
	information ParallelSolveInformation,
	solution Solution) (Solver, error)

// ParallelSolveInformation holds the information about the current parallel
// solve run.
type ParallelSolveInformation interface {
	// Cycle returns the current cycle. A cycle is a set of parallel runs.
	// In each cycle multiple runs are executed in parallel. Cycle identifies
	// how often a new solver has been created and started with the best
	// solution of the previous runs.
	Cycle() int

	// Random returns the random number generator from the solution.
	Random() *rand.Rand
	// Run returns the current run. A run is a single solve run. In each cycle
	// multiple runs are executed in parallel. Run identifies a run.
	Run() int
}
