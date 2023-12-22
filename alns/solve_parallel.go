package alns

import "math/rand"

// NewSolveOptionsFactory is a factory type for creating new solve options.
// This factory is used by the parallel solver to create new solve options for
// a new run of a solver.
type NewSolveOptionsFactory[T Solution[T]] func(
	information ParallelSolveInformation,
) SolveOptions

// NewSolverFactory is a factory type for creating new solver. This factory is
// used by the parallel solver to create new solver for a new run.
type NewSolverFactory[T Solution[T]] func(
	information ParallelSolveInformation,
	solution T) Solver[T, SolveOptions]

// ParallelSolveInformation holds the information about the current parallel
// solve run.
type ParallelSolveInformation interface {
	// Run returns the current run. A run is a single solve run. In each cycle
	// multiple runs are executed in parallel. Run identifies a run.
	Run() int
	// Cycle returns the current cycle. A cycle is a set of parallel runs.
	// In each cycle multiple runs are executed in parallel. Cycle identifies
	// how often a new solver has been created and started with the best
	// solution of the previous runs.
	Cycle() int

	// Random returns the random number generator from the solution.
	Random() *rand.Rand
}

// ParallelSolveOptions holds the options for the parallel solver.
type ParallelSolveOptions[T Solution[T]] interface {
	// Iterations returns the maximum number of iterations of the parallel
	// solver.
	Iterations() int
	// ParallelRuns returns the maximum number of parallel runs.
	ParallelRuns() int
	// RunDeterministically returns true if the parallel solver should run
	// deterministically. Deterministic mode will sacrifice performance for
	// determinism.
	RunDeterministically() bool
}

// ParallelSolver is the interface for parallel solver. The parallel solver will
// run multiple solver in parallel and return the best solution. The parallel
// solver will stop when the maximum duration is reached.
type ParallelSolver[T Solution[T]] interface {
	Progressioner
	BaseSolver[T, ParallelSolveOptions[T]]
	SetSolverFactory(NewSolverFactory[T])
	SetSolveOptionsFactory(NewSolveOptionsFactory[T])
	// SolveEvents returns the solve-events used by the solver.
	SolveEvents() SolveEvents[T]
}
