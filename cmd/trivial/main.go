package main

import (
	"strconv"

	"github.com/nextmv-io/sdk/hop/model"
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
	d := store.Var(root, model.Singleton(42))
	ds := store.Var(root, model.Repeat(10, model.Domain(model.Range(1, 10))))

	root = root.Value(
		x.Get,
	).Format(
		func(s types.Store) any {
			return map[string]any{
				"x":  x.Get(s),
				"y":  y.Slice(s),
				"z":  z.Map(s),
				"d":  d.Get(s),
				"ds": ds.Get(s),
			}
		},
	).Generate(
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
							d.Set(d.Get(s).Add(v)),
							ds.Set(ds.Get(s).Add(1, v)),
						)
					},
				)
			},
		),
		store.If(store.True).Return(),
	)

	return root.Maximizer(opt), nil
}
