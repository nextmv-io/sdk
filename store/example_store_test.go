package store_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/nextmv-io/sdk/store"
)

func ExampleFalse() {
	s := store.New()
	f := store.False(s)
	fmt.Println(f)
	// Output: false
}

func ExampleTrue() {
	s := store.New()
	f := store.True(s)
	fmt.Println(f)
	// Output: true
}

func ExampleAnd() {
	s := store.New()
	c := store.And(store.True, store.False)(s)
	fmt.Println(c)
	c = store.And(store.True, store.True)(s)
	fmt.Println(c)
	c = store.And(store.False, store.True)(s)
	fmt.Println(c)
	c = store.And(store.False, store.False)(s)
	fmt.Println(c)
	c = store.And(store.False, store.True, store.False)(s)
	fmt.Println(c)
	c = store.And(store.True, store.True, store.True)(s)
	fmt.Println(c)
	// Output:
	// false
	// true
	// false
	// false
	// false
	// true
}

func ExampleNot() {
	s := store.New()
	f := store.Not(store.True)(s)
	fmt.Println(f)
	t := store.Not(store.False)(s)
	fmt.Println(t)
	// Output:
	// false
	// true
}

func ExampleOr() {
	s := store.New()
	c := store.Or(store.True, store.False)(s)
	fmt.Println(c)
	c = store.Or(store.True, store.True)(s)
	fmt.Println(c)
	c = store.Or(store.False, store.True)(s)
	fmt.Println(c)
	c = store.Or(store.False, store.False)(s)
	fmt.Println(c)
	c = store.Or(store.False, store.True, store.False)(s)
	fmt.Println(c)
	c = store.Or(store.True, store.True, store.True)(s)
	fmt.Println(c)
	// Output:
	// true
	// true
	// true
	// false
	// true
	// true
}

func ExampleXor() {
	s := store.New()
	c := store.Xor(store.True, store.False)(s)
	fmt.Println(c)
	c = store.Xor(store.True, store.True)(s)
	fmt.Println(c)
	c = store.Xor(store.False, store.True)(s)
	fmt.Println(c)
	c = store.Xor(store.False, store.False)(s)
	fmt.Println(c)
	// Output:
	// true
	// false
	// true
	// false
}

// Changes can be applied to a store.
func ExampleChange() {
	// Original value.
	s := store.New()
	x := store.NewVar(s, 15)
	fmt.Println(x.Get(s))

	// Value after store changed.
	s = s.Apply(x.Set(42))
	fmt.Println(x.Get(s))
	// Output:
	// 15
	// 42
}

// Custom conditions can be defined.
func ExampleCondition() {
	s := store.New()
	a := func(store.Store) bool { return 1 < 2 }
	b := func(store.Store) bool { return 1 > 2 }
	c := store.And(a, b)(s)
	fmt.Println(c)
	// Output:
	// false
}

// DefaultOptions provide sensible defaults but they can (and should) be
// modified.
func ExampleDefaultOptions() {
	opt := store.DefaultOptions()
	b, err := json.MarshalIndent(opt, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	// Modify options
	opt.Diagram.Expansion.Limit = 1
	opt.Limits.Duration = time.Duration(4) * time.Second
	opt.Tags = map[string]any{"foo": 1, "bar": 2}
	b, err = json.MarshalIndent(opt, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "diagram": {
	//     "expansion": {
	//       "limit": 0
	//     },
	//     "width": 10
	//   },
	//   "limits": {
	//     "duration": "168h0m0s"
	//   },
	//   "search": {
	//     "buffer": 100
	//   },
	//   "sense": "minimize"
	// }
	// {
	//   "diagram": {
	//     "expansion": {
	//       "limit": 1
	//     },
	//     "width": 10
	//   },
	//   "limits": {
	//     "duration": "4s"
	//   },
	//   "search": {
	//     "buffer": 100
	//   },
	//   "sense": "minimize",
	//   "tags": {
	//     "bar": 2,
	//     "foo": 1
	//   }
	// }
}

// Get all the solutions from the Generate example.
func ExampleSolver_all() {
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

	solver := s.
		Value(func(s store.Store) int { return x.Get(s) }).
		Format(func(s store.Store) any { return x.Get(s) }).
		Maximizer(store.DefaultOptions())

	// Get all solutions.
	all := solver.All(context.Background())

	// Loop over the channel values to get the solutions.
	solutions := make([]store.Store, len(all))
	for solution := range all {
		solutions = append(solutions, solution.Store)
	}

	// Print the solutions.
	b, err := json.MarshalIndent(solutions, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// [
	//   0,
	//   1,
	//   2,
	//   3
	// ]
}

// Get the last (best) solutions from the Generate example.
func ExampleSolver_last() {
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

	solver := s.
		Value(func(s store.Store) int { return x.Get(s) }).
		Format(func(s store.Store) any { return x.Get(s) }).
		Maximizer(store.DefaultOptions())

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// 3
}

// Generate stores lazily: while an integer variable is less than or equal to
// 5, increase its value by 1. Lazy implementation of the Eager example.
func ExampleLazy() {
	s := store.New()
	x := store.NewVar(s, 0)
	solver := s.
		Generate(func(s store.Store) store.Generator {
			value := x.Get(s)
			return store.Lazy(
				func() bool {
					return value <= 5
				},
				func() store.Store {
					value++
					return s.Apply(x.Set(value))
				},
			)
		}).
		Value(func(s store.Store) int { return x.Get(s) }).
		Format(func(s store.Store) any { return x.Get(s) }).
		Maximizer(store.DefaultOptions())

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// 6
}

// Generate stores eagerly: create all stores from an integer variable by
// increasing its value in 1 each time. The value should never be greater than
// 5. Eager implementation of the Lazy example.
func ExampleEager() {
	s := store.New()
	x := store.NewVar(s, 0)
	solver := s.
		Generate(func(s store.Store) store.Generator {
			value := x.Get(s)
			var stores []store.Store
			for value <= 5 {
				value++
				if value > 5 {
					break
				}
				stores = append(stores, s.Apply(x.Set(value)))
			}
			return store.Eager(stores...)
		}).
		Value(func(s store.Store) int { return x.Get(s) }).
		Format(func(s store.Store) any { return x.Get(s) }).
		Maximizer(store.DefaultOptions())

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// 5
}

// Applying changes to a store updates it, e.g.: setting the value of a
// variable.
func ExampleStore_apply() {
	s := store.New()
	x := store.NewVar(s, 3.1416)
	s1 := s.Apply(
		x.Set(x.Get(s) * 2),
	)

	fmt.Println(x.Get(s))
	fmt.Println(x.Get(s1))
	// Output:
	// 3.1416
	// 6.2832
}

// Make an initial value approach a target by minimizing the absolute
// difference between them. The store is bounded near zero to help the solver
// look for the best solution. The resulting bounds are tightened.
func ExampleStore_bound() {
	s := store.New()
	initial := 10
	target := 16
	x := store.NewVar(s, initial)
	s = s.Bound(func(s store.Store) store.Bounds {
		return store.Bounds{
			Lower: -1,
			Upper: 1,
		}
	})

	solver := s.
		Value(func(s store.Store) int {
			diff := float64(target - x.Get(s))
			return int(math.Abs(diff))
		}).
		Generate(func(s store.Store) store.Generator {
			value := x.Get(s)
			return store.Lazy(
				func() bool {
					return value <= 2*target
				},
				func() store.Store {
					value++
					return s.Apply(x.Set(value))
				},
			)
		}).
		Format(func(s store.Store) any { return x.Get(s) }).
		Minimizer(store.DefaultOptions())

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())

	// Override this variable to have a consistent testable example.
	last.Statistics.Time = store.Time{}

	b, err := json.MarshalIndent(last, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "store": 16,
	//   "statistics": {
	//     "bounds": {
	//       "lower": -1,
	//       "upper": 0
	//     },
	//     "search": {
	//       "generated": 23,
	//       "filtered": 0,
	//       "expanded": 23,
	//       "reduced": 0,
	//       "restricted": 10,
	//       "deferred": 13,
	//       "explored": 1,
	//       "solutions": 2
	//     },
	//     "time": {
	//       "elapsed": "0s",
	//       "elapsed_seconds": 0,
	//       "start": "0001-01-01T00:00:00Z"
	//     },
	//     "value": 0
	//   }
	// }
}

// A store can be formatted to any JSON representation.
func ExampleStore_format() {
	s := store.New()
	x := store.NewVar(s, 10)
	s = s.Format(func(s store.Store) any {
		return map[string]int{"x": x.Get(s)}
	})

	b, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {"x":10}
}

// Given a parent, which is simply an integer variable, children are generated
// by adding 1. This is done until the value reaches a certain limit.
func ExampleStore_generate() {
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

	solver := s.
		Value(func(s store.Store) int { return x.Get(s) }).
		Format(func(s store.Store) any { return x.Get(s) }).
		Maximizer(store.DefaultOptions())

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// 3
}

// Validating that 1 is divisible by 2 results in an operational invalid store,
// represented as null.
func ExampleStore_validate() {
	s := store.New()
	x := store.NewVar(s, 1)
	s = s.Validate(func(s store.Store) bool {
		return x.Get(s)%2 == 0
	})

	solver := s.
		Format(func(s store.Store) any { return x.Get(s) }).
		Satisfier(store.DefaultOptions())

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// null
}

// A custom value can be set on a store. Using any solver, the store has the
// given value.
func ExampleStore_value() {
	s := store.New()
	x := store.NewVar(s, 6)
	s = s.Value(func(s store.Store) int {
		v := x.Get(s)
		return v * v
	})

	solver := s.Minimizer(store.DefaultOptions())

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Statistics.Value, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// 36
}

// Increase the value of a variable as much as possible.
func ExampleStore_maximizer() {
	s := store.New()
	x := store.NewVar(s, 10)
	maximizer := s.
		Value(x.Get).
		Format(func(s store.Store) any { return x.Get(s) }).
		Generate(func(s store.Store) store.Generator {
			value := x.Get(s)
			return store.Lazy(
				func() bool { return value <= 20 },
				func() store.Store {
					value += 5
					return s.Apply(x.Set(value))
				},
			)
		}).
		Maximizer(store.DefaultOptions())

	// Get the last solution of the problem and print it.
	last := maximizer.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// 25
}

// Decrease the value of a variable as much as possible.
func ExampleStore_minimizer() {
	s := store.New()
	x := store.NewVar(s, 10)
	minimizer := s.
		Value(x.Get).
		Format(func(s store.Store) any { return x.Get(s) }).
		Generate(func(s store.Store) store.Generator {
			value := x.Get(s)
			return store.Lazy(
				func() bool { return value >= 0 },
				func() store.Store {
					value -= 5
					return s.Apply(x.Set(value))
				},
			)
		}).
		Minimizer(store.DefaultOptions())

	// Get the last solution of the problem and print it.
	last := minimizer.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// -5
}

// Find the first number divisible by 6, starting from 100.
func ExampleStore_satisfier() {
	s := store.New()
	x := store.NewVar(s, 100)
	opt := store.DefaultOptions()
	opt.Limits.Solutions = 1
	opt.Diagram.Expansion.Limit = 1
	satisfier := s.
		Format(func(s store.Store) any { return x.Get(s) }).
		Validate(func(s store.Store) bool {
			return x.Get(s)%6 == 0
		}).
		Generate(func(s store.Store) store.Generator {
			value := x.Get(s)
			return store.Lazy(
				func() bool {
					return value > 0
				},
				func() store.Store {
					value--
					return s.Apply(x.Set(value))
				},
			)
		}).
		Satisfier(opt)

	// Get the last solution of the problem and print it.
	last := satisfier.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// 96
}
