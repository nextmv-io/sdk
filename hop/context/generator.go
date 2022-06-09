package context

// A Generator generates new contexts from an existing one.
type Generator interface {
	Condition() Condition
	Generate() func() Context
	Feasible() Condition
	With(Condition) Generator
}

// Action holds the actions for generating a new context.
type Action interface {
	Then(func() Context) Generator
	Return() Generator
	Discard() Generator
}

// If lets you specify under what condition you want to keep generating new
// contexts.
func If(c Condition) Action {
	return ifFunc(c)
}

// Scope lets you create a new generator using the same lexical scope.
func Scope(f func(Context) Generator) Generator {
	return scopeFunc(f)
}

// ifFunc holds the implementation of the If function.
var ifFunc func(Condition) Action

// scopeFunc holds the implementation of the Scope function.
var scopeFunc func(func(Context) Generator) Generator
