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
	d := store.Singleton(root, 42)
	ds := store.Repeat(root, 10, model.Domain(model.Range(1, 10)))

	root = root.Value(
		x.Get,
	).Format(
		func(s types.Store) any {
			return map[string]any{
				"x":  x.Get(s),
				"y":  y.Slice(s),
				"z":  z.Map(s),
				"d":  d.Domain(s),
				"ds": ds.Domains(s),
			}
		},
	).Generate(
		store.If(
			func(s types.Store) bool {
				return x.Get(s)%2 != 0
			},
		).Then(
			func(s types.Store) types.Store {
				v := x.Get(s) / 2
				child := s.Apply(
					x.Set(v),
					y.Prepend(v, v*2, v*v),
					y.Append(v/2, v/4, v/8),
					z.Set(v, strconv.Itoa(v)),
					d.Add(v),
					ds.Add(1, v),
					ds.Assign(2, 42),
				)
				return child
			},
		),

		store.If(store.True).Return(),
	)

	return root.Maximizer(opt), nil
}
