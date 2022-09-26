package mip

// Sense defines the constraint operator between the left-hand-side
// and the right-hand-side.
type Sense int64

// Sense of a Constraint.
const (
	// LessThanOrEqual is used to define a less than or equal constraint
	// 		c, _ := d.NewConstraint(mip.LessThanOrEqual, 123.4)
	//
	// 		c.NewTerm(1.0, x)  	 // results in 1.0 * x <= 123.4 in solver
	LessThanOrEqual Sense = iota
	// Equal is used to define an equality constraint
	// 		c, _ := d.NewConstraint(mip.Equal, 123.4)
	//
	// 		c.NewTerm(1.0, x)  	 // results in 1.0 * x = 123.4 in solver
	Equal
	// GreaterThanOrEqual is used to define a greater or equal constraint
	// 		c, _ := d.NewConstraint(mip.	GreaterThanOrEqual, 123.4)
	//
	// 		c.NewTerm(1.0, x)  	 // results in 1.0 * x >= 123.4 in solver
	GreaterThanOrEqual
)

// Constraint specifies a relation between variables a solution has to comply
// with. A constraint consists out of terms, a sense and a right hand side.
//
// For example:
//
//	2.5 * x + 3.5 * y <= 10.0
//
// The less than operator is the sense
// The value 10.0 is the right hand side
//
//	2.5 * x and 3.5 * y are 2 terms in this example
type Constraint interface {
	// Name returns assigned name. If no name has been set it will return
	// a unique auto-generated name.
	Name() string
	// NewTerm adds a term to the invoking constraint, invoking this API
	// multiple times for the same variable will take the sum of coefficients
	// of earlier added terms for that variable
	//
	// 		d := mip.NewModel()
	//
	// 		x, _ := d.NewContinuousVar(10.0, 100.0)
	//
	// 		c, _ := d.NewConstraint(mip.LessThanOrEqual, 123.4)
	//
	// 		c.NewTerm(1.0, x)  	 // results in 1.0 * x <= 123.4 in solver
	// 		c.NewTerm(2.0, x)    // results in 3.0 * x <= 123.4 in solver
	NewTerm(coefficient float64, variable Var) Term
	// RightHandSide returns the right-hand side of the invoking constraint.
	RightHandSide() float64
	// Sense returns the sense of the invoking constraint.
	Sense() Sense
	// SetName assigns name to invoking constraint
	SetName(name string)
	// Term returns a term for variable with the sum of all coefficients of
	// defined terms for variable. The second return argument defines how many
	// terms have been defined on the objective for variable.
	Term(variable Var) (Term, int)
	// Terms returns a copy slice of terms of the invoking constraint,
	// each variable is reported once. If the same variable has been
	// added multiple times the sum of coefficients is reported for that
	// variable.
	Terms() Terms
}

// Constraints slice of Constraint instances.
type Constraints []Constraint
