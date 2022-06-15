package store

import (
	"github.com/nextmv-io/sdk/hop/store/types"
)

// New returns a new Store.
func New() types.Store {
	connect()
	return newFunc()
}

// And uses the conditional "AND" logical operator on all given conditions. It
// returns true if all conditions are true.
func And(
	c1 types.Condition,
	c2 types.Condition,
	conditions ...types.Condition,
) types.Condition {
	return andFunc(c1, c2, conditions...)
}

// False is a convenience function that is always false.
func False(s types.Store) bool {
	return falseFunc(s)
}

// Not negates the given condition.
func Not(c types.Condition) types.Condition {
	return notFunc(c)
}

// Or uses the conditional "OR" logical operator on all given conditions. It
// returns true if at least one condition is true.
func Or(
	c1 types.Condition,
	c2 types.Condition,
	conditions ...types.Condition,
) types.Condition {
	return orFunc(c1, c2, conditions...)
}

// True is a convenience function that is always true.
func True(s types.Store) bool {
	return trueFunc(s)
}

// Xor uses the conditional "Exclusive OR" logical operator on all given
// conditions. It returns true if, and only if, the conditions are different.
func Xor(c1, c2 types.Condition) types.Condition {
	return xorFunc(c1, c2)
}

/*
If specifies under what condition a Generator can be used.

	s := store.New()
	x := store.Var(s, 1)
	s = s.Generate(
		// This generator can always be used.
		store.If(store.True).Return(),
		// This generator can never be used.
		store.If(store.False).Discard(),
		// This generator can only be used if the condition holds.
		store.If(func(s types.Store) bool { return x.Get(s) < 10 }).Return(),
	)
*/
func If(c types.Condition) types.Action {
	return ifFunc(c)
}

/*
Scope specifies a Generator that allows the use of the same lexical scope. This
is useful for reusing Variables and calculations among functions.

	s := store.New()
	x := store.Var(s, 1)
	s = s.Generate(
		store.Scope(func(s types.Store) types.Generator {
			v := x.Get(s)
			return store.If(func(s types.Store) bool {
				// v is used here.
				return v < 10
			}).Then(func(s types.Store) types.Store {
				// v is also used here.
				v++
				s.Apply(x.Set(v))
				return s
			})
		}),
	)
*/
func Scope(f func(types.Store) types.Generator) types.Generator {
	return scopeFunc(f)
}

var (
	newFunc func() types.Store
	andFunc func(
		types.Condition,
		types.Condition,
		...types.Condition,
	) types.Condition
	falseFunc func(types.Store) bool
	notFunc   func(types.Condition) types.Condition
	orFunc    func(
		types.Condition,
		types.Condition,
		...types.Condition,
	) types.Condition
	trueFunc  func(types.Store) bool
	xorFunc   func(types.Condition, types.Condition) types.Condition
	ifFunc    func(types.Condition) types.Action
	scopeFunc func(func(types.Store) types.Generator) types.Generator
)
