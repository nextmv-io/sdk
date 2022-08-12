package mip

// Objective specifies the objective of the definition. An objective
// consists out of terms and a specification if it should be maximized or
// minimized.
//
// For example:
//  maximize 2.5 * x + 3.5 * y
//
// 2.5 * x and 3.5 * y are 2 terms in this example.
type Objective interface {
	// AddTerm adds a term to the invoking objective, invoking this API
	// multiple times for the same variable will take the sum of coefficients
	// of earlier added terms for that variable
	//
	// 		d := mip.NewDefinition()
	// 		x, _ := d.AddContinuousVariable(10.0, 100.0)
	//
	// 		d.Objective().SetMaximize()			 // results in: maximize -
	// 		d.Objective().AddTerm(1.0, x)		// results in: maximize 1.0 * x
	// 		d.Objective().AddTerm(2.0, x)		// results in: maximize 3.0 * x
	AddTerm(coefficient float64, variable Variable) Term
	// IsMaximize return true if the invoking objective is a maximization objective.
	IsMaximize() bool
	// SetMaximize sets the invoking objective to be a maximization objective.
	SetMaximize()
	// SetMinimize sets the invoking objective to be a minimization objective.
	SetMinimize()
	// Terms returns a copy slice of terms of the invoking objective,
	// each variable is reported once. If the same variable has been
	// added multiple times the sum of coefficients is reported for that
	// variable.
	Terms() Terms
}
