package store_test

import (
	"fmt"

	"github.com/nextmv-io/sdk/store"
)

func ExampleMap_delete() {
	s1 := store.New()
	m := store.NewMap[int, string](s1)
	s1 = s1.Apply(
		m.Set(42, "foo"),
		m.Set(13, "bar"),
	)
	s2 := s1.Apply(m.Delete(42))
	fmt.Println(m.Map(s1))
	fmt.Println(m.Map(s2))
	// Output:
	// map[13:bar 42:foo]
	// map[13:bar]
}

func ExampleMap_get() {
	s1 := store.New()
	m := store.NewMap[int, string](s1)
	s2 := s1.Apply(m.Set(42, "foo"))
	fmt.Println(m.Get(s2, 42))
	fmt.Println(m.Get(s2, 88))
	// Output:
	// foo true
	//  false
}

func ExampleMap_len() {
	s1 := store.New()
	m := store.NewMap[int, string](s1)
	s2 := s1.Apply(
		m.Set(42, "foo"),
		m.Set(13, "bar"),
	)
	fmt.Println(m.Len(s1))
	fmt.Println(m.Len(s2))
	// Output:
	// 0
	// 2
}

func ExampleMap_map() {
	s1 := store.New()
	m := store.NewMap[int, string](s1)
	s2 := s1.Apply(
		m.Set(42, "foo"),
		m.Set(13, "bar"),
	)
	fmt.Println(m.Map(s2))
	// Output:
	// map[13:bar 42:foo]
}

func ExampleMap_set() {
	s1 := store.New()
	m := store.NewMap[int, string](s1)
	s2 := s1.Apply(m.Set(42, "foo"))
	s3 := s2.Apply(m.Set(42, "bar"))
	fmt.Println(m.Map(s2))
	fmt.Println(m.Map(s3))
	// Output:
	// map[42:foo]
	// map[42:bar]
}

func ExampleNewMap() {
	s := store.New()
	m1 := store.NewMap[int, [2]float64](s)
	m2 := store.NewMap[string, int](s)
	s1 := s.Apply(m1.Set(2, [2]float64{0.1, 3.1416}))
	s2 := s.Apply(m2.Set("a", 43))
	fmt.Println(m1.Map(s1))
	fmt.Println(m2.Map(s2))
	// Output:
	// map[2:[0.1 3.1416]]
	// map[a:43]
}
