package kmeans

import (
	"math/rand"
	"time"

	"github.com/nextmv-io/sdk/measure"
)

// Solver is the interface for a k-means solver.
type Solver interface {
	// Solve is the entrypoint to solve the model associated with
	// the invoking solver. Returns a solution when the invoking solver
	// reaches a conclusion.
	Solve(options SolveOptions) (Solution, error)
}

// SolveOptions is a set of options that can be used to influence
// the behavior of a solver.
type SolveOptions interface {
	// Candidates returns the number of candidate solutions to
	// consider when solving the model associated with the invoking
	// solver. Defaults to 1.
	Candidates() int
	// MaximumDuration returns the maximum duration to spend
	// solving the model associated with the invoking solver.
	// Defaults to 24 hour.
	MaximumDuration() time.Duration
	// Measure returns the measure used to calculate the distance
	// between points to derive the solution. Defaults to Euclidean
	// measure.
	Measure() measure.ByPoint
	// Random returns the random number generator used to derive the
	// solution. Defaults to a new random number generator seeded
	// with the current time.
	Random() *rand.Rand

	// SetCandidates sets the number of candidate solutions to
	// consider when solving the model associated with the invoking
	// solver. Returns the invoking solver options.
	SetCandidates(candidates int) SolveOptions
	// SetMaximumDuration sets the maximum duration to spend
	// solving the model associated with the invoking solver.
	// Returns the invoking solver options.
	SetMaximumDuration(maximumDuration time.Duration) SolveOptions
	// SetMeasure sets the measure used to calculate the distance
	// between points to derive the solution. Returns the invoking
	// solver options.
	SetMeasure(measure measure.ByPoint) SolveOptions
	// SetRandom sets the random number generator used to derive the
	// solution. Returns the invoking solver options.
	SetRandom(random *rand.Rand) SolveOptions
}
