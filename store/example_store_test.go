package store_test

import (
	"context"
	"encoding/json"
	"fmt"
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

// Maximize a store's value by increasing an integer variable. Discard the
// store as operationally invalid as soon as the value is 3. Although the
// generation logic allows for up to 1000 increments of operationally invalid
// stores, the store "stops on its tracks" when the value is 3 because the
// Discard function is used.
func ExampleAction_discard() {
	s := store.New()
	x := store.NewVar(s, 1)
	s = s.
		Generate(
			store.
				If(func(s store.Store) bool { return x.Get(s) == 3 }).
				Discard(),
			store.
				If(func(s store.Store) bool { return x.Get(s) < 1000 }).
				Then(func(s store.Store) store.Store {
					return s.Apply(x.Set(x.Get(s) + 1))
				}).
				With(store.False),
		).
		Format(func(s store.Store) any { return map[string]int{"x": x.Get(s)} }).
		Value(x.Get)

	// The solver type is a maximizer because the store should increase the
	// value of the number. Narrow down the search for performance.
	opt := store.DefaultOptions()
	opt.Diagram.Expansion.Limit = 1
	opt.Diagram.Width = 1
	opt.Search.Buffer = 1
	opt.Limits.Nodes = 3
	solver := s.Maximizer(opt)

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

// Return the store as operationally valid as soon as the value is 3, as in the
// Discard example.
func ExampleAction_return() {
	s := store.New()
	x := store.NewVar(s, 1)
	s = s.
		Generate(
			store.
				If(func(s store.Store) bool { return x.Get(s) == 3 }).
				Return(),
			store.
				If(func(s store.Store) bool { return x.Get(s) < 1000 }).
				Then(func(s store.Store) store.Store {
					return s.Apply(x.Set(x.Get(s) + 1))
				}),
		).
		Format(func(s store.Store) any { return map[string]int{"x": x.Get(s)} }).
		Value(x.Get)

	// The solver type is a maximizer because the store should increase the
	// value of the number. Narrow down the search for performance.
	opt := store.DefaultOptions()
	opt.Diagram.Expansion.Limit = 1
	opt.Diagram.Width = 1
	opt.Search.Buffer = 1
	opt.Limits.Nodes = 3
	solver := s.Maximizer(opt)

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "x": 3
	// }
}

// Create the Fibonacci sequence for the first 10 numbers. Using Then showcases
// how the store is transformed to include a new number.
func ExampleAction_then() {
	s := store.New()
	f := store.NewSlice(s, 1, 1)
	s = s.
		Generate(
			store.
				If(func(s store.Store) bool { return f.Len(s) < 10 }).
				Then(func(s store.Store) store.Store {
					last := f.Get(s, f.Len(s)-1)
					preceding := f.Get(s, f.Len(s)-2)
					return s.Apply(f.Append(last + preceding))
				}),
		).
		Value(f.Len).
		Format(func(s store.Store) any { return map[string]any{"f": f.Slice(s)} })

	// The solver type is a maximizer because the store should incorporate more
	// numbers into the sequence.
	solver := s.Maximizer(store.DefaultOptions())

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "f": [
	//     1,
	//     1,
	//     2,
	//     3,
	//     5,
	//     8,
	//     13,
	//     21,
	//     34,
	//     55
	//   ]
	// }
}

// If the value of an integer value is smaller than 10, increase its value by
// 1. The generating condition imposes a limit on the search.
func ExampleIf() {
	s := store.New()
	x := store.NewVar(s, 0)
	s = s.
		Generate(
			store.
				If(func(s store.Store) bool { return x.Get(s) < 10 }).
				Then(func(s store.Store) store.Store {
					return s.Apply(x.Set(x.Get(s) + 1))
				}),
		).
		Format(func(s store.Store) any { return map[string]int{"x": x.Get(s)} }).
		Value(x.Get)

	// The solver type is a maximizer because x should increase in value.
	solver := s.Maximizer(store.DefaultOptions())

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "x": 10
	// }
}

// Set custom bounds on the store to refine the search for the best value.
func ExampleBounder() {
	s := store.
		New().
		Bound(func(s store.Store) store.Bounds {
			return store.Bounds{Lower: 15, Upper: 42}
		})

	// The solver type is a minimizer to use bounds.
	solver := s.Minimizer(store.DefaultOptions())

	// Get the bounds of the solution and print them.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Statistics.Bounds, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "lower": 15,
	//   "upper": 42
	// }
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

func ExampleFormatter() {
	s := store.New()
	x := store.NewSlice(s, 1, 2, 3)

	// Apply some custom formatting.
	f1 := func(s store.Store) any { return x.Slice(s) }
	s1 := s.Format(f1)
	b, err := json.Marshal(s1)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	// Apply another custom formatting.
	f2 := func(s store.Store) any { return map[string][]int{"x": x.Slice(s)} }
	s2 := s.Format(f2)
	b, err = json.Marshal(s2)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// [1,2,3]
	// {"x":[1,2,3]}
}

// If a Generator is not used, new Stores are not created and thus there is no
// search.
func ExampleGenerator() {
	s := store.New().Generate()

	// Define any solver.
	solver := s.Satisfier(store.DefaultOptions())

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	fmt.Println(last.Store)
	// Output:
	// <nil>
}

// Define a custom operational validity.
func ExampleGenerator_with() {
	s := store.New()
	x := store.NewVar(s, 0)
	s = s.
		Generate(
			store.
				If(store.True).
				Then(func(s store.Store) store.Store {
					return s.Apply(x.Set(x.Get(s) + 1))
				}).
				With(func(s store.Store) bool {
					// The Store is operationally valid if x is even.
					return x.Get(s)%2 == 0
				}),
		).
		Value(x.Get)

	// The solver type is a maximizer to increase the value of x.
	opt := store.DefaultOptions()
	opt.Limits.Solutions = 1
	opt.Diagram.Expansion.Limit = 1
	opt.Limits.Nodes = 10
	solver := s.Maximizer(opt)

	// Get the stats of the solution and print them.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Statistics.Search, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "generated": 2,
	//   "filtered": 0,
	//   "expanded": 2,
	//   "reduced": 0,
	//   "restricted": 2,
	//   "deferred": 2,
	//   "explored": 0,
	//   "solutions": 1
	// }
}

// Define a Generator using a lexical scope, where variable definitions may be
// reused.
func ExampleScope() {
	s := store.New()
	x := store.NewVar(s, 1)
	s = s.
		Generate(
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
		).
		Value(x.Get)

	// The solver type is a maximizer to increase the value of x.
	opt := store.DefaultOptions()
	solver := s.Maximizer(opt)

	// Get the stats of the solution and print them.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Statistics.Search, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "generated": 9,
	//   "filtered": 0,
	//   "expanded": 9,
	//   "reduced": 0,
	//   "restricted": 9,
	//   "deferred": 0,
	//   "explored": 1,
	//   "solutions": 9
	// }
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
	//     "duration": "10s"
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
