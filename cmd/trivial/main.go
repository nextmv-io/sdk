package main

import (
	"github.com/nextmv-io/sdk/hop/context"
	"github.com/nextmv-io/sdk/hop/run"
	"github.com/nextmv-io/sdk/hop/solve"
)

func main() {
	run.Run(handler)
}

func handler(v int, opt solve.Options) (solve.Solver, error) {
	root := context.NewContext()
	x := context.Declare(root, v)
	y := context.NewVector(root, 4, 5, 6)

	root = root.Value(
		x.Get,
	).Format(
		func(ctx context.Context) any {
			return map[string]any{
				"x": x.Get(ctx),
				"y": y.Data(ctx),
			}
		},
	).Generate(
		/*
			context.Scope(
				func(ctx context.Context) context.Generator {
					return context.If(
						func() bool { return x.Get(ctx)%2 != 0 },
					).Then(
						func() context.Context {
							return ctx.Apply(x.Set(x.Get(ctx) / 2))
						},
					)
				},
			),
		*/

		/*
			context.Scope(
				func(ctx context.Context) context.Generator {
					v := x.Get(ctx)
					f := func() bool { return v%2 != 0 }
					return context.If(f).Then(
						func() context.Context {
							v = v / 2
							return ctx.Apply(x.Set(v))
						},
					)
				},
			),
		*/

		context.Scope(
			func(ctx context.Context) context.Generator {
				v := x.Get(ctx)
				f := func() bool { return v%2 != 0 }
				return context.If(f).Then(
					func() context.Context {
						v = v / 2
						return ctx.Apply(
							x.Set(v),
							y.Prepend(v, v*2, v*v),
							y.Append(v/2, v/4, v/8),
						)
					},
				)
			},
		),
		context.If(context.True).Return(),
	)

	return root.Maximizer(opt), nil
}
