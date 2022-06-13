package main

import (
	"strconv"

	"github.com/nextmv-io/sdk/hop/run"
	"github.com/nextmv-io/sdk/hop/solve"
	"github.com/nextmv-io/sdk/hop/store"
)

func main() {
	run.Run(handler)
}

func handler(v int, opt solve.Options) (solve.Solver, error) {
	root := store.NewStore()
	x := store.Var(root, v)
	y := store.NewSlice(root, 4, 5, 6)
	z := store.NewMap[int, string](root)

	root = root.Value(
		x.Get,
	).Format(
		func(s store.Store) any {
			return map[string]any{
				"x": x.Get(s),
				"y": y.Slice(s),
				"z": z.Map(s),
			}
		},
	).Generate(
		/*
			store.Scope(
				func(s store.Store) store.Generator {
					return store.If(
						func(s store.Store) bool { return x.Get(s)%2 != 0 },
					).Then(
						func(s store.Store) store.Store {
							return s.Apply(x.Set(x.Get(s) / 2))
						},
					)
				},
			),
			store.Scope(
				func(s store.Store) store.Generator {
					v := x.Get(s)
					f := func(s store.Store) bool { return v%2 != 0 }
					return store.If(f).Then(
						func(s store.Store) store.Store {
							v /= 2
							return s.Apply(x.Set(v))
						},
					)
				},
			),
		*/
		store.Scope(
			func(s store.Store) store.Generator {
				v := x.Get(s)
				f := func(s store.Store) bool { return v%2 != 0 }
				return store.If(f).Then(
					func(s store.Store) store.Store {
						v /= 2
						return s.Apply(
							x.Set(v),
							y.Prepend(v, v*2, v*v),
							y.Append(v/2, v/4, v/8),
							z.Set(v, strconv.Itoa(v)),
						)
					},
				)
			},
		),
		store.If(store.True).Return(),
	)

	return root.Maximizer(opt), nil
}
