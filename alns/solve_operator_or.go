package alns

// SolveOperatorOr is a solve-operator.
// A solver operator or is a solve-operator that executes n loops in each
// iteration of a solver. In each loop a random solve-operator is selected
// using the probability of the solve-operator. If there are 3 operators
// one with probability 0.1, one with probability 0.2 and one with probability
// 0.3 then the first operator has a 0.1/0.6 chance to be selected, the second
// operator has a 0.2/0.6 chance to be selected and the third operator has a
// 0.3/0.6 chance to be selected.
type SolveOperatorOr[T Solution[T]] interface {
	SolveOperator[T]

	// Operators returns the solve-operators one will be selected from in
	// each loop.
	Operators() SolveOperators[T]
}
