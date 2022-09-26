package mip

// Objective specifies the objective of the model. An objective
// consists out of terms and a specification if it should be maximized or
// minimized.
//
// For example:
//
//	maximize 2.5 * x + 3.5 * y
//
// 2.5 * x and 3.5 * y are 2 terms in this example.
type Objective interface {
	// IsMaximize return true if the invoking objective is a maximization
	// objective.
	IsMaximize() bool
	// NewTerm adds a term to the invoking objective, invoking this API
	// multiple times for the same variable will take the sum of coefficients
	// of earlier added terms for that variable
	//
	// 		d := mip.NewModel()
	// 		x, _ := d.NewContinuousVar(10.0, 100.0)
	//
	// 		d.Objective().SetMaximize()			 // results in: maximize -
	// 		d.Objective().NewTerm(1.0, x)		// results in: maximize 1.0 * x
	// 		d.Objective().NewTerm(2.0, x)		// results in: maximize 3.0 * x
	NewTerm(coefficient float64, variable Var) Term
	// SetMaximize sets the invoking objective to be a maximization objective.
	SetMaximize()
	// SetMinimize sets the invoking objective to be a minimization objective.
	SetMinimize()
	// Term returns a term for variable with the sum of all coefficients of
	// defined terms for variable. The second return argument defines how many
	// terms have been defined on the objective for variable.
	Term(variable Var) (Term, int)
	// Terms returns a copy slice of terms of the invoking objective,
	// each variable is reported once. If the same variable has been
	// added multiple times the sum of coefficients is reported for that
	// variable. The order of the terms is not specified and is not guaranteed
	// to be the same from one invocation to the next.
	Terms() Terms
}
