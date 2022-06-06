package context

import "github.com/nextmv-io/sdk/hop/solve"

// Context represents a context for making decisions.
type Context interface {
	// Methods that change search logic and representation of a context.
	Bound(Bounder) Context
	Check(...Condition) Context
	Format(Formatter) Context
	Value(Valuer) Context

	// Methods that change variables in a context.
	Apply(...Change) Context
	Propagate(...Propagator) Context

	// Methods that build a solver to search the space defined by a context.
	Maximizer() solve.Solver
	Minimizer() solve.Solver
	Satisfier() solve.Solver
}

// NewContext returns a new Context.
func NewContext() Context {
	connect()
	return newContextFunc()
}

var newContextFunc func() Context
