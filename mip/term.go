package mip

// Term is the building block of a constraint and an objective. A term consist
// of a coefficient and a variable and should be interpreted as the product
// of coefficient and the variable in the context of the constraint or
// objective.
type Term interface {
	// Coefficient return the coefficient value of the invoking term.
	Coefficient() float64
	// Variable return the variable of the invoking term.
	Variable() Variable
}

// Terms is a slice of Term instances.
type Terms []Term
