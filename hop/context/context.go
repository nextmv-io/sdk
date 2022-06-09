package context

import "github.com/nextmv-io/sdk/hop/solve"

// Context represents a context for making decisions.
type Context interface {
	// Apply changes to a context.
	Apply(...Change) Context
	// Bound the value of a context.
	Bound(Bounder) Context
	// Format a context into any structure prior to JSON encoding.
	Format(Formatter) Context
	// Propagate lets you propagate into multiple contexts.
	Propagate(...Propagator) Context
	// Value lets you specify the numerical value of a context.
	Value(Valuer) Context
	// Generate lets you genereate new contexts from different alternatives.
	Generate(...Generator) Context

	// Methods that build a solver to search the space defined by a context.
	Maximizer(solve.Options) solve.Solver
	Minimizer(solve.Options) solve.Solver
	Satisfier(solve.Options) solve.Solver
}

// NewContext returns a new Context.
func NewContext() Context {
	connect()
	return newContextFunc()
}

var newContextFunc func() Context
