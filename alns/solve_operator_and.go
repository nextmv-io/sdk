package alns

// SolveOperatorAnd is a solve-operator.
type SolveOperatorAnd[T Solution[T]] interface {
	SolveOperator[T]

	// Operators returns the solve-operators that will be executed in each
	// iteration.
	Operators() SolveOperators[T]
}
