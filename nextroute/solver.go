package nextroute

import (
	"context"
	"github.com/nextmv-io/sdk/connect"
)

type SolveOptions struct {
}

type Solver interface {
	Solve(ctx context.Context, solveOptions SolveOptions) (Solution, error)
}
type SolverOptions struct {
}

func NewSolver(
	solution Solution,
	options SolverOptions,
) (Solver, error) {
	connect.Connect(con, &newSolver)
	return newSolver(solution, options)
}
