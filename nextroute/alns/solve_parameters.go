package alns

import (
	sdkAlns "github.com/nextmv-io/sdk/alns"
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute"
)

// NewConstSolveParameter creates a new constant solve parameter for nextroute.
func NewConstSolveParameter(value int) sdkAlns.SolveParameter[nextroute.Solution] {
	connect.Connect(con, &newConstSolveParameter)
	return newConstSolveParameter(value)
}

// NewSolveParameter creates a new solve parameter for nextroute.
func NewSolveParameter(
	startValue int,
	deltaAfterIterations int,
	delta int,
	minValue int,
	maxValue int,
	snapBackAfterImprovement bool,
	zigzag bool,
) sdkAlns.SolveParameter[nextroute.Solution] {
	connect.Connect(con, &newSolveParameter)
	return newSolveParameter(
		startValue,
		deltaAfterIterations,
		delta,
		minValue,
		maxValue,
		snapBackAfterImprovement,
		zigzag,
	)
}
