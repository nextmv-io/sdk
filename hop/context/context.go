package context

// Context represents a context for making decisions.
type Context interface {
	Apply(...Change) Context
	Bound(Bounder) Context
	Check(...Condition) Context
	Format(Formatter) Context
	Propagate(...Propagator) Context
	Value(Valuer) Context
}

// NewContext returns a new Context.
func NewContext() Context {
	load()
	return newContextFunc()
}

var newContextFunc func() Context
