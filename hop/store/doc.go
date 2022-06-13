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

    s := store.NewStore()

Variables are stored in the Store.

    x := store.Var(s, 1)
    y := store.NewSlice(s, 2, 3, 4)
    z := store.NewMap[string, int](s)

The Format of the Store can be set and one can get the value of a Variable.

    s = s.Format(
        func(s store.Store) any {
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
        func(s store.Store) int {
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
            func(s store.Store) store.Generator {
                v := x.Get(s)
                f := func(s store.Store) bool { return v%2 != 0 }
                return store.If(f).
                    Then(func(s store.Store) store.Store {
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
            func(s store.Store) store.Generator {
                v := x.Get(s)
                f := func(s store.Store) bool { return v%2 == 0 }
                return store.If(f).
                    Then(func(s store.Store) store.Store {
                        return s.Apply(x.Set(v + 1))
                    })
            },
        ),
        // If x is greater than 75, then generate the same store with
        // operational validity based on x being divisible by 5.
        store.If(func(s store.Store) bool { return x.Get(s) > 75 }).
            Then(func(s store.Store) store.Store { return s }).
            With(func(s store.Store) bool { return y.Len(s)%5 == 0 }),
        // If x is greater than or equal to 100, return the same Store and it
        // is operationally valid.
        store.If(func(s store.Store) bool { return x.Get(s) >= 100 }).Return(),
    )

Run the Store to maximize or minimize the Value. Alternatively, operational
validity on the Store can be satisfied, in which case setting a Value is not
needed. Note that functions on the Store can be chained.

    package main

    import (
        "github.com/nextmv-io/sdk/hop/run"
        "github.com/nextmv-io/sdk/hop/solve"
        "github.com/nextmv-io/sdk/hop/store"
    )

    func main() {
        run.Run(handler)
    }

    func handler(v int, opt solve.Options) (solve.Solver, error) {
        s := store.NewStore()
        x := store.Var(s, 1)
        y := store.NewSlice(s, 2, 3, 4)
        z := store.NewMap[string, int](s)
        s = s.Format(
            func(s store.Store) any {
                return map[string]any{
                    "x": x.Get(s),
                    "y": y.Slice(s),
                    "z": z.Map(s),
                }
            },
        ).Value(
            func(s store.Store) int {
                sum := 0
                for i := 0; i < y.Len(s); i++ {
                    sum += y.Get(s, i)
                }
                return x.Get(s) + sum
            },
        ).Apply(
            x.Set(10),
            y.Append(5, 6),
        ).Generate(
            store.Scope(
                func(s store.Store) store.Generator {
                    v := x.Get(s)
                    f := func(s store.Store) bool { return v%2 != 0 }
                    return store.If(f).
                        Then(func(s store.Store) store.Store {
                            v /= 2
                            return s.Apply(
                                x.Set(v),
                                y.Prepend(v, v*2, v*v),
                                y.Append(v/2, v/4, v/8),
                            )
                        })
                },
            ),
            store.Scope(
                func(s store.Store) store.Generator {
                    v := x.Get(s)
                    f := func(s store.Store) bool { return v%2 == 0 }
                    return store.If(f).
                        Then(func(s store.Store) store.Store {
                            return s.Apply(x.Set(v + 1))
                        })
                },
            ),
            store.If(func(s store.Store) bool { return x.Get(s) > 75 }).
                Then(func(s store.Store) store.Store { return s }).
                With(func(s store.Store) bool { return y.Len(s)%5 == 0 }),
            store.If(func(s store.Store) bool { return x.Get(s) >= 100 }).
                Return(),
        )

        return s.Maximizer(opt), nil
        // return s.Minimizer(opt), nil
        // return s.Satisfier(opt), nil // Using Store.Value is unnecessary.
    }
*/
package store
