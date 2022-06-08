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
	)

	// child = child.Generate(
	// 	// Infeasible states to infinity.
	// 	context.If(context.True).Discard(),
	// 	// Feasible states to infinity.
	// 	context.If(context.True).Return(),
	// 	// Never generate anything.
	// 	context.If(context.True).Then(nil).With(context.True),
	// 	// Never generate anything.
	// 	context.If(context.True).Then(nil),
	// 	// Use a lexical scope.
	// 	context.Scope(func(ctx context.Context) context.Generator {
	// 		// Update variables and the model here!
	// 		return context.If(context.True).Then(nil).With(context.True)
	// 	}),
	// )

	return child.Maximizer(opt), nil
}
