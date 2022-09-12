package store_test

import (
	"fmt"

	"github.com/nextmv-io/sdk/model"
	"github.com/nextmv-io/sdk/store"
)

func ExampleDomains_add() {
	s1 := store.New()
	d := store.Repeat(s1, 3, model.Singleton(42))
	s2 := s1.Apply(d.Add(1, 41, 43))
	fmt.Println(d.Domains(s2))
	// Output:
	// [{[{42 42}]} {[{41 43}]} {[{42 42}]}]
}

func ExampleDomains_assign() {
	s1 := store.New()
	d := store.Repeat(s1, 3, model.Singleton(42))
	s2 := s1.Apply(d.Assign(0, 10))
	fmt.Println(d.Domains(s2))
	// Output:
	// [{[{10 10}]} {[{42 42}]} {[{42 42}]}]
}

func ExampleDomains_atLeast() {
	s1 := store.New()
	d := store.Repeat(s1, 2, model.NewDomain(model.NewRange(1, 100)))
	s2 := s1.Apply(d.AtLeast(1, 50))
	fmt.Println(d.Domains(s2))
	// Output:
	// [{[{1 100}]} {[{50 100}]}]
}

func ExampleDomains_atMost() {
	s1 := store.New()
	d := store.Repeat(s1, 2, model.NewDomain(model.NewRange(1, 100)))
	s2 := s1.Apply(d.AtMost(1, 50))
	fmt.Println(d.Domains(s2))
	// Output:
	// [{[{1 100}]} {[{1 50}]}]
}

func ExampleDomains_cmp() {
	s := store.New()
	d1 := store.Repeat(s, 2, model.Singleton(42))
	d2 := store.Repeat(s, 3, model.Singleton(43))
	fmt.Println(d1.Cmp(s, d2))
	// Output:
	// -1
}

func ExampleDomains_domain() {
	s := store.New()
	d := store.NewDomains(s, model.NewDomain(), model.Singleton(42))
	fmt.Println(d.Domain(s, 0))
	fmt.Println(d.Domain(s, 1))
	// Output:
	// {[]}
	// {[{42 42}]}
}

func ExampleDomains_domains() {
	s := store.New()
	d := store.NewDomains(s, model.NewDomain(), model.Singleton(42))
	fmt.Println(d.Domains(s))
	// Output:
	// [{[]} {[{42 42}]}]
}

func ExampleDomains_empty() {
	s := store.New()
	d := store.NewDomains(s, model.NewDomain())
	fmt.Println(d.Empty(s))
	// Output:
	// true
}

func ExampleDomains_len() {
	s := store.New()
	d := store.Repeat(s, 5, model.NewDomain())
	fmt.Println(d.Len(s))
	// Output:
	// 5
}

func ExampleDomains_remove() {
	s1 := store.New()
	d := store.NewDomains(s1, model.Multiple(42, 13))
	s2 := s1.Apply(d.Remove(0, []int{13}))
	fmt.Println(d.Domains(s2))
	// Output:
	// [{[{42 42}]}]
}

func ExampleDomains_singleton() {
	s := store.New()
	d := store.Repeat(s, 5, model.Singleton(42))
	fmt.Println(d.Singleton(s))
	// Output:
	// true
}

func ExampleDomains_slices() {
	s := store.New()
	d := store.NewDomains(s, model.NewDomain(), model.Multiple(1, 3))
	fmt.Println(d.Slices(s))
	// Output:
	// [[] [1 3]]
}

func ExampleDomains_values() {
	s1 := store.New()
	d := store.Repeat(s1, 3, model.Singleton(42))
	s2 := s1.Apply(d.Add(0, 41))
	fmt.Println(d.Values(s1))
	fmt.Println(d.Values(s2))
	// Output:
	// [42 42 42] true
	// [] false
}

func ExampleDomains_first() {
	s := store.New()
	d := store.NewDomains(
		s,
		model.Singleton(88),
		model.Multiple(1, 3),
		model.Multiple(4, 76),
	)
	fmt.Println(d.First(s))
	// Output:
	// 1 true
}

func ExampleDomains_largest() {
	s := store.New()
	d := store.NewDomains(
		s,
		model.Singleton(88),
		model.Multiple(1, 3),
		model.Multiple(4, 76, 97),
	)
	fmt.Println(d.Largest(s))
	// Output:
	// 2 true
}

func ExampleDomains_last() {
	s := store.New()
	d := store.NewDomains(
		s,
		model.Singleton(88),
		model.Multiple(1, 3),
		model.Multiple(4, 76, 97),
		model.Singleton(45),
	)
	fmt.Println(d.Last(s))
	// Output:
	// 2 true
}

func ExampleDomains_maximum() {
	s := store.New()
	d := store.NewDomains(
		s,
		model.Singleton(88),
		model.Multiple(4, 76, 97),
		model.Multiple(1, 3),
		model.Singleton(45),
	)
	fmt.Println(d.Maximum(s))
	// Output:
	// 1 true
}

func ExampleDomains_minimum() {
	s := store.New()
	d := store.NewDomains(
		s,
		model.Singleton(88),
		model.Multiple(4, 76, 97),
		model.Multiple(1, 3),
		model.Singleton(45),
	)
	fmt.Println(d.Minimum(s))
	// Output:
	// 2 true
}

func ExampleDomains_smallest() {
	s := store.New()
	d := store.NewDomains(
		s,
		model.Singleton(88),
		model.Multiple(1, 3),
		model.Multiple(4, 76, 97),
	)
	fmt.Println(d.Smallest(s))
	// Output:
	// 1 true
}

func ExampleNewDomains() {
	s := store.New()
	d := store.NewDomains( // [1 to 10, 42, odds]
		s,
		model.NewDomain(model.NewRange(1, 10)),
		model.Singleton(42),
		model.Multiple(1, 3, 5, 7),
	)
	fmt.Println(d.Domains(s))
	// Output:
	// [{[{1 10}]} {[{42 42}]} {[{1 1} {3 3} {5 5} {7 7}]}]
}

func ExampleRepeat() {
	s := store.New()
	d := store.Repeat(s, 3, model.NewDomain(model.NewRange(1, 10)))
	fmt.Println(d.Domains(s))
	// Output:
	// [{[{1 10}]} {[{1 10}]} {[{1 10}]}]
}
