package mip

// Definition manages the variables, constraints and objective.
type Definition interface {
	// AddBinaryVariable adds a binary variable to the invoking definition,
	// returns the newly constructed variable or error.
	AddBinaryVariable() (BinaryVariable, error)
	// AddContinuousVariable adds a continuous variable with bounds [loweBound,
	// upperBound] to the invoking definition, returns the newly constructed
	// variable or error.
	//
	// A continuous variable can take any value in the interval [lowerBound,
	// upperBound]
	AddContinuousVariable(
		lowerBound float64,
		upperBound float64,
	) (ContinuousVariable, error)
	// AddIntegerVariable adds an integer variable with bounds [loweBound,
	// upperBound] to the invoking definition, returns the newly constructed
	// variable or error.
	AddIntegerVariable(
		lowerBound int64,
		upperBound int64,
	) (IntegerVariable, error)
	// AddConstraint adds a constraint with sense and right-hand-side value rhs
	// to the invoking definition. All terms for existing and future variables
	// are initially zero. Returns the newly constructed constraint or error.
	// A constraint where all terms remain zero is ignored by the solver.
	AddConstraint(sense Sense, rhs float64) (Constraint, error)
	// Constraints returns a copy slice of all constraints.
	Constraints() Constraints
	// Objective returns the objective of the definition.
	Objective() Objective
	// Variables returns a copy slice of all variables.
	Variables() Variables
}
