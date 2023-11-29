package alns

import (
	sdkAlns "github.com/nextmv-io/sdk/alns"
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute"
)

// NewParallelSolver creates a new parallel solver for nextroute.
func NewParallelSolver() sdkAlns.ParallelSolver[nextroute.Solution] {
	connect.Connect(con, &newParallelSolver)
	return newParallelSolver()
}
