package context

// Alternative is a function used for generating new contexts from an existing
// one. It receives the existing context as the only argument. You can use this
// lexical scope to wrangle variables and return a Statement that holds the
// correct set of instructions for generating those new contexts.
type Alternative func(Context) Statement

// Generator is a function for generating a new context.
type Generator func() Context

// Action holds the actions for generating a new context.
type Action interface {
	Use(Generator) Extend
	Return() Statement
	Discard() Statement
}

// Extend holds the condition that indicates whether the new generated context
// is operationally valid or not.
type Extend interface {
	With(Condition) Statement
}

// Statement is a lexical scope with precise instructions for generating new
// contexts from an existing one.
type Statement interface {
	Condition() Condition
	Generator() Generator
	Feasible() Condition
}

// When lets you specify under what condition you want to keep generating new
// contexts.
func When(c Condition) Action {
	return whenFunc(c)
}

// whenFunc holds the implementation of the When function.
var whenFunc func(Condition) Action
