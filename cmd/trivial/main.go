package main

import (
	"strconv"

	"github.com/nextmv-io/sdk/hop/run"
	"github.com/nextmv-io/sdk/hop/store"
	"github.com/nextmv-io/sdk/hop/store/types"
)

func main() {
	run.Run(handler)
}

func handler(v int, opt types.Options) (types.Solver, error) {
	root := store.New()
	x := store.Var(root, v)
	y := store.Slice(root, 4, 5, 6)
	z := store.Map[int, string](root)

	root = root.Value(
		x.Get,
	).Format(
		func(s types.Store) any {
			return map[string]any{
				"x": x.Get(s),
				"y": y.Slice(s),
				"z": z.Map(s),
			}
		},
	).Generate(
		/*
			store.Scope(
				func(s types.Store) types.Generator {
					return store.If(
						func(s types.Store) bool { return x.Get(s)%2 != 0 },
					).Then(
						func(s types.Store) types.Store {
							return s.Apply(x.Set(x.Get(s) / 2))
						},
					)
				},
			),
			store.Scope(
				func(s types.Store) types.Generator {
					v := x.Get(s)
					f := func(s types.Store) bool { return v%2 != 0 }
					return store.If(f).Then(
						func(s types.Store) types.Store {
							v /= 2
							return s.Apply(x.Set(v))
						},
					)
				},
			),
		*/
		store.Scope(
			func(s types.Store) types.Generator {
				v := x.Get(s)
				f := func(s types.Store) bool { return v%2 != 0 }
				return store.If(f).Then(
					func(s types.Store) types.Store {
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
