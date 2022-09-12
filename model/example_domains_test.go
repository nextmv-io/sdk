package model_test

import (
	"fmt"

	"github.com/nextmv-io/sdk/model"
)

func ExampleDomains_add() {
	d1 := model.Repeat(3, model.Singleton(42))
	d2 := d1.Add(1, 41, 43)
	fmt.Println(d2)
	// Output:
	// [{[{42 42}]} {[{41 43}]} {[{42 42}]}]
}

func ExampleDomains_assign() {
	d1 := model.Repeat(3, model.Singleton(42))
	d2 := d1.Assign(0, 10)
	fmt.Println(d2)
	// Output:
	// [{[{10 10}]} {[{42 42}]} {[{42 42}]}]
}

func ExampleDomains_atLeast() {
	d1 := model.Repeat(2, model.NewDomain(model.NewRange(1, 100)))
	d2 := d1.AtLeast(1, 50)
	fmt.Println(d2)
	// Output:
	// [{[{1 100}]} {[{50 100}]}]
}

func ExampleDomains_atMost() {
	d1 := model.Repeat(2, model.NewDomain(model.NewRange(1, 100)))
	d2 := d1.AtMost(1, 50)
	fmt.Println(d2)
	// Output:
	// [{[{1 100}]} {[{1 50}]}]
}

func ExampleDomains_cmp() {
	d1 := model.Repeat(2, model.Singleton(42))
	d2 := model.Repeat(3, model.Singleton(43))
	fmt.Println(d1.Cmp(d2))
	// Output:
	// -1
}

func ExampleDomains_domain() {
	d := model.NewDomains(model.NewDomain(), model.Singleton(42))
	fmt.Println(d.Domain(0))
	fmt.Println(d.Domain(1))
	// Output:
	// {[]}
	// {[{42 42}]}
}

func ExampleDomains_empty() {
	d := model.NewDomains(model.NewDomain())
	fmt.Println(d.Empty())
	// Output:
	// true
}

func ExampleDomains_len() {
	d := model.Repeat(5, model.NewDomain())
	fmt.Println(d.Len())
	// Output:
	// 5
}

func ExampleDomains_remove() {
	d1 := model.NewDomains(model.Multiple(42, 13))
	d2 := d1.Remove(0, []int{13})
	fmt.Println(d2)
	// Output:
	// [{[{42 42}]}]
}

func ExampleDomains_singleton() {
	d := model.Repeat(5, model.Singleton(42))
	fmt.Println(d.Singleton())
	// Output:
	// true
}

func ExampleDomains_slices() {
	d := model.NewDomains(model.NewDomain(), model.Multiple(1, 3))
	fmt.Println(d.Slices())
	// Output:
	// [[] [1 3]]
}

func ExampleDomains_values() {
	d1 := model.Repeat(3, model.Singleton(42))
	d2 := d1.Add(0, 41)
	fmt.Println(d1.Values())
	fmt.Println(d2.Values())
	// Output:
	// [42 42 42] true
	// [] false
}

func ExampleDomains_first() {
	d := model.NewDomains(
		model.Singleton(88),
		model.Multiple(1, 3),
		model.Multiple(4, 76),
	)
	fmt.Println(d.First())
	// Output:
	// 1 true
}

func ExampleDomains_largest() {
	d := model.NewDomains(
		model.Singleton(88),
		model.Multiple(1, 3),
		model.Multiple(4, 76, 97),
	)
	fmt.Println(d.Largest())
	// Output:
	// 2 true
}

func ExampleDomains_last() {
	d := model.NewDomains(
		model.Singleton(88),
		model.Multiple(1, 3),
		model.Multiple(4, 76, 97),
		model.Singleton(45),
	)
	fmt.Println(d.Last())
	// Output:
	// 2 true
}

func ExampleDomains_maximum() {
	d := model.NewDomains(
		model.Singleton(88),
		model.Multiple(4, 76, 97),
		model.Multiple(1, 3),
		model.Singleton(45),
	)
	fmt.Println(d.Maximum())
	// Output:
	// 1 true
}

func ExampleDomains_minimum() {
	d := model.NewDomains(
		model.Singleton(88),
		model.Multiple(4, 76, 97),
		model.Multiple(1, 3),
		model.Singleton(45),
	)
	fmt.Println(d.Minimum())
	// Output:
	// 2 true
}

func ExampleDomains_smallest() {
	d := model.NewDomains(
		model.Singleton(88),
		model.Multiple(1, 3),
		model.Multiple(4, 76, 97),
	)
	fmt.Println(d.Smallest())
	// Output:
	// 1 true
}

func ExampleNewDomains() {
	d := model.NewDomains( // [1 to 10, 42, odds]
		model.NewDomain(model.NewRange(1, 10)),
		model.Singleton(42),
		model.Multiple(1, 3, 5, 7),
	)
	fmt.Println(d)
	// Output:
	// [{[{1 10}]} {[{42 42}]} {[{1 1} {3 3} {5 5} {7 7}]}]
}

func ExampleRepeat() {
	d := model.Repeat(3, model.NewDomain(model.NewRange(1, 10)))
	fmt.Println(d)
	// Output:
	// [{[{1 10}]} {[{1 10}]} {[{1 10}]}]
}
