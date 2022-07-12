package store_test

import (
	"fmt"

	"github.com/nextmv-io/sdk/model"
	"github.com/nextmv-io/sdk/store"
)

func ExampleDomain_add() {
	s1 := store.New()
	d := store.Multiple(s1, 1, 3, 5)
	s2 := s1.Apply(d.Add(2, 4))
	fmt.Println(d.Domain(s1))
	fmt.Println(d.Domain(s2))
	// Output:
	// {[{1 1} {3 3} {5 5}]}
	// {[{1 5}]}
}

func ExampleDomain_atLeast() {
	s1 := store.New()
	d := store.NewDomain(s1, model.NewRange(1, 10), model.NewRange(101, 110))
	s2 := s1.Apply(d.AtLeast(50))
	fmt.Println(d.Domain(s1))
	fmt.Println(d.Domain(s2))
	// Output:
	// {[{1 10} {101 110}]}
	// {[{101 110}]}
}

func ExampleDomain_atMost() {
	s1 := store.New()
	d := store.NewDomain(s1, model.NewRange(1, 10), model.NewRange(101, 110))
	s2 := s1.Apply(d.AtMost(50))
	fmt.Println(d.Domain(s1))
	fmt.Println(d.Domain(s2))
	// Output:
	// {[{1 10} {101 110}]}
	// {[{1 10}]}
}

func ExampleDomain_cmp() {
	s := store.New()
	d1 := store.NewDomain(s, model.NewRange(1, 5), model.NewRange(8, 10))
	d2 := store.Multiple(s, -1, 1)
	fmt.Println(d1.Cmp(s, d2))
	// Output:
	// 1
}

func ExampleDomain_contains() {
	s := store.New()
	d := store.NewDomain(s, model.NewRange(1, 10))
	fmt.Println(d.Contains(s, 5))
	fmt.Println(d.Contains(s, 15))
	// Output:
	// true
	// false
}

func ExampleDomain_domain() {
	s := store.New()
	d := store.NewDomain(s, model.NewRange(1, 10))
	fmt.Println(d.Domain(s))
	// Output:
	// {[{1 10}]}
}

func ExampleDomain_empty() {
	s := store.New()
	d1 := store.NewDomain(s)
	d2 := store.Singleton(s, 42)
	fmt.Println(d1.Empty(s))
	fmt.Println(d2.Empty(s))
	// Output:
	// true
	// false
}

func ExampleDomain_len() {
	s := store.New()
	d := store.NewDomain(s, model.NewRange(1, 10), model.NewRange(-5, -1))
	fmt.Println(d.Len(s))
	// Output:
	// 15
}

func ExampleDomain_max() {
	s := store.New()
	d1 := store.NewDomain(s)
	d2 := store.NewDomain(s, model.NewRange(1, 10), model.NewRange(-5, -1))
	fmt.Println(d1.Max(s))
	fmt.Println(d2.Max(s))
	// Output:
	// 9223372036854775807 false
	// 10 true
}

func ExampleDomain_min() {
	s := store.New()
	d1 := store.NewDomain(s)
	d2 := store.NewDomain(s, model.NewRange(1, 10), model.NewRange(-5, -1))
	fmt.Println(d1.Min(s))
	fmt.Println(d2.Min(s))
	// Output:
	// -9223372036854775808 false
	// -5 true
}

func ExampleDomain_remove() {
	s1 := store.New()
	d := store.NewDomain(s1, model.NewRange(1, 5))
	s2 := s1.Apply(d.Remove(2, 4))
	fmt.Println(d.Domain(s1))
	fmt.Println(d.Domain(s2))
	// Output:
	// {[{1 5}]}
	// {[{1 1} {3 3} {5 5}]}
}

func ExampleDomain_slice() {
	s := store.New()
	d := store.NewDomain(s, model.NewRange(1, 5))
	fmt.Println(d.Slice(s))
	// Output:
	// [1 2 3 4 5]
}

func ExampleDomain_value() {
	s := store.New()
	d1 := store.NewDomain(s)
	d2 := store.Singleton(s, 42)
	d3 := store.Multiple(s, 1, 3, 5)
	fmt.Println(d1.Value(s))
	fmt.Println(d2.Value(s))
	fmt.Println(d3.Value(s))
	// Output:
	// 0 false
	// 42 true
	// 0 false
}

func ExampleNewDomain() {
	s := store.New()
	d1 := store.NewDomain(s, model.NewRange(1, 10))
	d2 := store.NewDomain(s, model.NewRange(1, 10), model.NewRange(20, 29))
	fmt.Println(d1.Domain(s))
	fmt.Println(d2.Domain(s))
	// Output:
	// {[{1 10}]}
	// {[{1 10} {20 29}]}
}

func ExampleSingleton() {
	s := store.New()
	fortyTwo := store.Singleton(s, 42)
	fmt.Println(fortyTwo.Domain(s))
	// Output:
	// {[{42 42}]}
}

func ExampleMultiple() {
	s := store.New()
	even := store.Multiple(s, 2, 4, 6, 8)
	fmt.Println(even.Domain(s))
	// Output:
	// {[{2 2} {4 4} {6 6} {8 8}]}
}
