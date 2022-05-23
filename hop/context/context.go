package context

type Context interface {
	Apply(...Change) Context
	Bound(Bounder) Context
	Check(...Condition) Context
	Format(Formatter) Context
	Propagate(...Propagator) Context
	Value(Valuer) Context
}

func NewContext() Context {
	load()
	return newContextFunc()
}

var newContextFunc func() Context
