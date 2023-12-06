package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/model"
)

const (
	// Minimize indicates the solution space is being searched to find the
	// smallest possible value.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Minimize Sense = iota
	// Maximize indicates the solution space is being searched to find the
	// biggest possible value.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Maximize
	// Satisfy indicates the solution space is being searched to find
	// operationally valid Stores.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Satisfy
)

/*
Store represents a store of variables and logic to solve decision automation
problems. Adding logic to the Store updates it (functions may be called
directly and chained):

	s := store.New()    // s := store.New().
	s = s.Apply(...)    // 	   Apply(...).
	s = s.Bound(...)    // 	   Bound(...).
	s = s.Format(...)   // 	   Format(...).
	s = s.Generate(...) // 	   Generate(...)

The variables and logic stored define a solution space. This space is searched
to make decisions.

Deprecated: This package is deprecated and will be removed in the future.
*/
type Store interface {
	/*
		Apply changes to a Store. A change happens when a stored variable is
		updated:

			s := store.New()
			x := store.NewVar(s, 3.1416)
			s1 := s.Apply(
				x.Set(x.Get(s) * 2),
			)

		Deprecated: This package is deprecated and will be removed in the future.
	*/
	Apply(...Change) Store

	/*
		Bound the value of a Store. The solver can use this information to more
		efficiently find the best Store. The lower and upper bounds can be set:

			s := store.New()
			x := store.NewVar(s, initial)
			s = s.Bound(func(s store.Store) store.Bounds {
				return store.Bounds{
					Lower: -1,
					Upper: 1,
				}
			})

		Deprecated: This package is deprecated and will be removed in the future.
	*/
	Bound(Bounder) Store

	/*
		Format a Store into any structure prior to JSON encoding.

			s := store.New()
			x := store.NewVar(s, 10)
			s = s.Format(func(s store.Store) any {
				return map[string]int{"x": x.Get(s)}
			})

		Deprecated: This package is deprecated and will be removed in the future.
	*/
	Format(Formatter) Store

	/*
		Generate new Stores (children) from the existing one (parent). A
		callback function provides a lexical scope that can be used to perform
		and update calculations.

			s := store.New()
			x := store.NewVar(s, 0)
			s = s.Generate(func(s store.Store) store.Generator {
				value := x.Get(s)
				return store.Lazy(
					func() bool {
						return value <= 2
					},
					func() store.Store {
						value++
						return s.Apply(x.Set(value))
					},
				)
			})

		Deprecated: This package is deprecated and will be removed in the future.
	*/
	Generate(func(Store) Generator) Store

	/*
		Propagate changes into a Store and is re-invoked until no further
		changes need to be made and an empty slice of changes is returned.

			s := store.New()
			x := store.NewVar(s, 1)
			s = s.Propagate(func(s store.Store) []store.Change {
				if x.Get(s) <= 1 {
					return []store.Change{
						x.Set(2),
						x.Set(42),
					}
				}
				return []store.Change{}
			})

		Deprecated: This package is deprecated and will be removed in the future.
	*/
	Propagate(...Propagator) Store

	/*
		Validate the Store. A Store is operationally valid if all decisions
		have been made and those decisions fulfill certain requirements; e.g.:
		all stops have been assigned to vehicles, all shifts are covered with
		the necessary personnel, all assignment have been made, quantity
		respects an alloted capacity, etc. Setting operational validity is
		optional and the default is true.

			s := store.New()
			x := store.NewVar(s, 1)
			s = s.Validate(func(s store.Store) bool {
				return x.Get(s)%2 == 0
			})

		Deprecated: This package is deprecated and will be removed in the future.
	*/
	Validate(Condition) Store

	/*
		Value sets the integer value of a Store. When maximizing or minimizing,
		this is the value that is optimized.

			s := store.New()
			x := store.NewVar(s, 6)
			s = s.Value(func(s store.Store) int {
				v := x.Get(s)
				return v * v
			})

		Deprecated: This package is deprecated and will be removed in the future.
	*/
	Value(Valuer) Store

	// Maximizer builds a solver that searches the space defined by the Store
	// to maximize a value.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Maximizer(Options) Solver

	// Minimizer builds a solver that searches the space defined by the Store
	// to minimize a value.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Minimizer(Options) Solver

	// Satisfier builds a solver that searches the space defined by the Store
	// to satisfy operational validity.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Satisfier(Options) Solver
}

// A Generator is used to generate new Stores (children) from an existing one
// (parent). It is meant to be used with the store.Generate function.
//
// Deprecated: This package is deprecated and will be removed in the future.
type Generator any

// Condition represents a logical condition on a Store.
//
// Deprecated: This package is deprecated and will be removed in the future.
type Condition func(Store) bool

// Change a Store.
//
// Deprecated: This package is deprecated and will be removed in the future.
type Change func(Store)

// Formatter maps a Store to any type with a JSON representation. It is meant
// to be used with the store.Format function.
//
// Deprecated: This package is deprecated and will be removed in the future.
type Formatter func(Store) any

// Bounder maps a Store to monotonically tightening bounds. It is meant to be
// used with the store.Bound function.
//
// Deprecated: This package is deprecated and will be removed in the future.
type Bounder func(Store) Bounds

// Propagator propagates Changes to a Store. It is meant to be used with the
// store.Propagate function.
//
// Deprecated: This package is deprecated and will be removed in the future.
type Propagator func(Store) []Change

// Valuer maps a Store to an integer value. It is meant to be used with the
// store.Value function.
//
// Deprecated: This package is deprecated and will be removed in the future.
type Valuer func(Store) int

// Sense specifies whether one is maximizing, minimizing, or satisfying.
// Default is set to minimization.
//
// Deprecated: This package is deprecated and will be removed in the future.
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
//
// Deprecated: This package is deprecated and will be removed in the future.
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
		Buffer int `usage:"solution buffer (0 = none)" default:"100"`
	}
	Limits Limits
	// Options for random number generation.
	Random struct {
		// Seed for generating random numbers.
		Seed int64 `json:"seed,omitempty" usage:"random seed"`
	}
	// Pool that is used in specific engines.
	Pool struct {
		// Maximum Size of the Pool.
		Size int `json:"size,omitempty" usage:"pool size (only specific engines)"`
	}
}

// MarshalJSON Options.
//
// Deprecated: This package is deprecated and will be removed in the future.
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
//
// Deprecated: This package is deprecated and will be removed in the future.
type Diagram struct {
	// Maximum Width of the Decision Diagram.
	Width int `usage:"diagram width" default:"10"`
	// Maximum Expansion that can be generated from a Store.
	Expansion struct {
		// Limit represents the maximum number of children Stores that can
		// be generated from a parent.
		Limit int `json:"limit" usage:"diagram expansion limit" default:"1"`
	}
}

// MarshalJSON Diagram.
//
// Deprecated: This package is deprecated and will be removed in the future.
func (d Diagram) MarshalJSON() ([]byte, error) {
	m := map[string]any{"width": d.Width}
	m["expansion"] = d.Expansion

	return json.Marshal(m)
}

// Limits when performing a search. The search will stop if any one of these
// limits are encountered.
type Limits struct {
	// Time Duration.
	Duration time.Duration `usage:"time limit, e.g. 10ms"`
	// Nodes represent active Stores in the search.
	Nodes int `usage:"node limit"`
	// Solutions represent operationally valid Stores.
	Solutions int `usage:"solution limit"`
}

// MarshalJSON Limits.
//
// Deprecated: This package is deprecated and will be removed in the future.
func (l Limits) MarshalJSON() ([]byte, error) {
	m := map[string]any{}
	m["duration"] = l.Duration.String()
	if l.Nodes != model.MaxInt {
		m["nodes"] = l.Nodes
	}
	if l.Solutions != model.MaxInt {
		m["solutions"] = l.Solutions
	}

	return json.Marshal(m)
}

// A Solver searches a space and finds the best Solution possible, this is, the
// best collection of variable assignments in an operationally valid Store.
//
// Deprecated: This package is deprecated and will be removed in the future.
type Solver interface {
	// All Solutions found by the Solver. Loop over the channel values to get
	// the solutions.
	All(context.Context) <-chan Solution

	// Last Solution found by the Solver. When running a Maximizer or
	// Minimizer, the last Solution is the best one found (highest or smallest
	// value, respectively) with the given options. Using this function is
	// equivalent to getting the last element when using All.
	Last(context.Context) Solution

	// Options provided to the Solver.
	Options() Options
}

// Solution of a decision automation problem. A Solution is an operationally
// valid Store.
//
// Deprecated: This package is deprecated and will be removed in the future.
type Solution struct {
	// Store of the Solution. If nil, it means that the solution is
	// operationally invalid.
	Store      Store      `json:"store"`
	Statistics Statistics `json:"statistics"`
}

// Statistics of the search.
//
// Deprecated: This package is deprecated and will be removed in the future.
type Statistics struct {
	// Bounds of the store. Nil when using a Satisfier.
	Bounds *Bounds `json:"bounds,omitempty"`
	Search Search  `json:"search"`
	Time   Time    `json:"time"`
	// Value of the store. Nil when using a Satisfier.
	Value *int `json:"value,omitempty"`
}

// Search statistics of the Store generation.
//
// Deprecated: This package is deprecated and will be removed in the future.
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
//
// Deprecated: This package is deprecated and will be removed in the future.
type Time struct {
	Start   time.Time     `json:"start"`
	Elapsed time.Duration `json:"elapsed"`
}

// MarshalJSON Time.
//
// Deprecated: This package is deprecated and will be removed in the future.
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
//
// Deprecated: This package is deprecated and will be removed in the future.
type Bounds struct {
	Lower int `json:"lower"`
	Upper int `json:"upper"`
}

// New returns a new Store.
//
// Deprecated: This package is deprecated and will be removed in the future.
func New() Store {
	connect.Connect(con, &newFunc)
	return newFunc()
}

// And uses the conditional "AND" logical operator on all given conditions. It
// returns true if all conditions are true.
//
// Deprecated: This package is deprecated and will be removed in the future.
func And(c1 Condition, c2 Condition, conditions ...Condition) Condition {
	connect.Connect(con, &andFunc)
	return andFunc(c1, c2, conditions...)
}

// False is a convenience function that is always false.
//
// Deprecated: This package is deprecated and will be removed in the future.
func False(s Store) bool {
	connect.Connect(con, &falseFunc)
	return falseFunc(s)
}

// Not negates the given condition.
//
// Deprecated: This package is deprecated and will be removed in the future.
func Not(c Condition) Condition {
	connect.Connect(con, &notFunc)
	return notFunc(c)
}

// Or uses the conditional "OR" logical operator on all given conditions. It
// returns true if at least one condition is true.
//
// Deprecated: This package is deprecated and will be removed in the future.
func Or(c1 Condition, c2 Condition, conditions ...Condition) Condition {
	connect.Connect(con, &orFunc)
	return orFunc(c1, c2, conditions...)
}

// True is a convenience function that is always true.
//
// Deprecated: This package is deprecated and will be removed in the future.
func True(s Store) bool {
	connect.Connect(con, &trueFunc)
	return trueFunc(s)
}

// Xor uses the conditional "Exclusive OR" logical operator on all given
// conditions. It returns true if, and only if, the conditions are different.
//
// Deprecated: This package is deprecated and will be removed in the future.
func Xor(c1, c2 Condition) Condition {
	connect.Connect(con, &xorFunc)
	return xorFunc(c1, c2)
}

/*
DefaultOptions for running a solver. Options can be customized after using
these sensitive defaults.

	opt := store.DefaultOptions()
	opt.Limits.Duration = time.Duration(5) * time.Second

Deprecated: This package is deprecated and will be removed in the future.
*/
func DefaultOptions() Options {
	connect.Connect(con, &defaultOptionsFunc)
	return defaultOptionsFunc()
}

// Eager way of generating new Stores. The Generator uses the list of Stores
// upfront in the order they are provided.
//
// Deprecated: This package is deprecated and will be removed in the future.
func Eager(s ...Store) Generator {
	connect.Connect(con, &eagerFunc)
	return eagerFunc(s...)
}

// Lazy way of generating new Stores. While the condition holds, the function
// is called to generate new Stores. If the condition is no longer true or a
// nil Store is returned, the generator is not used anymore by the current
// parent.
//
// Deprecated: This package is deprecated and will be removed in the future.
func Lazy(c func() bool, f func() Store) Generator {
	connect.Connect(con, &lazyFunc)
	return lazyFunc(c, f)
}

var (
	con                = connect.NewConnector("sdk", "Store")
	newFunc            func() Store
	andFunc            func(Condition, Condition, ...Condition) Condition
	falseFunc          func(Store) bool
	notFunc            func(Condition) Condition
	orFunc             func(Condition, Condition, ...Condition) Condition
	trueFunc           func(Store) bool
	xorFunc            func(Condition, Condition) Condition
	defaultOptionsFunc func() Options
	eagerFunc          func(...Store) Generator
	lazyFunc           func(func() bool, func() Store) Generator
)
