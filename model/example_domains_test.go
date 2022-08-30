package model_test

import (
	"fmt"

	"github.com/nextmv-io/sdk/model"
)

func ExampleNewDomains_first() {
	domains := model.NewDomains()

	fmt.Println(domains.Len())
	fmt.Println(domains.Empty())

	// Output:
	// 0
	// true
}

func ExampleNewDomains_second() {
	domains := model.NewDomains(
		model.Singleton(0),
		model.Singleton(1),
	)

	fmt.Println(domains.Len())
	fmt.Println(domains.Domain(0).Value())
	fmt.Println(domains.Domain(1).Value())

	// Output:
	// 2
	// 0 true
	// 1 true
}

func ExampleRepeat_first() {
	domains := model.Repeat(0, model.NewDomain())

	fmt.Println(domains.Len())
	fmt.Println(domains.Empty())

	// Output:
	// 0
	// true
}

func ExampleRepeat_second() {
	domains := model.Repeat(1, model.NewDomain())

	fmt.Println(domains.Len())
	fmt.Println(domains.Empty())

	// Output:
	// 1
	// true
}

func ExampleRepeat_third() {
	domains := model.Repeat(2, model.Singleton(1))

	fmt.Println(domains.Len())
	fmt.Println(domains.Domain(0).Value())
	fmt.Println(domains.Domain(1).Value())

	// Output:
	// 2
	// 1 true
	// 1 true
}

func ExampleDomains_Add_first() {
	domains := model.NewDomains(
		model.Singleton(0),
		model.Singleton(1),
	)

	domains = domains.Add(0, 1, 2)

	fmt.Println(domains.Domain(0).Len())

	// Output:
	// 3
}

func ExampleDomains_Add_second() {
	domains := model.NewDomains(
		model.Singleton(0),
		model.Singleton(1),
	)

	domains = domains.Add(1, 1, 2)

	fmt.Println(domains.Domain(1).Len())

	// Output:
	// 2
}

func ExampleDomains_Assign_first() {
	domains := model.NewDomains(
		model.NewDomain(model.NewRange(0, 2)),
	)

	domains = domains.Assign(0, 1)

	fmt.Println(domains.Domain(0).Value())

	// Output:
	// 1 true
}

func ExampleDomains_Assign_second() {
	domains := model.NewDomains(
		model.NewDomain(model.NewRange(0, 2)),
	)

	domains = domains.Assign(0, 3)

	fmt.Println(domains.Domain(0).Value())

	// Output:
	// 3 true
}

func ExampleDomains_AtLeast_first() {
	domains := model.NewDomains(
		model.NewDomain(model.NewRange(0, 2)),
	)

	domains = domains.AtLeast(0, 1)

	fmt.Println(domains.Domain(0).Min())
	fmt.Println(domains.Domain(0).Max())

	// Output:
	// 1 true
	// 2 true
}

func ExampleDomains_AtLeast_second() {
	domains := model.NewDomains(
		model.NewDomain(model.NewRange(0, 2)),
	)

	domains = domains.AtLeast(0, 3)

	fmt.Println(domains.Domain(0).Empty())

	// Output:
	// true
}

func ExampleDomains_AtMost_first() {
	domains := model.NewDomains(
		model.NewDomain(model.NewRange(0, 2)),
	)

	domains = domains.AtMost(0, 1)

	fmt.Println(domains.Domain(0).Min())
	fmt.Println(domains.Domain(0).Max())

	// Output:
	// 0 true
	// 1 true
}

func ExampleDomains_AtMost_second() {
	domains := model.NewDomains(
		model.NewDomain(model.NewRange(0, 2)),
	)

	domains = domains.AtMost(0, -1)

	fmt.Println(domains.Domain(0).Empty())

	// Output:
	// true
}

func ExampleDomains_Remove_first() {
	domains := model.NewDomains(
		model.NewDomain(model.NewRange(0, 2)),
	)

	domains = domains.Remove(0, 1)

	fmt.Println(domains.Domain(0).Min())
	fmt.Println(domains.Domain(0).Max())
	fmt.Println(domains.Domain(0).Contains(1))

	// Output:
	// 0 true
	// 2 true
	// false
}

func ExampleDomains_Slices_first() {
	domains := model.NewDomains(
		model.NewDomain(model.NewRange(5, 7)),
		model.NewDomain(model.NewRange(3, 4)),
	)

	slices := domains.Slices()

	fmt.Println(len(slices))
	fmt.Println(len(slices[0]))
	fmt.Println(slices[0][0])
	fmt.Println(slices[0][1])
	fmt.Println(slices[0][2])
	fmt.Println(len(slices[1]))
	fmt.Println(slices[1][0])
	fmt.Println(slices[1][1])

	// Output:
	// 2
	// 3
	// 5
	// 6
	// 7
	// 2
	// 3
	// 4
}

func ExampleDomains_Values_first() {
	domains := model.NewDomains(
		model.NewDomain(model.NewRange(5, 7)),
		model.NewDomain(model.NewRange(3, 3)),
	)

	values, singleton := domains.Values()

	fmt.Println(len(values))
	fmt.Println(singleton)

	// Output:
	// 0
	// false
}

func ExampleDomains_Values_second() {
	domains := model.NewDomains(
		model.NewDomain(model.NewRange(5, 5)),
		model.NewDomain(model.NewRange(3, 3)),
	)

	values, singleton := domains.Values()

	fmt.Println(len(values))
	fmt.Println(singleton)
	fmt.Println(values[0])
	fmt.Println(values[1])

	// Output:
	// 2
	// true
	// 5
	// 3
}

func ExampleDomains_First_first() {
	domains := model.NewDomains(
		model.NewDomain(model.NewRange(5, 7)),
		model.NewDomain(model.NewRange(3, 3)),
	)

	fmt.Println(domains.First())

	// Output:
	// 0 true
}

func ExampleDomains_First_second() {
	domains := model.NewDomains(
		model.NewDomain(model.NewRange(5, 5)),
		model.NewDomain(model.NewRange(3, 4)),
	)

	fmt.Println(domains.First())

	// Output:
	// 1 true
}

func ExampleDomains_First_third() {
	domains := model.NewDomains(
		model.NewDomain(model.NewRange(5, 5)),
		model.NewDomain(model.NewRange(3, 3)),
	)

	_, found := domains.First()

	fmt.Println(found)

	// Output:
	// false
}

func ExampleDomains_Largest_first() {
	domains := model.NewDomains(
		model.NewDomain(model.NewRange(5, 7)),
		model.NewDomain(model.NewRange(3, 3)),
	)

	fmt.Println(domains.Largest())

	// Output:
	// 0 true
}

func ExampleDomains_Largest_second() {
	domains := model.NewDomains(
		model.NewDomain(model.NewRange(5, 5)),
		model.NewDomain(model.NewRange(3, 4)),
	)

	fmt.Println(domains.Largest())

	// Output:
	// 1 true
}

func ExampleDomains_Largest_third() {
	domains := model.NewDomains(
		model.NewDomain(model.NewRange(5, 5)),
		model.NewDomain(model.NewRange(3, 3)),
	)

	_, found := domains.Largest()

	fmt.Println(found)

	// Output:
	// false
}

func ExampleDomains_Last_first() {
	domains := model.NewDomains(
		model.NewDomain(model.NewRange(5, 7)),
		model.NewDomain(model.NewRange(3, 3)),
	)

	fmt.Println(domains.Last())

	// Output:
	// 0 true
}

func ExampleDomains_Last_second() {
	domains := model.NewDomains(
		model.NewDomain(model.NewRange(5, 5)),
		model.NewDomain(model.NewRange(3, 4)),
		model.NewDomain(model.NewRange(3, 4)),
	)

	fmt.Println(domains.Last())

	// Output:
	// 2 true
}

func ExampleDomains_Last_third() {
	domains := model.NewDomains(
		model.NewDomain(model.NewRange(5, 5)),
		model.NewDomain(model.NewRange(3, 3)),
	)

	_, found := domains.Last()

	fmt.Println(found)

	// Output:
	// false
}

func ExampleDomains_Maximum_first() {
	domains := model.NewDomains(
		model.NewDomain(model.NewRange(5, 7)),
		model.NewDomain(model.NewRange(3, 3)),
	)

	fmt.Println(domains.Maximum())

	// Output:
	// 0 true
}

func ExampleDomains_Maximum_second() {
	domains := model.NewDomains(
		model.NewDomain(model.NewRange(5, 5)),
		model.NewDomain(model.NewRange(3, 7)),
		model.NewDomain(model.NewRange(3, 4)),
	)

	fmt.Println(domains.Maximum())

	// Output:
	// 1 true
}

func ExampleDomains_Maximum_third() {
	domains := model.NewDomains(
		model.NewDomain(model.NewRange(5, 5)),
		model.NewDomain(model.NewRange(3, 3)),
	)

	_, found := domains.Maximum()

	fmt.Println(found)

	// Output:
	// false
}

func ExampleDomains_Minimum_first() {
	domains := model.NewDomains(
		model.NewDomain(model.NewRange(5, 7)),
		model.NewDomain(model.NewRange(3, 3)),
	)

	fmt.Println(domains.Minimum())

	// Output:
	// 0 true
}

func ExampleDomains_Minimum_second() {
	domains := model.NewDomains(
		model.NewDomain(model.NewRange(5, 5)),
		model.NewDomain(model.NewRange(3, 7)),
		model.NewDomain(model.NewRange(2, 4)),
	)

	fmt.Println(domains.Minimum())

	// Output:
	// 2 true
}

func ExampleDomains_Minimum_third() {
	domains := model.NewDomains(
		model.NewDomain(model.NewRange(5, 5)),
		model.NewDomain(model.NewRange(3, 3)),
	)

	_, found := domains.Minimum()

	fmt.Println(found)

	// Output:
	// false
}

func ExampleDomains_Smallest_first() {
	domains := model.NewDomains(
		model.NewDomain(model.NewRange(5, 7)),
		model.NewDomain(model.NewRange(3, 3)),
	)

	fmt.Println(domains.Smallest())

	// Output:
	// 0 true
}

func ExampleDomains_Smallest_second() {
	domains := model.NewDomains(
		model.NewDomain(model.NewRange(5, 5)),
		model.NewDomain(model.NewRange(3, 7)),
		model.NewDomain(model.NewRange(2, 4)),
	)

	fmt.Println(domains.Smallest())

	// Output:
	// 2 true
}

func ExampleDomains_Smallest_third() {
	domains := model.NewDomains(
		model.NewDomain(model.NewRange(5, 5)),
		model.NewDomain(model.NewRange(3, 3)),
	)

	_, found := domains.Smallest()

	fmt.Println(found)

	// Output:
	// false
}
