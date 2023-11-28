package alns

import (
	"context"
)

// SolveOperator is a solve-operator. A solve-operator is a function that
// modifies the current solution. The function is executed with a certain
// probability. The probability is set by the SetProbability method. The
// probability is used by the solver to determine if the solve-operator should
// be executed. The probability is a value between 0 and 1. The manipulation of
// the solution is implemented in the Execute method. The Execute method is
// called by the solver. The Execute method is passed a SolveInformation
// instance that contains information about the current solution and the
// solver. The Execute method should modify the current solution. The
// Execute method should not modify the SolveInformation instance. The Execute
// method should not modify the SolveOperator instance.
type SolveOperator[T Solution[T]] interface {
	// CanResultInImprovement returns true if the solve-operator can result in
	// an improvement compared to the best solution. The solver uses this
	// information to determine if the best solution should be replaced by
	// the new solution.
	CanResultInImprovement() bool

	// Execute executes the solve-operator. The Execute method is called by the
	// solver. The Execute method is passed a SolveInformation instance that
	// contains information about the current solution and the solver. The
	// Execute method should modify the current solution.
	Execute(context.Context, SolveInformation[T])

	// Index returns the unique index of the solve-operator.
	Index() int

	// Probability returns the probability of the solve-operator.
	// The probability is a value between 0 and 1. The solver uses the
	// probability to determine if the solve-operator should be executed in
	// the current iteration. Each iteration the solver will execute a
	// solve-operator with a probability equal to the probability of the
	// solve-operator. The probability is set by the SetProbability method.
	Probability() float64

	// SetProbability sets the probability of the solve-operator.
	SetProbability(probability float64) SolveOperator[T]
	// SolveParameters returns the solve-parameters of the solve-operator.
	SolveParameters() SolveParameters[T]
}

// InterestedInBetterSolution is an interface that can be implemented by
// solve-operators that are interested in being notified when a better solution
// is found. The solver will call the OnBetterSolution method when a better
// best-solution is found.
type InterestedInBetterSolution[T Solution[T]] interface {
	OnBetterSolution(SolveInformation[T])
}

// InterestedInStartSolve is an interface that can be implemented by
// solve-operators that are interested in being notified when the solver starts
// solving. The solver will call the OnStartSolve method when the solver starts
// solving.
type InterestedInStartSolve[T Solution[T]] interface {
	OnStartSolve(SolveInformation[T])
}

// SolveOperators is a slice of solve-operators.
type SolveOperators[T Solution[T]] []SolveOperator[T]

// SolveOperatorOption is a function that configures a solve-operator.
type SolveOperatorOption[T Solution[T]] func(s SolveOperator[T]) error
