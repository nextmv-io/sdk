/*
Package store provides a modeling kit for decision automation problems. It is
based on the paradigm of "decisions as code". The base interface is the Store: a
space defined by variables and logic. The underlying algorithms search that
space and find the best solution possible, this is, the best collection of
variable assignments. The Store is the root node of a search tree. Child Stores
(nodes) inherit both logic and variables from the parent and may also add new
variables and logic, or overwrite existing ones. Changes to a child do not
impact its parent.

A new Store is defined.

    s := store.New()

Variables are stored in the Store.

    x := store.NewVar(s, 1)
    y := store.NewSlice(s, 2, 3, 4)
    z := store.NewMap[string, int](s)

The Format of the Store can be set and one can get the value of a variable.

    s = s.Format(
        func(s store.Store) any {
            return map[string]any{
                "x": x.Get(s),
                "y": y.Slice(s),
                "z": z.Map(s),
            }
        },
    )

The Value of the Store can be set. When maximizing or minimizing, variable
assignments are chosen so that this value increases or decreases, respectively.

    s = s.Value(
        func(s store.Store) int {
            sum := 0
            for i := 0; i < y.Len(s); i++ {
                sum += y.Get(s, i)
            }
            return x.Get(s) + sum
        },
    )

Changes, like setting a new value on a variable, can be applied to the Store.

    s = s.Apply(
        x.Set(10),
        y.Append(5, 6),
    )

To broaden the search space, new Stores can be generated.

    s = s.Generate(func(s store.Store) store.Generator {
        value := x.Get(s)
        return store.Lazy(
            func() bool {
                return value <= 10
            },
            func() store.Store {
                value++
                return s.Apply(x.Set(value))
            },
        )
    })

To check the operational validity of the Store (all decisions have been made and
they are valid), use the provided function.

    s = s.Validate(func(s store.Store) bool {
        return x.Get(s)%2 == 0
    })

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

To find the best collection of variable assignments in the Store, the last
Solution can be obtained from the given Solver. Alternatively, all Solutions can
be retrieved to debug the search mechanics of the Solver.

    solver := s.Maximizer(opt)
    last := solver.Last(context.Background())
    // all := solver.All(context.Background())
    best := x.Get(last.Store)
    stats := last.Statistics

Runners are provided for convenience when running the Store. They read data and
options and manage the call to the Solver. The `NEXTMV_RUNNER` environment
variable defines the type of runner used.

  - "cli": (Default) Command Line Interface runner. Useful for running from a
    terminal. Can read from a file or stdin and write to a file or stdout.
  - "http": HTTP runner. Useful for sending requests and receiving responses on
    the specified port.

The runner receives a handler that specifies the data type and expects a Solver.

    package main

    import (
        "github.com/nextmv-io/sdk/run"
        "github.com/nextmv-io/sdk/store"
    )

    func main() {
        handler := func(v int, opt store.Options) (store.Solver, error) {
            s := store.New()
            x := store.NewVar(s, v)      // Initialized from the runner.
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
