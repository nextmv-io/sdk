package store

const (
	// Minimizer indicates the solution space is being searched to find the
	// smallest possible value.
	Minimizer Sense = iota
	// Maximizer indicates the solution space is being searched to find the
	// biggest possible value.
	Maximizer
	// Satisfier indicates the solution space is being searched to find
	// operationally valid Stores.
	Satisfier
)

func (s Sense) String() string {
	switch s {
	case Minimizer:
		return "minimizer"
	case Maximizer:
		return "maximizer"
	case Satisfier:
		return "satisfier"
	default:
		return ""
	}
}

// New returns a new Store.
func New() Store {
	connect()
	return newFunc()
}

// And uses the conditional "AND" logical operator on all given conditions. It
// returns true if all conditions are true.
func And(c1 Condition, c2 Condition, conditions ...Condition) Condition {
	return andFunc(c1, c2, conditions...)
}

// False is a convenience function that is always false.
func False(s Store) bool {
	return falseFunc(s)
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
func True(s Store) bool {
	return trueFunc(s)
}

// Xor uses the conditional "Exclusive OR" logical operator on all given
// conditions. It returns true if, and only if, the conditions are different.
func Xor(c1, c2 Condition) Condition {
	return xorFunc(c1, c2)
}

/*
If specifies under what condition a Generator can be used.

	s := store.New()
	x := store.NewVar(s, 1)
	s = s.Generate(
		// This generator can always be used.
		store.If(store.True).Return(),
		// This generator can never be used.
		store.If(store.False).Discard(),
		// This generator can only be used if the condition holds.
		store.If(func(s store.Store) bool { return x.Get(s) < 10 }).Return(),
	)
*/
func If(c Condition) Action {
	return ifFunc(c)
}

/*
Scope specifies a Generator that allows the use of the same lexical scope. This
is useful for reusing Variables and calculations among functions.

	s := store.New()
	x := store.NewVar(s, 1)
	s = s.Generate(
		store.Scope(func(s store.Store) store.Generator {
			v := x.Get(s)
			return store.If(func(s store.Store) bool {
				// v is used here.
				return v < 10
			}).Then(func(s store.Store) store.Store {
				// v is also used here.
				v++
				return s.Apply(x.Set(v))
			})
		}),
	)
*/
func Scope(f func(Store) Generator) Generator {
	return scopeFunc(f)
}

/*
DefaultOptions for running a solver. Options can be customized after using
these sensitive defaults.

	opt := store.DefaultOptions()
	opt.Limits.Duration = time.Duration(5) * time.Second
*/
func DefaultOptions() Options {
	return defaultOptionsFunc()
}

var (
	newFunc            func() Store
	andFunc            func(Condition, Condition, ...Condition) Condition
	falseFunc          func(Store) bool
	notFunc            func(Condition) Condition
	orFunc             func(Condition, Condition, ...Condition) Condition
	trueFunc           func(Store) bool
	xorFunc            func(Condition, Condition) Condition
	ifFunc             func(Condition) Action
	scopeFunc          func(func(Store) Generator) Generator
	defaultOptionsFunc func() Options
)
