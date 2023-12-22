package alns

import (
	"github.com/nextmv-io/sdk/connect"
)

// NewConstSolveParameter creates a new constant solve parameter for nextroute.
func NewConstSolveParameter(value int) SolveParameter {
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
) SolveParameter {
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
