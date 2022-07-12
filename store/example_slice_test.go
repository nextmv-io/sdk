package store_test

import (
	"fmt"

	"github.com/nextmv-io/sdk/store"
)

func ExampleSlice_append() {
	s1 := store.New()
	x := store.NewSlice(s1, 1, 2, 3) // [1, 2, 3]
	s2 := s1.Apply(x.Append(4, 5))
	fmt.Println(x.Slice(s2))
	// Output:
	// [1 2 3 4 5]
}

func ExampleSlice_get() {
	s := store.New()
	x := store.NewSlice(s, 1, 2, 3)
	fmt.Println(x.Get(s, 2))
	// Output:
	// 3
}
