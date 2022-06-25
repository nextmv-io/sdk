package store

import (
	"context"
	"time"

	"github.com/nextmv-io/sdk/model"
)

/*
Store represents a store of Variables and logic to solve decision automation
problems. Adding logic to the Store updates it (functions may be called
directly and chained):

    s := store.New()    // s := store.New().
    s = s.Apply(...)    // 	   Apply(...).
    s = s.Bound(...)    // 	   Generate(...).
    s = s.Format(...)   // 	   Bound(...).
    s = s.Generate(...) // 	   Format(...)

The Variables and logic stored define a solution space. This space is searched
to make decisions.
*/
type Store interface {
	/*
		Apply changes to a Store. A change happens when a stored Variable is
		updated:

			s := store.New()
			x := store.NewVar(s, 3.1416)
			s = s.Apply(
				x.Set(x.Get(s) * 2),
			)
	*/
	Apply(...Change) Store

	/*
		Bound the value of a Store. The lower and upper bounds can be set:

				s := store.New()
				x := store.NewVar(s, 10)
				s = s.Bound(func(s store.Store) store.Bounds {
					return store.Bounds{
						Lower: x.Get(s),
						Upper: math.MaxInt,
					}
				})
	*/
	Bound(Bounder) Store

	/*
		Format a Store into any structure prior to JSON encoding.

				s := store.New()
				x := store.NewVar(s, 10)
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

		     s := store.New()
		     x := store.NewVar(s, 10)
		     s = s.Generate(
		         store.If(func(s store.Store) bool {
		             // If the value of x is smaller than 10.
		             return x.Get(s) < 10
		         }).Then(func(s store.Store) store.Store {
		             // Then increase the value of x by 1.
		             v := x.Get(s)
		             v++
		             return s.Apply(x.Set(v))
		         }).With(func(s store.Store) bool {
		             // Operationally valid if x is even.
		             return x.Get(s)%2 == 0
		         }),
		     )

		There are convenience exit functions to stop generating new Stores
		and set their operational validity.

		     s := store.New()
		     x := store.NewVar(s, 10)
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

			s := store.New()
			x := store.NewVar(s, 10)
			s = s.Generate(
				store.Scope(func(s store.Store) store.Generator {
					v := x.Get(s)
					return store.If(func(s store.Store) bool {
						return v < 10
					}).Then(func(s store.Store) store.Store {
						return s.Apply(x.Set(v + 1))
					})
				}),
			)
	*/
	Generate(...Generator) Store

	/*
		Propagate changes into a Store until it reaches a fixed point.

				s := store.New()
				x := store.NewVar(s, 1)
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

			s := store.New()
			x := store.NewVar(s, 1)
			s = s.Value(func(s store.Store) int {
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
	// to satisfy operational validity.
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
			x := store.NewVar(s, 1)
			s = s.Generate(
				store.If(store.True).
					Then(func(s store.Store) store.Store {
						return s.Apply(x.Set(x.Get(s) + 1))
					}).
					With(func(s store.Store) bool {
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
			x := store.NewVar(s, 1)
			s = s.Generate(
				store.If(store.True).
					Then(func(s store.Store) store.Store {
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

// Variable stored in a Store.
type Variable[T any] interface {
	/*
		Get the current value of the Variable in the Store.

			s := store.New()
			x := store.NewVar(s, 10)
			s = s.Format(func(s store.Store) any {
				return map[string]int{"x": x.Get(s)}
			})

	*/
	Get(Store) T

	/*
		Set a new value on the Variable.

			s := store.New()
			x := store.NewVar(s, 10)
			s = s.Apply(x.Set(15))
	*/
	Set(T) Change
}

// Slice manages an immutable slice container of some type in a Store.
type Slice[T any] interface {
	/*
		Append one or more values to the end of a Slice.

			s1 := store.New()
			x := store.NewSlice(s, 1, 2, 3)   // [1, 2, 3]
			s2 := s1.Apply(x.Append(4, 5)) // [1, 2, 3, 4, 5]
	*/
	Append(value T, values ...T) Change

	/*
		Get an index of a Slice.

			s := store.New()
			x := store.NewSlice(s, 1, 2, 3)
			x.Get(s, 2) // 3
	*/
	Get(Store, int) T

	/*
		Insert one or more values at an index in a Slice.

			s1 := store.New()
			x := store.NewSlice(s, "a", "b", "c")
			s2 := s1.Apply(s.Insert(2, "d", "e")) // [a, b, d, e, c]
	*/
	Insert(index int, value T, values ...T) Change

	/*
		Len returns the length of a Slice.

			s := store.New()
			x := store.NewSlice(s, 1, 2, 3)
			x.Len(s) // 3
	*/
	Len(Store) int

	/*
		Prepend one or more values at the beginning of a Slice.

			s1 := store.New()
			x := store.NewSlice(s, 1, 2, 3)    // [1, 2, 3]
			s2 := s1.Apply(x.Prepend(4, 5)) // [4, 5, 1, 2, 3]
	*/
	Prepend(value T, values ...T) Change

	/*
		Remove a sub-Slice from a starting to an ending index.

			s1 := store.New()
			x := store.NewSlice(s, 1, 2, 3) // [1, 2, 3]
			s2 := s1.Apply(x.Remove(1))  // [1, 3]
	*/
	Remove(start, end int) Change

	/*
		Set a value by index.

			s1 := store.New()
			x := store.NewSlice(s, "a", "b", "c") // [a, b, c]
			s2 := s1.Apply(x.Set(1, "d"))      // [a, d, c]
	*/
	Set(int, T) Change

	/*
		Slice representation that is mutable.

			s := store.New()
			x := store.NewSlice(s, 1, 2, 3)
			x.Slice(s) // []int{1, 2, 3}
	*/
	Slice(Store) []T
}

// A Key for a Map.
type Key interface{ int | string }

// A Map stores key-value pairs in a Store.
type Map[K Key, V any] interface {
	/*
		Delete a Key from the Map.

			s1 := store.New()
			m := store.NewMap[int, string](s1)
			s1 = s1.Apply( // {42: foo, 13: bar}
				m.Set(42, "foo"),
				m.Set(13, "bar"),
			)
			s2 := s1.Apply(m.Delete("foo")) // {13: bar}
	*/
	Delete(K) Change

	/*
		Get a value for a Key.

			s1 := store.New()
			m := store.NewMap[int, string](s1)
			s2 := s1.Apply(m.Set(42, "foo"))
			m.Get(s2) // (foo, true)
			m.Get(s2) // (_, false)
	*/
	Get(Store, K) (V, bool)

	/*
		Len returns the number of Keys in a Map.

			s1 := store.New()
			m := store.NewMap[int, string](s1)
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
			m := store.NewMap[int, string](s1)
			s2 := s1.Apply(
				m.Set(42, "foo"),
				m.Set(13, "bar"),
			)
			m.Map(s2) // map[int]string{42: "foo", 13: "bar"}
	*/
	Map(Store) map[K]V

	/*
		Set a Key to a Value.

			s1 := store.New()
			m := store.NewMap[int, string](s1)
			s2 := s1.Apply(m.Set(42, "foo")) // 42 -> foo
			s3 := s2.Apply(m.Set(42, "bar")) // 42 -> bar
	*/
	Set(K, V) Change
}

// A Domain of integers.
type Domain interface {
	/*
		Add values to a Domain.

			s1 := store.New()
			d := store.Multiple(s, 1, 3, 5)
			s2 := s1.Apply(d.Add(2, 4))

			d.Domain(s1) // {1, 3, 5}}
			d.Domain(s2) // [1, 5]]
	*/
	Add(...int) Change

	/*
		AtLeast updates the Domain to the sub-Domain of at least some value.

			s1 := store.New()
			d := store.NewDomain(s, model.NewRange(1, 10), model.NewRange(101, 110))
			s2 := s1.Apply(d.AtLeast(50))

			d.Domain(s1) // {[1, 10], [101, 110]}
			d.Domain(s2) // [101, 110]
	*/
	AtLeast(int) Change

	/*
		AtMost updates the Domain to the sub-Domain of at most some value.

			s1 := store.New()
			d := store.NewDomain(s, model.NewRange(1, 10), model.NewRange(101, 110))
			s2 := s1.Apply(d.AtMost(50))

			d.Domain(s1) // {[1, 10], [101, 110]}
			d.Domain(s2) // [1, 10]
	*/
	AtMost(int) Change

	/*
		Cmp lexically compares two integer Domains. It returns a negative value
		if the receiver is less, 0 if they are equal, and a positive value if
		the receiver Domain is greater.

			s := store.New()
			d1 := store.NewDomain(s, model.NewRange(1, 5), model.NewRange(8, 10))
			d2 := store.Multiple(s, -1, 1)
			d1.Cmp(s, d2) // > 0
	*/
	Cmp(Store, Domain) int

	/*
		Contains returns true if a Domain contains a given value.

			s := store.New()
			d := store.NewDomain(s, model.NewRange(1, 10))
			d.Contains(s, 5)  // true
			d.Contains(s, 15) // false
	*/
	Contains(Store, int) bool

	/*
		Domain returns a Domain unattached to a Store.

			s := store.New()
			d := store.NewDomain(s, model.NewRange(1, 10))
			d.Domain(s) // model.NewDomain(model.NewRange(1, 10))
	*/
	Domain(Store) model.Domain

	/*
		Empty is true if a Domain is empty for a Store.

			s := store.New()
			d1 := store.NewDomain(s)
			d2 := store.Singleton(s, 42)
			d1.Empty() // true
			d2.Empty() // false
	*/
	Empty(Store) bool

	/*
		Len of a Domain, counting all values within ranges.

			s := store.New()
			d := store.NewDomain(s, model.NewRange(1, 10), model.NewRange(-5, -1))
			d.Len(s) // 15
	*/
	Len(Store) int

	/*
		Max of a Domain and a boolean indicating it is non-empty.

			s := store.New()
			d1 := store.NewDomain(s)
			d2 := store.NewDomain(s, model.NewRange(1, 10), model.NewRange(-5, -1))
			d1.Max() // returns (_, false)
			d2.Max() // returns (10, true)
	*/
	Max(Store) (int, bool)

	/*
		Min of a Domain and a boolean indicating it is non-empty.

			s := store.New()
			d1 := store.NewDomain(s)
			d2 := store.NewDomain(s, model.NewRange(1, 10), model.NewRange(-5, -1))
			d1.Min() // returns (_, false)
			d2.Min() // returns (-5, true)

	*/
	Min(Store) (int, bool)

	/*
		Remove values from a Domain.

			s1 := store.New()
			d := store.NewDomain(s, model.NewRange(1, 5))
			s2 := s1.Apply(d.Remove(2, 4))

			d.Domain(s1) // [1, 5]
			d.Domain(s2) // {1, 3, 5}
	*/
	Remove(...int) Change

	/*
		Slice representation of a Domain.

			s := store.New()
			d := store.NewDomain(s, model.NewRange(1, 5))
			d.Slice(s) // [1, 2, 3, 4, 5]
	*/
	Slice(Store) []int

	/*
		Value returns an int and true if a Domain is Singleton.

			s := store.New()
			d1 := store.NewDomain(s)
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
		Add values to a Domain by index.

			s1 := store.New()
			d := store.Repeat(s1, 1, model.Singleton(42)) // [42, 42, 42]
			s2 := s1.Apply(d.Add(1, 41, 43))              // [42, [41,43], 42]
	*/
	Add(int, ...int) Change

	/*
		Assign a Singleton value to a Domain by index.

			s1 := store.New()
			d := store.Repeat(s1, 3, model.Singleton(42)) // [42, 42, 42]
			s2 := s1.Apply(d.Assign(0, 10))               // [10, 42, 42]
	*/
	Assign(int, int) Change

	/*
		AtLeast updates the Domain to the sub-Domain of at least some value.

			s1 := store.New()
			d := store.Repeat( // [[1, 100], [1, 100]]
				s1,
				2,
				model.NewDomain(model.NewRange(1, 100)),
			)
			s2 := s1.Apply(d.AtLeast(1, 50)) // [[1, 100], [50, 100]]
	*/
	AtLeast(int, int) Change

	/*
		AtMost updates the Domain to the sub-Domain of at most some value.

			s1 := store.New()
			d := store.Repeat( // [[1, 100], [1, 100]]
				s1,
				2,
				model.NewDomain(model.NewRange(1, 100)),
			)
			s2 := s1.Apply(d.AtMost(1, 50)) // [[1, 100], [1, 50]]
	*/
	AtMost(int, int) Change

	/*
		Cmp lexically compares two sequences of integer Domains.

			s := store.New()
			d1 := store.Repeat(s, 3, model.Singleton(42)) // [42, 42, 42]
			d2 := store.Repeat(s, 2, model.Singleton(43)) // [43, 43]]
			d1.Cmp(s, d2) // < 0
	*/
	Cmp(Store, Domains) int

	/*
		Domain by index.

			s := store.New()
			d := store.NewDomains(
				s,
				model.NewDomain(),
				model.Singleton(42),
			)
			d.Domain(s, 0) // {}
			d.Domain(s, 1) // 42
	*/
	Domain(Store, int) model.Domain

	/*
		Domains in the sequence.

			s := store.New()
			d := store.NewDomains(
				s,
				model.NewDomain(),
				model.Singleton(42),
			)
			d.Domains(s) // [{}, 42}
	*/
	Domains(Store) model.Domains

	/*
		Empty is true if all Domains are empty.

			s := store.New()
			d := store.NewDomains(s, model.NewDomain())
			d.Empty(s) // true
	*/
	Empty(Store) bool

	/*
		Len returns the number of Domains.

			s := store.New()
			d := store.Repeat(s, 5, model.NewDomain())
			d.Len(s) // 5
	*/
	Len(Store) int

	/*
		Remove values from a Domain by index.

			s1 := store.New()
			d := store.NewDomains(s1, model.Multiple(42, 13)) // {13, 42}
			s2 := s1.Apply(d.Remove(13)) // {42}
	*/
	Remove(int, ...int) Change

	/*
		Singleton is true if all Domains are Singleton.

			s := store.New()
			d := store.Repeat(s, 5, model.Singleton(42))
			d.Singleton(s) // true
	*/
	Singleton(Store) bool

	/*
		Slices converts Domains to a slice of int slices.

			s := store.New()
			d := store.NewDomains(s, model.NewDomain(), model.Multiple(1, 3))
			d.Slices(s) // [[], [1, 2, 3]]
	*/
	Slices(Store) [][]int

	/*
		Values returns the values of a sequence of Singleton Domains.

			s1 := store.New()
			d := store.Repeat(s1, 3, model.Singleton(42))
			s2 := store.Apply(d.Add(0, 41))
			d.Values(s1) // ([42, 42, 42], true)
			d.Values(s2) // (_, false)
	*/
	Values(Store) ([]int, bool)

	/* Domain selectors */

	// First returns the first Domain index with length above 1.
	First(Store) (int, bool)
	// Largest returns the index of the largest Domain with length above 1 by
	// number of elements.
	Largest(Store) (int, bool)
	// Last returns the last Domain index with length above 1.
	Last(Store) (int, bool)
	// Maximum returns the index of the Domain containing the maximum value
	// with length above 1.
	Maximum(Store) (int, bool)
	// Minimum returns the index of the Domain containing the minimum value
	// with length above 1.
	Minimum(Store) (int, bool)
	// Smallest returns the index of the smallest Domain with length above 1 by
	// number of elements.
	Smallest(Store) (int, bool)
}

// Sense specifies whether one is maximizing, minimizing, or satisfying.
type Sense int

// Options for a solver.
type Options struct {
	Sense Sense
	// Tags are custom key-value pairs that the user defines for
	// record-keeping.
	Tags    map[string]any
	Diagram Diagram
	// Search options.
	Search struct {
		// Buffer represents the maximum number of Stores that can be buffered
		// when generating more Stores.
		Buffer int
	}
	Limits Limits
	// Options for random number generation.
	Random struct {
		// Seed for generating random numbers.
		Seed int64 `json:"seed,omitempty"`
	}
	// Pool that is used in specific engines.
	Pool struct {
		// Maximum Size of the Pool.
		Size int `json:"size,omitempty"`
	}
}

// Diagram options. The Store search is based on Decision Diagrams. These
// options configure the mechanics of using DD.
type Diagram struct {
	// Maximum Width of the Decision Diagram.
	Width int
	// Maximum Expansion that can be generated from a Store.
	Expansion struct {
		// Limit represents the maximum number of children Stores that can
		// be generated from a parent.
		Limit int `json:"limit"`
	}
}

// Limits when performing a search. The search will stop if any one of these
// limits are encountered.
type Limits struct {
	// Time Duration.
	Duration time.Duration
	// Nodes reprent active Stores in the search.
	Nodes int
	// Solutions represent operationally valid Stores.
	Solutions int
}

// A Solver searches a space and finds the best Solution possible, this is, the
// best collection of Variable assignments in an operationally valid Store.
type Solver interface {
	// All Solutions found by the Solver.
	All(context.Context) []Solution

	/*
		Last Solution found by the Solver. When running a Maximizer or
		Minimizer, the last Solution is the best one found (highest or smallest
		value, respectively) with the given options. Using this function is
		equivalent to getting the last element when using All:

		    s := store.New()
		    x := store.NewVar(s, 0)
		    opt := store.DefaultOptions()
		    // Minimizer and Satisfier may also be used
		    solver = s.Generate(...).Value(...).Format(...).Maximizer(opt)
		    all := solver.All(context.Background())
		    last := all[len(all)-1]

	*/
	Last(context.Context) Solution

	// Options provided to the Solver.
	Options() Options
}

// Solution of a decision automation problem. A Solution is an operationally
// valid Store.
type Solution struct {
	// Store of the Solution. If nil, it means that the solution is
	// operationally invalid.
	Store      Store      `json:"store"`
	Statistics Statistics `json:"statistics"`
}

// Statistics of the search.
type Statistics struct {
	// Bounds of the store. Nil when using a Satisfier.
	Bounds *Bounds `json:"bounds,omitempty"`
	Search Search  `json:"search"`
	Time   Time    `json:"time"`
	// Value of the store. Nil when using a Satisfier.
	Value *int `json:"value,omitempty"`
}

// Search statistics of the Store generation.
type Search struct {
	// Generated stores in the search.
	Generated int `json:"generated"`
	// Filtered stores in the search.
	Filtered int `json:"filtered"`
	// Expanded stores in the search.
	Expanded int `json:"expanded"`
	// Reduced stores in the search.
	Reduced int `json:"reduced"`
	// Restricted stores in the search.
	Restricted int `json:"restricted"`
	// Deferred stores in the search.
	Deferred int `json:"deferred"`
	// Explored stores in the search.
	Explored int `json:"explored"`
	// Operationally valid stores in the search.
	Solutions int `json:"solutions"`
}

// Time statistics.
type Time struct {
	Start   time.Time     `json:"start"`
	Elapsed time.Duration `json:"elapsed"`
}

// Bounds on an objective value at some node in the search tree consist of a
// lower value and an upper value. If the lower and upper value are the same,
// the bounds have converged.
type Bounds struct {
	Lower int `json:"lower"`
	Upper int `json:"upper"`
}
