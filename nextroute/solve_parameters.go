package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
)

// SolveParameter is an interface for a parameter that can change
// during the solving. The parameter can be used to control the
// behavior of the solver and it's operators.
type SolveParameter interface {
	// Update updates the parameter based on the given solve information.
	// Update is invoked after each iteration of the solver.
	Update(SolveInformation)

	// Value returns the current value of the parameter.
	Value() int
}

// SolveParameters is a slice of solve parameters.
type SolveParameters []SolveParameter

// NewConstSolveParameter creates a new constant solve parameter for nextroute.
// A const solve parameter is a parameter that does not change during the
// solving. SolveParameter.Value will always return the same value.
func NewConstSolveParameter(value int) SolveParameter {
	connect.Connect(con, &newConstSolveParameter)
	return newConstSolveParameter(value)
}

// NewSolveParameter creates a new solve parameter for nextroute.
// A solve parameter is a parameter that can change during the solving. The
// parameter can be used to control the behavior of the solver and it's
// operators. This solve-parameter will change its value after a given number
// of iterations. The parameter will change its value by a given delta. The
// parameter will never be smaller than the given min value and never be bigger
// than the given max value. If snapBackAfterImprovement is true the parameter
// will snap back to the start value after an improvement of the best solution.
// If zigzag is true the parameter will zigzag between min and max value like a
// jig saw.
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
