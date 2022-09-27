package mip

// Var represents the entities on which the solver has to make a decision
// without violating constraints and while optimizing the objective.
// Vars can be of a certain type, bool, float or int.
//
// Float vars can take a value of a float quantity
// Int vars are vars that must take an integer value
// (0, 1, 2, ...)
// Bool vars can take two values, zero or one.
type Var interface {
	// Index is a unique number assigned to the var. The index corresponds
	// to the location in the slice returned by Model.Variables().
	Index() int
	// IsBool returns true if the invoking variable is a bool variable,
	// otherwise it returns false.
	IsBool() bool
	// IsFloat returns true if the invoking variable is a float
	// variable otherwise false.
	IsFloat() bool
	// IsInt returns true if the invoking variable is an int variable
	// otherwise false.
	IsInt() bool
	// LowerBound returns the lowerBound of the invoking variable.
	//
	// Lower bounds of variables are limited by the lower bounds of the
	// underlying solver technology. The lower bound used will be the maximum
	// of the specification and the lower bound of the solver used.
	LowerBound() float64
	// Name returns assigned name. If no name has been set it will return
	// a unique auto-generated name.
	Name() string
	// SetName assigns name to invoking var
	SetName(name string)
	// UpperBound returns the upperBound of the invoking variable.
	//
	// Upper bounds of variables are limited by the upper bounds of the
	// underlying solver technology. The upper bound used will be the minimum
	// of the specification and the upper bound of the solver used.
	UpperBound() float64
}

// Vars is a slice of Var instances.
type Vars []Var

// Float a Var which can take any value in an interval.
type Float interface {
	Var
	ensureFloat() bool
}

// Int a Var which can take any integer value in an interval.
type Int interface {
	Var
	ensureInt() bool
}

// Bool a Var which can take two values, zero or one. A bool
// variable is also an int variable which can have two values zero and
// one.
type Bool interface {
	Int
	ensureBool() bool
}
