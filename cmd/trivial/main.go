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
	)

	return child.Maximizer(opt), nil
}
