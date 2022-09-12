package mip

// Term is the building block of a constraint and an objective. A term consist
// of a coefficient and a var and should be interpreted as the product
// of coefficient and the var in the context of the constraint or
// objective.
type Term interface {
	// Coefficient returns the coefficient value of the invoking term.
	Coefficient() float64
	// Var returns the var of the invoking term.
	Var() Var
}

// Terms is a slice of Term instances.
type Terms []Term
