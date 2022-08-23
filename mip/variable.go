package mip

// Variable represents the entities on which the solver has to make a decision
// without violating constraints and optimize the objective.
// Variables can be of a certain type, continuous, C or binary.
//
// Continuous variables can take a value of a continuous quantity
// Integer variables are variables that must take an integer value
// (0, 1, 2, ...)
// Binary variables can take two values, zero or one.
type Variable interface {
	// Index is a unique number assigned to the variable. The index corresponds
	// to the location in the slice returned by Model.Variables().
	Index() int
	// IsBinary returns true if the invoking variable is a binary variable.
	// otherwise false
	IsBinary() bool
	// IsContinuous returns true if the invoking variable is a continuous
	// variable otherwise false.
	IsContinuous() bool
	// IsInteger returns true if the invoking variable is an integer variable
	// otherwise false.
	IsInteger() bool
	// LowerBound returns the lowerBound of the invoking variable. By definition
	// this is 0.0 for a binary variable
	//
	// Lower bounds of variables are limited by the lower bounds of the
	// underlying solver technology. The lower bound used will be the maximum
	// of the specification and the lower bound of the solver used.
	LowerBound() float64
	// UpperBound returns the upperBound of the invoking variable. By definition
	// this is 1.0 for a binary variable
	//
	// Upper bounds of variables are limited by the upper bounds of the
	// underlying solver technology. The upper bound used will be the minimum
	// of the specification and the upper bound of the solver used.
	UpperBound() float64
}

// Variables slice of Variable instances.
type Variables []Variable

// ContinuousVariable a Variable which can take any value in an interval.
type ContinuousVariable interface {
	Variable
	ensureContinuous() bool
}

// IntegerVariable a Variable which can take any integer value in an interval.
type IntegerVariable interface {
	Variable
	ensureInteger() bool
}

// BinaryVariable a Variable which can take two values, zero or one. A binary
// variable is also an integer variable which can have two values zero and
// one.
type BinaryVariable interface {
	IntegerVariable
	ensureBinary() bool
}
