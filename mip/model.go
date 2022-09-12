package mip

// Model manages the variables, constraints and objective.
type Model interface {
	// Constraints returns a copy slice of all constraints.
	Constraints() Constraints
	// NewBinaryVar adds a binary variable to the invoking model,
	// returns the newly constructed variable or error.
	NewBinaryVar() (BinaryVar, error)
	// NewContinuousVar adds a continuous var with bounds [loweBound,
	// upperBound] to the invoking model, returns the newly constructed
	// var or error.
	//
	// A continuous var can take any value in the interval [lowerBound,
	// upperBound]
	NewContinuousVar(
		lowerBound float64,
		upperBound float64,
	) (ContinuousVar, error)
	// NewIntegerVar adds an integer var with bounds [loweBound,
	// upperBound] to the invoking model, returns the newly constructed
	// var or error.
	NewIntegerVar(
		lowerBound int64,
		upperBound int64,
	) (IntegerVar, error)
	// NewConstraint adds a constraint with sense and right-hand-side value rhs
	// to the invoking model. All terms for existing and future variables
	// are initially zero. Returns the newly constructed constraint or error.
	// A constraint where all terms remain zero is ignored by the solver.
	NewConstraint(sense Sense, rhs float64) (Constraint, error)
	// Objective returns the objective of the model.
	Objective() Objective
	// Vars returns a copy slice of all vars.
	Vars() Vars
}
