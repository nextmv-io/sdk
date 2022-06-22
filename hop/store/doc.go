/*
Package store provides a modeling kit for decision automation problems. It is
based on the paradigm of "decisions as code". The base interface is the Store:
a space defined by Variables and logic. Hop is an engine provided to search
that space and find the best solution possible, this is, the best collection of
Variable assignments. The Store is the root node of a search tree. Child Stores
(nodes) inherit both logic and Variables from the parent and may also add new
Variables and logic, or overwrite existing ones. Changes to a child do not
impact its parent.

A new Store is defined.

    s := store.New()

Variables are stored in the Store.

    x := store.Var(s, 1)
    y := store.Slice(s, 2, 3, 4)
    z := store.Map[string, int](s)

The Format of the Store can be set and one can get the value of a Variable.

    s = s.Format(
        func(s types.Store) any {
            return map[string]any{
                "x": x.Get(s),
                "y": y.Slice(s),
                "z": z.Map(s),
            }
        },
    )

The Value of the Store can be set. When maximizing or minimizing, Variable
assignments are chosen so that this value increases or decreases, respectively.

    s = s.Value(
        func(s types.Store) int {
            sum := 0
            for i := 0; i < y.Len(s); i++ {
                sum += y.Get(s, i)
            }
            return x.Get(s) + sum
        },
    )

Changes, like setting a new value on a Variable, can be applied to the Store.

    s = s.Apply(
        x.Set(10),
        y.Append(5, 6),
    )

To broaden the search space, new Stores can be generated.

    s = s.Generate(
		// If x is odd, divide the value in half and modify y. Operationally
		// valid.
		store.Scope(
			func(s types.Store) types.Generator {
				v := x.Get(s)
				f := func(s types.Store) bool { return v%2 != 0 }
				return store.If(f).
					Then(func(s types.Store) types.Store {
						v /= 2
						return s.Apply(
							x.Set(v),
							y.Prepend(v, v*2, v*v),
							y.Append(v/2, v/4, v/8),
						)
					})
			},
		),
		// If x is even, increase the value by 1. Operationally valid.
		store.Scope(
			func(s types.Store) types.Generator {
				v := x.Get(s)
				f := func(s types.Store) bool { return v%2 == 0 }
				return store.If(f).
					Then(func(s types.Store) types.Store {
						return s.Apply(x.Set(v + 1))
					})
			},
		),
		// If x is greater than 75, then generate the same store with
		// operational validity based on x being divisible by 5.
		store.If(func(s types.Store) bool { return x.Get(s) > 75 }).
			Then(func(s types.Store) types.Store { return s }).
			With(func(s types.Store) bool { return y.Len(s)%5 == 0 }),
		// If x is greater than or equal to 100, return the same Store and it
		// is operationally valid.
		store.If(func(s types.Store) bool { return x.Get(s) >= 100 }).Return(),
	)

When setting a Value, it can be maximized or minimized. Alternatively,
operational validity on the Store can be satisfied, in which case setting a
Value is not needed. Options are required to specify the search mechanics.

    // DefaultOptions provide sensible defaults.
    opt := store.DefaultOptions()
    // Options can be modified, e.g.: changing the duration.
    // opt.Limits.Duration = time.Duration(4) * time.Second
	solver := s.Value(...).Minimizer(opt)
	// solver := s.Value(...).Minimizer(opt)
	// solver := s.Satisfier(opt)

To find the best collection of Variable assignments in the Store, the last
Solution can be obtained from the given Solver. Alternatively, all Solutions
can be retrieved to debug the search mechanics of the Solver.

    solver := s.Maximizer(opt)
	last := solver.Last(context.Background())
	// all := solver.All(context.Background())
	best := x.Get(last.Store)
	stats := last.Statistics

Runners are provided for convenience when running the Store. They read data and
options and manage the call to the Solver. The `NEXTMV_RUNNER` environment
variable defines the type of runner used:

    - "cli": Command Line Interface runner. Useful for running from a terminal.
    Can read from a file or stdin and write to a file or stdout.
    - "http": HTTP runner. Useful for sending requests and receiving responses
    on the specified port.

The runner receives a handler that specifies the data type and expects a Solver.

    func main() {
        handler := func(v int, opt types.Options) (types.Solver, error) {
            s := store.New()
            x := store.Var(s, v) // Initialized from the runner.
            s = s.Value(...).Format(...) // Modify the Store.

            return s.Maximizer(opt), nil // Options are passed by the runner.
        }
        run.Run(handler)
    }

Compile the binary and use the -h flag to see available options to configure a
runner. You can use command-line flags or environment variables. When using
environment variables, use all caps and snake case. For example, using the
command-line flag `-hop.solver.limits.duration` is equivalent to setting the
environment variable `HOP_SOLVER_LIMITS_DURATION`.

Using the cli runner for example:

    echo 0 | go run main.go -hop.solver.limits.duration 2s

Writes this output to stdout:

    {
      "hop": {
        "version": "..."
      },
      "options": {
        "diagram": {
          "expansion": {
            "limit": 0
          },
          "width": 10
        },
        "limits": {
          "duration": "2s"
        },
        "search": {
          "buffer": 100
        },
        "sense": "maximizer"
      },
      "store": {
        "x": 10
      },
      "statistics": {
        "bounds": {
          "lower": 10,
          "upper": 9223372036854776000
        },
        "search": {
          "generated": 10,
          "filtered": 0,
          "expanded": 10,
          "reduced": 0,
          "restricted": 10,
          "deferred": 0,
          "explored": 1,
          "solutions": 5
        },
        "time": {
          "elapsed": "93.417µs",
          "elapsed_seconds": 9.3417e-05,
          "start": "..."
        },
        "value": 10
      }
    }

*/
package store
