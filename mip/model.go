package mip

// Model manages the variables, constraints and objective.
type Model interface {
	// Constraints returns a copy slice of all constraints.
	Constraints() Constraints
	// NewBinaryVar adds a binary variable to the invoking model,
	// returns the newly constructed variable.
	NewBinaryVar() BinaryVar
	// NewContinuousVar adds a continuous var with bounds [loweBound,
	// upperBound] to the invoking model, returns the newly constructed
	// var.
	NewContinuousVar(
		lowerBound float64,
		upperBound float64,
	) ContinuousVar
	// NewIntegerVar adds an integer var with bounds [loweBound,
	// upperBound] to the invoking model, returns the newly constructed
	// var.
	NewIntegerVar(
		lowerBound int64,
		upperBound int64,
	) IntegerVar
	// NewConstraint adds a constraint with sense and right-hand-side value rhs
	// to the invoking model. All terms for existing and future variables
	// are initially zero. Returns the newly constructed constraint.
	// A constraint where all terms remain zero is ignored by the solver.
	NewConstraint(sense Sense, rhs float64) Constraint
	// Objective returns the objective of the model.
	Objective() Objective
	// Vars returns a copy slice of all vars.
	Vars() Vars
}
