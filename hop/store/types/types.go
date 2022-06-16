// Package types holds type definitions.
package types

import "github.com/nextmv-io/sdk/hop/model/types"

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
	/*
		Append one or more values to the end of a slice.

			s1 := store.New()
			x := store.Slice(s, 1, 2, 3)   // [1, 2, 3]
			s2 := s1.Apply(x.Append(4, 5)) // [1, 2, 3, 4, 5]
	*/
	Append(value T, values ...T) Change

	/*
		Get an index of a slice.

			s := store.New()
			x := store.Slice(s, 1, 2, 3)
			x.Get(s, 2) // 3
	*/
	Get(Store, int) T

	/*
		Insert one or more values at an index in a slice.

			s1 := store.New()
			x := store.Slice(s, "a", "b", "c")
			s2 := s1.Apply(s.Insert(2, "d", "e")) // [a, b, d, e, c]
	*/
	Insert(index int, value T, values ...T) Change

	/*
		Len returns the length of a slice.

			s := store.New()
			x := store.Slice(s, 1, 2, 3)
			x.Len(s) // 3
	*/
	Len(Store) int

	/*
		Prepend one or more values at the beginning of a slice.

			s1 := store.New()
			x := store.Slice(s, 1, 2, 3)    // [1, 2, 3]
			s2 := s1.Apply(x.Prepend(4, 5)) // [4, 5, 1, 2, 3]
	*/
	Prepend(value T, values ...T) Change

	/*
		Remove a subslice from a start to an end index.

			s1 := store.New()
			x := store.Slice(s, 1, 2, 3) // [1, 2, 3]
			s2 := s1.Apply(x.Remove(1))  // [1, 3]
	*/
	Remove(start, end int) Change

	/*
		Set a value by index.

			s1 := store.New()
			x := store.Slice(s, "a", "b", "c") // [a, b, c]
			s2 := s1.Apply(x.Set(1, "d"))      // [a, d, c]
	*/
	Set(int, T) Change

	/*
		Slice representation that is mutable.

			s := store.New()
			x := store.Slice(s, 1, 2, 3)
			x.Slice(s) // []int{1, 2, 3}
	*/
	Slice(Store) []T
}

// A Key for a Map.
type Key interface{ int | string }

// A Map stores key-value pairs in a Store.
type Map[K Key, V any] interface {
	/*
		Delete a key from the map.

			s1 := store.New()
			m := store.Map[int, string](s1)
			s1 = s1.Apply( // {42: foo, 13: bar}
				m.Set(42, "foo"),
				m.Set(13, "bar"),
			)
			s2 := s1.Apply(m, Delete("foo")) // {13: bar}
	*/
	Delete(K) Change

	/*
		Get a value for a key.

			s1 := store.New()
			m := store.Map[int, string](s1)
			s2 := s1.Apply(m.Set(42, "foo"))
			m.Get(s2) // (foo, true)
			m.Get(s2) // (_, false)
	*/
	Get(Store, K) (V, bool)

	/*
		Len returns the number of keys in a map,

			s1 := store.New()
			m := store.Map[int, string](s1)
			s2 := s1.Apply(
				m.Set(42, "foo"),
				m.Set(13, "bar"),
			)
			m.Len(s1) // 0
			m.Len(s2) // 2
	*/
	Len(Store) int

	/*
		Map representation that is mutable.

			s1 := store.New()
			m := store.Map[int, string](s1)
			s2 := s1.Apply(
				m.Set(42, "foo"),
				m.Set(13, "bar"),
			)
			m.Map(s2) // map[int]string{42: "foo", 13: "bar"}
	*/
	Map(Store) map[K]V

	/*
		Set a key to a value.

			s1 := store.New()
			m := store.Map[int, string](s1)
			s2 := s1.Apply(m.Set(42, "foo")) // 42 -> foo
			s3 := s2.Apply(m.Set(42, "bar")) // 42 -> bar
	*/
	Set(K, V) Change
}

// A Domain of integers.
type Domain interface {
	/*
		Add values to a domain.

			s1 := store.New()
			d := store.Multiple(s, 1, 3, 5)
			s2 := s1.Apply(d.Add(2, 4))

			d.Domain(s1) // {1, 3, 5}}
			d.Domain(s2) // [1, 5]]
	*/
	Add(...int) Change

	/*
		AtLeast updates the domain to the subdomain of at least some value.

			s1 := store.New()
			d := store.Domain(s, model.Range(1, 10), model.Range(101, 110))
			s2 := s1.Apply(d.AtLeast(50))

			d.Domain(s1) // {[1, 10], [101, 110]}
			d.Domain(s2) // [101, 110]
	*/
	AtLeast(int) Change

	/*
		AtMost updates the domain to the subdomain of at most some value.

			s1 := store.New()
			d := store.Domain(s, model.Range(1, 10), model.Range(101, 110))
			s2 := s1.Apply(d.AtMost(50))

			d.Domain(s1) // {[1, 10], [101, 110]}
			d.Domain(s2) // [1, 10]
	*/
	AtMost(int) Change

	/*
		Cmp lexically compares two integer domains. It returns a negative value
		if the receiver is less, 0 if they are equal, and a positive value if
		the receiver domain is greater.

			s := store.New()
			d1 := store.Domain(s, model.Range(1, 5), model.Range(8, 10))
			d2 := store.Multiple(s, -1, 1)
			d1.Cmp(s, d2) // > 0
	*/
	Cmp(Store, Domain) int

	/*
		Contains returns true if a domain contains a given value.

			s := store.New()
			d := store.Domain(s, model.Range(1, 10))
			d.Contains(s, 5)  // true
			d.Contains(s, 15) // false
	*/
	Contains(Store, int) bool

	/*
		Domain returns a domain unattached to a store.

			s := store.New()
			d := store.Domain(s, model.Range(1, 10))
			d.Domain(s) // model.Domain(model.Range(1, 10))
	*/
	Domain(Store) types.Domain

	/*
		Empty is true if a domain is empty for a store.

			s := store.New()
			d1 := store.Domain(s)
			d2 := store.Singleton(s, 42)
			d1.Empty() // true
			d2.Empty() // false
	*/
	Empty(Store) bool

	/*
		Len of a domain, counting all values within ranges.

			s := store.New()
			d := store.Domain(s, model.Range(1, 10), model.Range(-5, -1))
			d.Len(s) // 15
	*/
	Len(Store) int

	/*
		Max of a domain and a boolean indicating it is nonempty.

			s := store.New()
			d1 := store.Domain(s)
			d2 := store.Domain(s, model.Range(1, 10), model.Range(-5, -1))
			d1.Max() // returns (_, false)
			d2.Max() // returns (10, true)
	*/
	Max(Store) (int, bool)

	/*
		Min of a domain and a boolean indicating it is nonempty.

			s := store.New()
			d1 := store.Domain(s)
			d2 := store.Domain(s, model.Range(1, 10), model.Range(-5, -1))
			d1.Min() // returns (_, false)
			d2.Min() // returns (-5, true)

	*/
	Min(Store) (int, bool)

	/*
		Remove values from a domain.

			s1 := store.New()
			d := store.Domain(s, model.Range(1, 5))
			s2 := s1.Apply(d.Remove(2, 4))

			d.Domain(s1) // [1, 5]
			d.Domain(s2) // {1, 3, 5}
	*/
	Remove(...int) Change

	/*
		Slice representation of a domain.

			s := store.New()
			d := store.Domain(s, model.Range(1, 5))
			d.Slice(s) // [1, 2, 3, 4, 5]
	*/
	Slice(Store) []int

	/*
		Value returns an int and true if a domain is singleton.

			s := store.New()
			d1 := store.Domain(s)
			d2 := store.Singleton(s, 42)
			d3 := store.Multiple(s, 1, 3, 5)
			d1.Value() // returns (_, false)
			d2.Value() // returns (42, true)
			d3.Value() // returns (_, false)
	*/
	Value(Store) (int, bool)
}

// Domains of integers.
type Domains interface {
	/*
		Add values to a domain by index.

			s1 := store.New()
			d := store.Repeat(s1, 1, model.Singleton(42)) // [42, 42, 42]
			s2 := s1.Apply(d.Add(1, 41, 43))              // [42, [41,43], 42]
	*/
	Add(int, ...int) Change

	/*
		Assign a singleton value to a domain by index.

			s1 := store.New()
			d := store.Repeat(s1, 3, model.Singleton(42)) // [42, 42, 42]
			s2 := s1.Apply(d.Assign(0, 10))               // [10, 42, 42]
	*/
	Assign(int, int) Change

	/*
		AtLeast updates the domain to the subdomain of at least some value.

			s1 := store.New()
			d := store.Repeat( // [[1, 100], [1, 100]]
				s1,
				2,
				model.Domain(model.Range(1, 100)),
			)
			s2 := s1.Apply(d.AtLeast(1, 50)) // [[1, 100], [50, 100]]
	*/
	AtLeast(int, int) Change

	/*
		AtMost updates the domain to the subdomain of at most some value.

			s1 := store.New()
			d := store.Repeat( // [[1, 100], [1, 100]]
				s1,
				2,
				model.Domain(model.Range(1, 100)),
			)
			s2 := s1.Apply(d.AtMost(1, 50)) // [[1, 100], [1, 50]]
	*/
	AtMost(int, int) Change

	/*
		Cmp lexically compares two sequences of integer domains.

			s := store.New()
			d1 := store.Repeat(s, 3, model.Singleton(42)) // [42, 42, 42]
			d2 := store.Repeat(s, 2, model.Singleton(43)) // [43, 43]]
			d1.Cmp(s, d2) // < 0
	*/
	Cmp(Store, Domains) int

	/*
		Domain by index.

			s := store.New()
			d := store.Domains(
				s,
				model.Domain(),
				model.Singleton(42),
			)
			d.Domain(s, 0) // {}
			d.Domain(s, 1) // 42
	*/
	Domain(Store, int) types.Domain

	/*
		Domains in the sequence.

			s := store.New()
			d := store.Domains(
				s,
				model.Domain(),
				model.Singleton(42),
			)
			d.Domains(s) // [{}, 42}
	*/
	Domains(Store) types.Domains

	/*
		Empty is true if all domains are empty.

			s := store.New()
			d := store.Domains(s, model.Domain())
			d.Empty(s) // true
	*/
	Empty(Store) bool

	/*
		Len returns the number of domains.

			s := store.New()
			d := store.Repeat(s, 5, model.Domain())
			d.Len(s) // 5
	*/
	Len(Store) int

	/*
		Remove values from a domain by index.

			s1 := store.New()
			d := store.Domains(s1, model.Multiple(42, 13)) // {13, 42}
			s2 := s1.Apply(d.Remove(13))                   // {42}
	*/
	Remove(int, ...int) Change

	/*
		Singleton is true if all domains are Singleton.

			s := store.New()
			d := store.Repeat(s, 5, model.Singleton(42))
			d.Singleton(s) // true
	*/
	Singleton(Store) bool

	/*
		Slices converts domains to a slice of int slices.

			s := store.New()
			d := store.Domains(s, model.Domain(), model.Multiple(1, 3))
			d.Slices(s) // [[], [1, 2, 3]]
	*/
	Slices(Store) [][]int

	/*
		Values returns the values of a sequence of singleton domains.

			s1 := store.New()
			d := store.Repeat(s1, 3, model.Singleton(42))
			s2 := store.Apply(d.Add(0, 41))
			d.Values(s1) // ([42, 42, 42], true)
			d.Values(s2) // (_, false)
	*/
	Values(Store) ([]int, bool)

	/* Domain selectors */

	// First returns the first domain index with length above 1.
	First(Store) (int, bool)
	// Largest returns the index of the largest domain with length above 1 by
	// number of elements.
	Largest(Store) (int, bool)
	// Last returns the last domain index with length above 1.
	Last(Store) (int, bool)
	// Maximum returns the index of the domain containing the maximum value with
	// length above 1.
	Maximum(Store) (int, bool)
	// Minimum returns the index of the domain containing the minimum value with
	// length above 1.
	Minimum(Store) (int, bool)
	// Smallest returns the index of the smallest domain with length above 1 by
	// number of elements.
	Smallest(Store) (int, bool)
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
