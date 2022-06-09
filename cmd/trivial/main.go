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

	child := root.Apply(
		x.Set(x.Get(root) / 2),
	).Value(
		x.Get,
	).Format(
		func(ctx context.Context) any {
			return map[string]any{"x": x.Get(ctx)}
		},
	).Value(
		x.Get,
	).Format(
		func(ctx context.Context) any {
			return map[string]any{"x": x.Get(ctx)}
		},
	).Generate(
		context.Scope(func(ctx context.Context) context.Generator {
			value := x.Get(ctx)
			f := func() bool { return value < 100 }
			return context.If(f).
				Then(func() context.Context {
					value++
					x.Set(value)
					return ctx
				}).
				With(f)
		}),
	)

	return child.Maximizer(opt), nil
}
