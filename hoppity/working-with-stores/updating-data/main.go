package main

import (
	"fmt"

	"github.com/nextmv-io/sdk/hop/store"
)

func main() {
	s1 := store.New()
	x := store.Var(s1, 42)
	y := store.Var(s1, "foo")

	s2 := s1.Apply(y.Set("bar"))
	pi := store.Var(s2, 3.14)

	fmt.Println("s1:", x.Get(s1), y.Get(s1))
	fmt.Println("s2:", x.Get(s2), y.Get(s2), pi.Get(s2))
}
