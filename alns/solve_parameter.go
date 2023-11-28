package alns

// SolveParameter is an interface for a parameter that can change
// during the solving. The parameter can be used to control the
// behavior of the solver and it's operators.
type SolveParameter[T Solution[T]] interface {
	// Update updates the parameter based on the given solve information.
	// Update is invoked after each iteration of the solver.
	Update(SolveInformation[T])

	// Value returns the current value of the parameter.
	Value() int
}

// SolveParameters is a slice of solve parameters.
type SolveParameters[T Solution[T]] []SolveParameter[T]
