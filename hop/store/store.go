package store

import (
	"github.com/nextmv-io/sdk/hop/model"
	"github.com/nextmv-io/sdk/hop/solve"
)

/*
Store represents a store of variables and logic to solve decision automation
problems. Adding logic to the store updates it:

    s := store.NewStore()
    s = s.Apply(...)
    s = s.Bound(...)
    s = s.Format(...)
    s = s.Generate(...)

functions may be called directly and chained.

    s := store.NewStore().
        Apply(...).
        Generate(...).
        Bound(...).
        Format(...)

The variables and logic stored define a solution space and Hop searches this
space to make decisions.
*/
type Store interface {
	/*
		Apply changes to a Store. A change happens when a stored variable is
		updated:

			s := store.NewStore()
			x := store.Var(s, 3.1416)
			s = s.Apply(
				x.Set(x.Get(s) * 2),
			)
	*/
	Apply(...Change) Store

	/*
		Bound the value of a Store. The lower and upper bounds can be set:

				s := store.NewStore()
				x := store.Var(s, 10)
				s = s.Bound(func(s store.Store) model.Bounds {
					return model.Bounds{
						Lower: x.Get(s),
						Upper: math.MaxInt,
					}
				})
	*/
	Bound(Bounder) Store

	/*
		Format a Store into any structure prior to JSON encoding.

				s := store.NewStore()
				x := store.Var(s, 10)
				s = s.Format(func(s store.Store) any {
					return map[string]int{"x": x.Get(s)}
				})

		Output after running:

				{"x": 10}
	*/
	Format(Formatter) Store

	/*
		Generate new Stores. A Generator is used until the generating condition
		no longer holds. In that case, the next Generator is evaluated. This is
		repeated until all Generators are exhausted.

		     s := store.NewStore()
		     x := store.Var(s, 10)
		     s = s.Generate(
		         store.If(func(s store.Store) bool {
		             // If the value of x is smaller than 10.
		             return x.Get(s) < 10
		         }).Then(func(s store.Store) store.Store {
		             // Then increase the value of x by 1.
		             v := x.Get(s)
		             v++
		             s.Apply(x.Set(v))
		             return s
		         }).With(func(s store.Store) bool {
		             // Operationally valid if x is even.
		             return x.Get(s)%2 == 0
		         }),
		     )

		There are convenience exit functions to stop generating new Stores
		and set their operational validity.

		     s := store.NewStore()
		     x := store.Var(s, 10)
		     s = s.Generate(
		         // If x is greater than 20 stop the search and make it
				 // operationally invalid.
		         store.If(func(s store.Store) bool { return x.Get(s) > 20 }).
		             Discard(),
		         // If x is a multiple of 5 and 3 stop the search and make it
				 // operationally valid.
		         store.If(store.And(
		             func(s store.Store) bool { return x.Get(s)%5 == 0 },
		             func(s store.Store) bool { return x.Get(s)%3 == 0 },
		         )).Return(),
		     )

		The same lexical scope can be used to avoid making calls multiple
		times.

			s := store.NewStore()
			x := store.Var(s, 10)
			s = s.Generate(
				store.Scope(func(s store.Store) store.Generator {
					v := x.Get(s)
					return store.If(func(s store.Store) bool {
						return v < 10
					}).Then(func(s store.Store) store.Store {
						s.Apply(x.Set(v + 1))
						return s
					})
				}),
			)
	*/
	Generate(...Generator) Store

	/*
		Propagate changes into a Store until it reaches a fixed point.

				s := store.NewStore()
				x := store.Var(s, 1)
				s = s.Propagate(
					func(s store.Store) []store.Change {
						return []store.Change{
							x.Set(2),
							x.Set(42),
						}
					},
				)
	*/
	Propagate(...Propagator) Store

	/*
		Value sets the integer value of a Store. When maximizing or minimizing,
		this is the value that is optimized.

			s := store.NewStore()
			x := store.Var(s, 1)
			s = s.Value(func(s store.Store) int {
				v := x.Get(s)
				return v * v
			})
	*/
	Value(Valuer) Store

	// Maximizer builds a solver that searches the space defined by the Store
	// to maximize a value.
	Maximizer(solve.Options) solve.Solver

	// Minimizer builds a solver that searches the space defined by the Store
	// to minimize a value.
	Minimizer(solve.Options) solve.Solver

	// Satisfier builds a solver that searches the space defined by the Store
	// to satisfy operational validity
	Satisfier(solve.Options) solve.Solver
}

// Generator is an interface to generate a new Store from an existing one. A
// Generator follows this logic: `If` a generating condition holds, `Then`
// generate a new Store `With` a condition for operational validity. A Store is
// operationally valid if all decisions have been made and those decisions
// fulfill certain requirements; e.g., all stops have been assigned to
// vehicles, all shifts are covered with the necessary personnel, all
// assignment have been made, quantity respects an alloted capacity, etc.
// Setting operational validity is optional and the default is true.
type Generator interface {
	Condition() Condition
	Generate() func(Store) Store
	Feasible() Condition
	With(Condition) Generator
}

// Action holds the actions for generating a new Store.
type Action interface {
	Then(func(Store) Store) Generator
	Return() Generator
	Discard() Generator
}

// Condition represents a logical condition on a context.
type Condition func(Store) bool

// Change a Store.
type Change func(Store)

// Formatter maps a Store to any type with a JSON representation.
type Formatter func(Store) any

// Bounder maps a context to monotonically tightening bounds.
type Bounder func(Store) model.Bounds

// Propagator propagates Changes to a Store.
type Propagator func(Store) []Change

// Valuer maps a Store to an integer value.
type Valuer func(Store) int

// NewStore returns a new Store.
func NewStore() Store {
	connect()
	return newStoreFunc()
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

	s := store.NewStore()
	x := store.Var(s, 1)
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
Scope specifies a Generator that allows the use of the same lexical scope.

	s := store.NewStore()
	x := store.Var(s, 1)
	s = s.Generate(
		store.Scope(func(s store.Store) store.Generator {
			v := x.Get(s)
			return store.If(func(s store.Store) bool {
				// v is used here.
				return v < 10
				}).Then(func(s store.Store) store.Store {
					// v is also used here.
					v++
					s.Apply(x.Set(v))
					return s
				})
			}),
		)
*/
func Scope(f func(Store) Generator) Generator {
	return scopeFunc(f)
}

var (
	newStoreFunc func() Store
	andFunc      func(Condition, Condition, ...Condition) Condition
	falseFunc    func(Store) bool
	notFunc      func(Condition) Condition
	orFunc       func(Condition, Condition, ...Condition) Condition
	trueFunc     func(Store) bool
	xorFunc      func(Condition, Condition) Condition
	ifFunc       func(Condition) Action
	scopeFunc    func(func(Store) Generator) Generator
)
