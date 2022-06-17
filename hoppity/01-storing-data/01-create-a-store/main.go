package main

import (
	"fmt"

	"github.com/nextmv-io/sdk/hop/store"
)

func main() {
	s := store.New()

	x := store.Var(s, 42)
	y := store.Var(s, []float64{3.14, 2.72})

	fmt.Println(
		x.Get(s)*10,
		y.Get(s)[0],
	)
}
