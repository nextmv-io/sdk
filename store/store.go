package store

import (
	"context"
	"encoding/json"
	"math"
	"time"
)

const (
	// Minimize indicates the solution space is being searched to find the
	// smallest possible value.
	Minimize Sense = iota
	// Maximize indicates the solution space is being searched to find the
	// biggest possible value.
	Maximize
	// Satisfy indicates the solution space is being searched to find
	// operationally valid Stores.
	Satisfy
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
// optional and the default is true. A user is not encouraged to implement the
// Generator interface.
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
		with a generated Store. When not used, the default operational validity
		is true.

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

// Action holds the actions for generating a new Store. They are meant to be
// used with the store.Generate function.
type Action interface {
	/*
		Discard is a convenience function that discards the existing Store,
		setting it as operationally invalid.

			// Discard an operationally invalid Store.
			s := store.New().
				Generate(store.If(store.True).Discard())
	*/
	Discard() Generator

	/*
		Return is a convenience function that returns the existing Store,
		setting it as operationally valid.

			// Return an operationally valid Store.
			s := store.New().
				Generate(store.If(store.True).Return())
	*/
	Return() Generator

	/*
		Then receives a function that takes in a Store and transforms it to
		generate a new one. The default operational validity of the new Store
		is true and it can be modified by calling `With`.

			s := store.New()
			x := store.NewVar(s, 1)
			s = s.Generate(
				store.If(store.True).
					Then(func(s store.Store) store.Store {
		                // Generate a new store in which the value of x
						// increases by 1.
						return s.Apply(x.Set(x.Get(s) + 1))
					}), // Can use .With(...) to set operational validity
			)
	*/
	Then(func(Store) Store) Generator
}

// Condition represents a logical condition on a Store.
type Condition func(Store) bool

// Change a Store.
type Change func(Store)

// Formatter maps a Store to any type with a JSON representation. It is meant
// to be used with the store.Format function.
type Formatter func(Store) any

// Bounder maps a Store to monotonically tightening bounds. It is meant to be
// used with the store.Bound function.
type Bounder func(Store) Bounds

// Propagator propagates Changes to a Store. It is meant to be used with the
// store.Propagate function.
type Propagator func(Store) []Change

// Valuer maps a Store to an integer value. It is meant to be used with the
// store.Value function.
type Valuer func(Store) int

// Sense specifies whether one is maximizing, minimizing, or satisfying.
// Default is set to minimization.
type Sense int

func (s Sense) String() string {
	switch s {
	case Maximize:
		return "maximize"
	case Satisfy:
		return "satisfy"
	case Minimize:
		return "minimize"
	default:
		panic("sense not defined")
	}
}

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

// MarshalJSON Options.
func (o Options) MarshalJSON() ([]byte, error) {
	search := map[string]any{}
	search["buffer"] = o.Search.Buffer
	m := map[string]any{
		"diagram": o.Diagram,
		"search":  search,
	}
	if o.Limits != (Limits{}) {
		m["limits"] = o.Limits
	}
	if o.Random.Seed != 0 {
		m["random"] = o.Random
	}
	if o.Sense.String() != "" {
		m["sense"] = o.Sense.String()
	}
	if len(o.Tags) > 0 {
		m["tags"] = o.Tags
	}
	if o.Pool.Size != 0 {
		m["pool"] = o.Pool
	}

	return json.Marshal(m)
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

// MarshalJSON Diagram.
func (d Diagram) MarshalJSON() ([]byte, error) {
	m := map[string]any{"width": d.Width}
	m["expansion"] = d.Expansion

	return json.Marshal(m)
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

// MarshalJSON Limits.
func (l Limits) MarshalJSON() ([]byte, error) {
	m := map[string]any{}
	m["duration"] = l.Duration.String()
	if l.Nodes != math.MaxInt {
		m["nodes"] = l.Nodes
	}
	if l.Solutions != math.MaxInt {
		m["solutions"] = l.Solutions
	}

	return json.Marshal(m)
}

// A Solver searches a space and finds the best Solution possible, this is, the
// best collection of Variable assignments in an operationally valid Store.
type Solver interface {
	/*
		All Solutions found by the Solver. Loop over the channel values to get
		the solutions.

			s := store.New()
		    x := store.NewVar(s, 0)
		    opt := store.DefaultOptions()
		    // Minimizer and Satisfier may also be used
		    solver = s.Generate(...).Value(...).Format(...).Maximizer(opt)
		    all := solver.All(context.Background())
			solutions := make([]store.Solution, len(all))
			for solution := range all {
				solutions = append(solutions, solution)
			}
	*/
	All(context.Context) <-chan Solution

	/*
		Last Solution found by the Solver. When running a Maximizer or
		Minimizer, the last Solution is the best one found (highest or smallest
		value, respectively) with the given options:

			s := store.New()
		    x := store.NewVar(s, 0)
		    opt := store.DefaultOptions()
		    // Minimizer and Satisfier may also be used
		    solver = s.Generate(...).Value(...).Format(...).Maximizer(opt)
			last := solve.Last(context.Background())


		Using this function is equivalent to getting the last element when
		using All:

		    last := solutions[len(solutions)-1]

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

// MarshalJSON Time.
func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"start":           t.Start,
		"elapsed":         t.Elapsed.String(),
		"elapsed_seconds": t.Elapsed.Seconds(),
	})
}

// Bounds on an objective value at some node in the search tree consist of a
// lower value and an upper value. If the lower and upper value are the same,
// the bounds have converged.
type Bounds struct {
	Lower int `json:"lower"`
	Upper int `json:"upper"`
}

// New returns a new Store.
func New() Store {
	connect()
	return newFunc()
}

// And uses the conditional "AND" logical operator on all given conditions. It
// returns true if all conditions are true.
func And(c1 Condition, c2 Condition, conditions ...Condition) Condition {
	connect()
	return andFunc(c1, c2, conditions...)
}

// False is a convenience function that is always false.
func False(s Store) bool {
	connect()
	return falseFunc(s)
}

// Not negates the given condition.
func Not(c Condition) Condition {
	connect()
	return notFunc(c)
}

// Or uses the conditional "OR" logical operator on all given conditions. It
// returns true if at least one condition is true.
func Or(c1 Condition, c2 Condition, conditions ...Condition) Condition {
	connect()
	return orFunc(c1, c2, conditions...)
}

// True is a convenience function that is always true.
func True(s Store) bool {
	connect()
	return trueFunc(s)
}

// Xor uses the conditional "Exclusive OR" logical operator on all given
// conditions. It returns true if, and only if, the conditions are different.
func Xor(c1, c2 Condition) Condition {
	connect()
	return xorFunc(c1, c2)
}

/*
If specifies under what condition a Generator can be used. It is meant to be
used with the store.Generate function.

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
	connect()
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
	connect()
	return scopeFunc(f)
}

/*
DefaultOptions for running a solver. Options can be customized after using
these sensitive defaults.

	opt := store.DefaultOptions()
	opt.Limits.Duration = time.Duration(5) * time.Second
*/
func DefaultOptions() Options {
	connect()
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
