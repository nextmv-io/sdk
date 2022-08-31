package model_test

import (
	"fmt"

	"github.com/nextmv-io/sdk/model"
)

func ExampleDomain_add() {
	d1 := model.Multiple(1, 3, 5)
	d2 := d1.Add(2, 4)
	fmt.Println(d1)
	fmt.Println(d2)
	// Output:
	// {[{1 1} {3 3} {5 5}]}
	// {[{1 5}]}
}

func ExampleDomain_atLeast() {
	d1 := model.NewDomain(
		model.NewRange(1, 10),
		model.NewRange(101, 110),
	)
	d2 := d1.AtLeast(50)
	fmt.Println(d1)
	fmt.Println(d2)
	// Output:
	// {[{1 10} {101 110}]}
	// {[{101 110}]}
}

func ExampleDomain_atMost() {
	d1 := model.NewDomain(
		model.NewRange(1, 10),
		model.NewRange(101, 110),
	)
	d2 := d1.AtMost(50)
	fmt.Println(d1)
	fmt.Println(d2)
	// Output:
	// {[{1 10} {101 110}]}
	// {[{1 10}]}
}

func ExampleDomain_cmp() {
	d1 := model.NewDomain(
		model.NewRange(1, 5),
		model.NewRange(8, 10),
	)
	d2 := model.Multiple(-1, 1)
	fmt.Println(d1.Cmp(d2))
	// Output:
	// 1
}

func ExampleDomain_contains() {
	d := model.NewDomain(model.NewRange(1, 10))
	fmt.Println(d.Contains(5))
	fmt.Println(d.Contains(15))
	// Output:
	// true
	// false
}

func ExampleDomain_empty() {
	d1 := model.NewDomain()
	d2 := model.Singleton(42)
	fmt.Println(d1.Empty())
	fmt.Println(d2.Empty())
	// Output:
	// true
	// false
}

func ExampleDomain_len() {
	d := model.NewDomain(
		model.NewRange(1, 10),
		model.NewRange(-5, -1),
	)
	fmt.Println(d.Len())
	// Output:
	// 15
}

func ExampleDomain_max() {
	d1 := model.NewDomain()
	d2 := model.NewDomain(
		model.NewRange(1, 10),
		model.NewRange(-5, -1),
	)
	fmt.Println(d1.Max())
	fmt.Println(d2.Max())
	// Output:
	// 9223372036854775807 false
	// 10 true
}

func ExampleDomain_min() {
	d1 := model.NewDomain()
	d2 := model.NewDomain(
		model.NewRange(1, 10),
		model.NewRange(-5, -1),
	)
	fmt.Println(d1.Min())
	fmt.Println(d2.Min())
	// Output:
	// -9223372036854775808 false
	// -5 true
}

func ExampleDomain_remove() {
	domain := model.NewDomain(model.NewRange(0, 3))
	domain = domain.Remove(2)
	fmt.Println(domain.Min())
	fmt.Println(domain.Contains(2))
	fmt.Println(domain.Max())
	// Output:
	// 0 true
	// false
	// 3 true
}

func ExampleDomain_slice() {
	d := model.NewDomain(model.NewRange(1, 5))
	fmt.Println(d.Slice())
	// Output:
	// [1 2 3 4 5]
}

func ExampleDomain_value() {
	d1 := model.NewDomain()
	d2 := model.Singleton(42)
	d3 := model.Multiple(1, 3, 5)
	fmt.Println(d1.Value())
	fmt.Println(d2.Value())
	fmt.Println(d3.Value())
	// Output:
	// 0 false
	// 42 true
	// 0 false
}

func ExampleNewDomain() {
	domain := model.NewDomain(
		model.NewRange(0, 1),
		model.NewRange(3, 4),
	)
	fmt.Println(domain.Min())
	fmt.Println(domain.Max())
	fmt.Println(domain.Contains(2))
	// Output:
	// 0 true
	// 4 true
	// false
}

func ExampleSingleton() {
	domain := model.Singleton(1)
	fmt.Println(domain.Value())
	// Output:
	// 1 true
}

func ExampleMultiple() {
	domain := model.Multiple(3, 5, 1)
	fmt.Println(domain.Min())
	fmt.Println(domain.Max())
	fmt.Println(domain.Len())
	fmt.Println(domain.Contains(3))
	// Output:
	// 1 true
	// 5 true
	// 3
	// true
}
