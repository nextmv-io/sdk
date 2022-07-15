package store_test

import (
	"fmt"

	"github.com/nextmv-io/sdk/store"
)

func ExampleSlice_append() {
	s1 := store.New()
	x := store.NewSlice(s1, 1, 2, 3)
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

func ExampleSlice_insert() {
	s1 := store.New()
	x := store.NewSlice(s1, "a", "b", "c")
	s2 := s1.Apply(x.Insert(2, "d", "e"))
	fmt.Println(x.Slice(s2))
	// Output:
	// [a b d e c]
}

func ExampleSlice_len() {
	s := store.New()
	x := store.NewSlice(s, 1, 2, 3)
	fmt.Println(x.Len(s))
	// Output:
	// 3
}

func ExampleSlice_prepend() {
	s1 := store.New()
	x := store.NewSlice(s1, 1, 2, 3)
	s2 := s1.Apply(x.Prepend(4, 5))
	fmt.Println(x.Slice(s2))
	// Output:
	// [4 5 1 2 3]
}

func ExampleSlice_remove() {
	s1 := store.New()
	x := store.NewSlice(s1, 1, 2, 3)
	s2 := s1.Apply(x.Remove(1, 1))
	fmt.Println(x.Slice(s2))
	// Output:
	// [1 3]
}

func ExampleSlice_set() {
	s1 := store.New()
	x := store.NewSlice(s1, "a", "b", "c")
	s2 := s1.Apply(x.Set(1, "d"))
	fmt.Println(x.Slice(s2))
	// Output:
	// [a d c]
}

func ExampleSlice_slice() {
	s := store.New()
	x := store.NewSlice(s, 1, 2, 3)
	fmt.Println(x.Slice(s))
	// Output:
	// [1 2 3]
}

func ExampleNewSlice() {
	s := store.New()
	x := store.NewSlice[int](s)
	y := store.NewSlice(s, 3.14, 2.72)
	fmt.Println(x.Slice(s))
	fmt.Println(y.Slice(s))
	// Output:
	// []
	// [3.14 2.72]
}
