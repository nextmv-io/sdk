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
	enc.Encode(s)

	x := store.Var(s, 42)
	y := store.Var(s, "foo")
	pi := store.Var(s, 3.14)
	enc.Encode(s)

	s = s.Format(func(s types.Store) any {
		return map[string]any{
			"x":  x.Get(s),
			"y":  y.Get(s),
			"pi": pi.Get(s),
		}
	})
	enc.Encode(s)

	enc.Encode(s.Apply(y.Set("bar")))
}
