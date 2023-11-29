package alns

import (
	sdkAlns "github.com/nextmv-io/sdk/alns"
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute"
)

// NewSolver creates a new solver for nextroute.
func NewSolver() (sdkAlns.Solver[nextroute.Solution, sdkAlns.SolveOptions], error) {
	connect.Connect(con, &newSolver)
	return newSolver()
}
