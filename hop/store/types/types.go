// Package types holds type definitions.
package types

/*
Store represents a store of variables and logic to solve decision automation
problems. Adding logic to the store updates it (functions may be called
directly and chained):

    s := store.New()    // s := store.New().
    s = s.Apply(...)    // 	   Apply(nil).
    s = s.Bound(...)    // 	   Generate(nil).
    s = s.Format(...)   // 	   Bound(nil).
    s = s.Generate(...) // 	   Format(nil)

The variables and logic stored define a solution space and Hop searches this
space to make decisions.
*/
type Store interface {
	/*
		Apply changes to a Store. A change happens when a stored variable is
		updated:

			s := store.New()
			x := store.Var(s, 3.1416)
			s = s.Apply(
				x.Set(x.Get(s) * 2),
			)
	*/
	Apply(...Change) Store

	/*
		Bound the value of a Store. The lower and upper bounds can be set:

				s := store.New()
				x := store.Var(s, 10)
				s = s.Bound(func(s types.Store) types.Bounds {
					return types.Bounds{
						Lower: x.Get(s),
						Upper: math.MaxInt,
					}
				})
	*/
	Bound(Bounder) Store

	/*
		Format a Store into any structure prior to JSON encoding.

				s := store.New()
				x := store.Var(s, 10)
				s = s.Format(func(s types.Store) any {
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

		     s := store.New()
		     x := store.Var(s, 10)
		     s = s.Generate(
		         store.If(func(s types.Store) bool {
		             // If the value of x is smaller than 10.
		             return x.Get(s) < 10
		         }).Then(func(s types.Store) types.Store {
		             // Then increase the value of x by 1.
		             v := x.Get(s)
		             v++
		             s.Apply(x.Set(v))
		             return s
		         }).With(func(s types.Store) bool {
		             // Operationally valid if x is even.
		             return x.Get(s)%2 == 0
		         }),
		     )

		There are convenience exit functions to stop generating new Stores
		and set their operational validity.

		     s := store.New()
		     x := store.Var(s, 10)
		     s = s.Generate(
		         // If x is greater than 20 stop the search and make it
				 // operationally invalid.
		         store.If(func(s types.Store) bool { return x.Get(s) > 20 }).
		             Discard(),
		         // If x is a multiple of 5 and 3 stop the search and make it
				 // operationally valid.
		         store.If(store.And(
		             func(s types.Store) bool { return x.Get(s)%5 == 0 },
		             func(s types.Store) bool { return x.Get(s)%3 == 0 },
		         )).Return(),
		     )

		The same lexical scope can be used to avoid making calls multiple
		times.

			s := store.New()
			x := store.Var(s, 10)
			s = s.Generate(
				store.Scope(func(s types.Store) types.Generator {
					v := x.Get(s)
					return store.If(func(s types.Store) bool {
						return v < 10
					}).Then(func(s types.Store) types.Store {
						s.Apply(x.Set(v + 1))
						return s
					})
				}),
			)
	*/
	Generate(...Generator) Store

	/*
		Propagate changes into a Store until it reaches a fixed point.

				s := store.New()
				x := store.Var(s, 1)
				s = s.Propagate(
					func(s types.Store) []types.Change {
						return []types.Change{
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

			s := store.New()
			x := store.Var(s, 1)
			s = s.Value(func(s types.Store) int {
				v := x.Get(s)
				return v * v
			})
	*/
	Value(Valuer) Store

	// Maximizer builds a solver that searches the space defined by the Store
	// to maximize a value.
	Maximizer(Options) Solver

	// Minimizer builds a solver that searches the space defined by the Store
	// to minimize a value.
	Minimizer(Options) Solver

	// Satisfier builds a solver that searches the space defined by the Store
	// to satisfy operational validity
	Satisfier(Options) Solver
}

// A Generator is used to generate new Stores from an existing one. It follows
// this logic: `If` a generating condition holds, `Then` generate a new Store
// `With` a condition for operational validity. A Store is operationally valid
// if all decisions have been made and those decisions fulfill certain
// requirements; e.g., all stops have been assigned to vehicles, all shifts are
// covered with the necessary personnel, all assignment have been made,
// quantity respects an alloted capacity, etc. Setting operational validity is
// optional and the default is true.
type Generator interface {
	// Condition returns the generating Condition. The Generator may generate
	// new Stores as long as this Condition holds.
	Condition() Condition

	// Generate returns the generating function. The function takes in a Store
	// and transforms it to generate a new one.
	Generate() func(Store) Store

	// Valid returns the Condition that determines operational validity.
	Valid() Condition

	/*
		With establishes the Condition for operational validity that is used
		with a generated Store.

			s := store.New()
			x := store.Var(s, 1)
			s = s.Generate(
				store.If(store.True).
					Then(func(s types.Store) types.Store {
						return s.Apply(x.Set(x.Get(s) + 1))
					}).
					With(func(s types.Store) bool {
						// The Store is operationally valid if x is even.
						return x.Get(s)%2 == 0
					}),
			)
	*/
	With(Condition) Generator
}

// Action holds the actions for generating a new Store.
type Action interface {
	/*
		Then receives a function that takes in a Store and transforms it to
		generate a new one.

			s := store.New()
			x := store.Var(s, 1)
			s = s.Generate(
				store.If(store.True).
					Then(func(s types.Store) types.Store {
		                // Generate a new store in which the value of x
						// increases by 1.
						return s.Apply(x.Set(x.Get(s) + 1))
					}),
			)
	*/
	Then(func(Store) Store) Generator

	/*
		Return is a convenience function that returns the existing Store,
		setting it as operationally valid.

			// Return an operationally valid Store.
			s := store.New().
				Generate(store.If(store.True).Return())
	*/
	Return() Generator

	/*
		Discard is a convenience function that discards the existing Store,
		setting it as operationally invalid.

			// Discard an operationally invalid Store.
			s := store.New().
				Generate(store.If(store.True).Discard())
	*/
	Discard() Generator
}

// Condition represents a logical condition on a context.
type Condition func(Store) bool

// Change a Store.
type Change func(Store)

// Formatter maps a Store to any type with a JSON representation.
type Formatter func(Store) any

// Bounder maps a context to monotonically tightening bounds.
type Bounder func(Store) Bounds

// Propagator propagates Changes to a Store.
type Propagator func(Store) []Change

// Valuer maps a Store to an integer value.
type Valuer func(Store) int

// Variable that can be stored in a Store.
type Variable[T any] interface {
	/*
		Get gets the current value of the Variable in the Store.

			s := store.New()
			x := store.Var(s, 10)
			s = s.Format(func(s types.Store) any {
				return map[string]int{"x": x.Get(s)}
			})

	*/
	Get(Store) T

	/*
		Set sets a new value on the Variable.

			s := store.New()
			x := store.Var(s, 10)
			s = s.Apply(x.Set(15))
	*/
	Set(T) Change
}

// Slice manages an immutable slice container of some type in a Store.
type Slice[T any] interface {
	// Append one or more values to the end of a slice.
	Append(value T, values ...T) Change
	// Get an index of a slice.
	Get(Store, int) T
	// Insert one or more values at an index in a slice.
	Insert(index int, value T, values ...T) Change
	// Len returns the length of a slice.
	Len(Store) int
	// Prepend one or more values at the beginning of a slice.
	Prepend(value T, values ...T) Change
	// Remove a subslice from a start to an end index.
	Remove(start, end int) Change
	// Set a value by index.
	Set(int, T) Change
	// Slice representation that is mutable.
	Slice(Store) []T
}

// A Key for a Map.
type Key interface{ int | string }

// A Map stores key-value pairs in a Store.
type Map[K Key, V any] interface {
	// Delete a key from the map.
	Delete(K) Change
	// Get an index of a vector.
	Get(Store, K) (V, bool)
	// Len returns the number of keys in a map,
	Len(Store) int
	// Map representation that is mutable.
	Map(Store) map[K]V
	// Set a key to a value.
	Set(K, V) Change
}

// Options for a solver.
type Options any

// A Solver searches a space.
type Solver any

// Bounds on an objective value at some node in the search tree consist of a
// lower value and an upper value. If the lower and upper value are the same,
// the bounds have converged.
type Bounds struct {
	Lower int
	Upper int
}
