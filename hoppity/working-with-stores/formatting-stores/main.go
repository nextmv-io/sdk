package main

import (
	"encoding/json"
	"os"

	"github.com/nextmv-io/sdk/hop/store"
	"github.com/nextmv-io/sdk/hop/store/types"
)

func main() {
	enc := json.NewEncoder(os.Stdout)

	s := store.New()
	if err := enc.Encode(s); err != nil {
		panic(err)
	}

	x := store.Var(s, 42)
	y := store.Var(s, "foo")
	pi := store.Var(s, 3.14)
	if err := enc.Encode(s); err != nil {
		panic(err)
	}

	s = s.Format(func(s types.Store) any {
		return map[string]any{
			"x":  x.Get(s),
			"y":  y.Get(s),
			"pi": pi.Get(s),
		}
	})
	if err := enc.Encode(s); err != nil {
		panic(err)
	}

	if err := enc.Encode(s.Apply(y.Set("bar"))); err != nil {
		panic(err)
	}
}
