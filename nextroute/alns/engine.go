// Package alns is a package
package alns

import (
	sdkAlns "github.com/nextmv-io/sdk/alns"
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute"
)

var (
	con = connect.NewConnector("sdk", "NextRouteAlns")

	newSolverOperatorAnd func(
		probability float64,
		operators sdkAlns.SolveOperators[nextroute.Solution],
	) (sdkAlns.SolveOperatorAnd[nextroute.Solution], error)

	newSolverOperatorOr func(
		loops int,
		probability float64,
		operators sdkAlns.SolveOperators[nextroute.Solution],
	) (sdkAlns.SolveOperatorOr[nextroute.Solution], error)

	newSolveOperatorIndex func() int

	newSolveOperator func(
		probability float64,
		canResultInImprovement bool,
		parameters sdkAlns.SolveParameters[nextroute.Solution],
	) sdkAlns.SolveOperator[nextroute.Solution]

	newParallelSolver func() sdkAlns.ParallelSolver[nextroute.Solution]

	newConstSolveParameter func(int) SolveParameter

	newSolveParameter func(
		startValue int,
		deltaAfterIterations int,
		delta int,
		minValue int,
		maxValue int,
		snapBackAfterImprovement bool,
		zigzag bool,
	) SolveParameter

	newSolver func() (sdkAlns.Solver[nextroute.Solution, sdkAlns.SolveOptions], error)
)
