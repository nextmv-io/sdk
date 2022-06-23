package main

import (
	"fmt"

	"github.com/nextmv-io/sdk/hop/store"
)

func main() {
	s1 := store.New()
	x := store.Slice(s1, "c", "d", "e")

	s2 := s1.Apply(
		x.Append("h", "i", "j"),
		x.Prepend("a", "y", "z"),
	)

	s3 := s2.Apply(
		x.Insert(6, "f", "g"),
		x.Remove(2, 2),
		x.Set(1, "b"),
	)

	fmt.Println(x.Slice(s1))
	fmt.Println(x.Slice(s2))
	fmt.Println(x.Slice(s3))
}
