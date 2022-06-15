package store_test

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nextmv-io/sdk/hop/store"
	"github.com/nextmv-io/sdk/hop/store/types"
)

func Example_basic() {
	s := store.New()
	x := store.Var(s, 0)
	opt := store.DefaultOptions()
	solver := s.Value(x.Get).
		Generate(
			store.If(
				func(s types.Store) bool {
					v := x.Get(s)
					return v%3 != 0 || v == 0
				},
			).Then(
				func(s types.Store) types.Store {
					return s.Apply(x.Set(x.Get(s) + 1))
				},
			).With(
				func(s types.Store) bool {
					return x.Get(s)%3 == 0
				},
			),
		).
		Format(func(s types.Store) any {
			return map[string]int{"x": x.Get(s)}
		}).
		Maximizer(opt)
	solution := solver.Last(context.Background())
	b, err := json.MarshalIndent(solution, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
