package mip

// Term is the building block of a constraint and an objective. A term consist
// of a coefficient and a var and should be interpreted as the product
// of coefficient and the var in the context of the constraint or
// objective.
type Term interface {
	// Coefficient returns the coefficient value of the invoking term.
	Coefficient() float64
	// Var returns the variable of the term.
	Var() Var
}

// Terms is a slice of Term instances.
type Terms []Term

// QuadraticTerm consists of a coefficient and two vars. It should be
// interpreted as the product of a coefficient and the two vars.
type QuadraticTerm interface {
	// Coefficient returns the coefficient value of the invoking term.
	Coefficient() float64
	// Var1 returns the first variable.
	Var1() Var
	// Var2 returns the second variable.
	Var2() Var
}

// QuadraticTerms is a slice of QuadraticTerm instances.
type QuadraticTerms []QuadraticTerm
