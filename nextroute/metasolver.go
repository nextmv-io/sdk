package nextroute

import (
	"time"

	"github.com/nextmv-io/sdk/connect"
)

type MetaSolveOptions struct {
	MaximumDuration time.Duration `json:"maximum_duration"  usage:"maximum duration of solver in seconds"`
}

// MetaSolver is the interface for a meta solver.
type MetaSolver interface {
	Solve(solveOptions MetaSolveOptions) (Solution, error)
}

type MetaSolverOptions struct {
}

func NewMetaSolver(
	solution Solution,
) (MetaSolver, error) {
	connect.Connect(con, &newMetaSolver)
	return newMetaSolver(solution)
}
