package context

// Condition represents a logical condition on a context.
type Condition func(Context) bool

// And uses the conditional "AND" logical operator on all given conditions. It
// returns true if all conditions are true.
func And(c1 Condition, c2 Condition, conditions ...Condition) Condition {
	return andFunc(c1, c2, conditions...)
}

// False is a convenience function that is always false.
func False(ctx Context) bool {
	return falseFunc(ctx)
}

// Not negates the given condition.
func Not(c Condition) Condition {
	return notFunc(c)
}

// Or uses the conditional "OR" logical operator on all given conditions. It
// returns true if at least one condition is true.
func Or(c1 Condition, c2 Condition, conditions ...Condition) Condition {
	return orFunc(c1, c2, conditions...)
}

// True is a convenience function that is always true.
func True(ctx Context) bool {
	return trueFunc(ctx)
}

// Xor uses the conditional "Exclusive OR" logical operator on all given
// conditions. It returns true if, and only if, the conditions are different.
func Xor(c1, c2 Condition) Condition {
	return xorFunc(c1, c2)
}

var (
	andFunc   func(Condition, Condition, ...Condition) Condition
	falseFunc func(Context) bool
	notFunc   func(Condition) Condition
	orFunc    func(Condition, Condition, ...Condition) Condition
	trueFunc  func(Context) bool
	xorFunc   func(Condition, Condition) Condition
)
