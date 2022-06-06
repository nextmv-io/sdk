package main

import (
	"encoding/json"
	"os"

	"github.com/nextmv-io/sdk/hop/context"
)

func main() {
	root := context.NewContext()
	x := context.Declare(root, 42)

	child := root.Apply(
		x.Set(x.Get(root) / 2),
	).Check(
		context.False,
	).Value(
		x.Get,
	).Format(
		func(ctx context.Context) any {
			return map[string]any{"x": x.Get(ctx)}
		},
	)

	json.NewEncoder(os.Stdout).Encode(child)
}
