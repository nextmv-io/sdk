package model_test

import (
	"fmt"

	"github.com/nextmv-io/sdk/model"
)

func ExampleNewDomain_empty() {
	domain := model.NewDomain()

	fmt.Println(domain.Empty())

	// Output:
	// true
}

func ExampleSingleton() {
	domain := model.Singleton(1)

	fmt.Println(domain.Value())

	// Output:
	// 1 true
}

func ExampleMultiple_first() {
	domain := model.Multiple()

	fmt.Println(domain.Empty())

	// Output:
	// true
}

func ExampleMultiple_second() {
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

/*
func ExampleNewDomain_singleRange_first() {
	domain := model.NewDomain(model.NewRange(1, 0))

	fmt.Println(domain.Empty())

	// Output:
	// true
}
*/

func ExampleNewDomain_singleRange_second() {
	domain := model.NewDomain(model.NewRange(0, 0))

	fmt.Println(domain.Empty())

	fmt.Println(domain.Value())

	fmt.Println(domain.Min())

	fmt.Println(domain.Max())

	iterator := domain.Iterator()
	for iterator.Next() {
		fmt.Println(iterator.Value())
	}
	fmt.Println("done")

	fmt.Println(domain.Contains(-1))

	fmt.Println(domain.Contains(0))

	fmt.Println(domain.Contains(1))

	// Output:
	// false
	// 0 true
	// 0 true
	// 0 true
	// 0
	// done
	// false
	// true
	// false
}

func ExampleNewDomain_singleRange_third() {
	domain := model.NewDomain(model.NewRange(0, 1))

	fmt.Println(domain.Empty())

	_, singleton := domain.Value()

	fmt.Println(singleton)

	fmt.Println(domain.Min())

	fmt.Println(domain.Max())

	// Output:
	// false
	// false
	// 0 true
	// 1 true
}

func ExampleNewDomain_multipleRange_first() {
	domain := model.NewDomain(
		model.NewRange(0, 0),
		model.NewRange(0, 0),
	)

	fmt.Println(domain.Empty())

	fmt.Println(domain.Value())

	fmt.Println(domain.Min())

	fmt.Println(domain.Max())

	// Output:
	// false
	// 0 true
	// 0 true
	// 0 true
}

func ExampleNewDomain_multipleRange_second() {
	domain := model.NewDomain(
		model.NewRange(0, 1),
		model.NewRange(0, 0),
	)

	fmt.Println(domain.Empty())

	_, singleton := domain.Value()

	fmt.Println(singleton)

	fmt.Println(domain.Min())

	fmt.Println(domain.Max())

	// Output:
	// false
	// false
	// 0 true
	// 1 true
}

func ExampleNewDomain_multipleRange_third() {
	domain := model.NewDomain(
		model.NewRange(0, 1),
		model.NewRange(0, 1),
	)

	fmt.Println(domain.Empty())

	_, singleton := domain.Value()

	fmt.Println(singleton)

	fmt.Println(domain.Min())

	fmt.Println(domain.Max())

	// Output:
	// false
	// false
	// 0 true
	// 1 true
}

func ExampleNewDomain_multipleRange_fourth() {
	domain := model.NewDomain(
		model.NewRange(0, 1),
		model.NewRange(0, 2),
	)

	fmt.Println(domain.Empty())

	_, singleton := domain.Value()

	fmt.Println(singleton)

	fmt.Println(domain.Min())

	fmt.Println(domain.Max())

	// Output:
	// false
	// false
	// 0 true
	// 2 true
}

func ExampleNewDomain_multipleRange_fifth() {
	domain := model.NewDomain(
		model.NewRange(0, 2),
		model.NewRange(0, 1),
	)

	fmt.Println(domain.Empty())

	_, singleton := domain.Value()

	fmt.Println(singleton)

	fmt.Println(domain.Min())

	fmt.Println(domain.Max())

	// Output:
	// false
	// false
	// 0 true
	// 2 true
}

func ExampleNewDomain_multipleRange_sixth() {
	domain := model.NewDomain(
		model.NewRange(0, 1),
		model.NewRange(2, 3),
	)

	fmt.Println(domain.Min())

	fmt.Println(domain.Max())

	// Output:
	// 0 true
	// 3 true
}

func ExampleNewDomain_multipleRange_seventh() {
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

func ExampleNewDomain_multipleRange_eighth() {
	domain := model.NewDomain(
		model.NewRange(1, 1),
		model.NewRange(0, 2),
	)

	fmt.Println(domain.Min())

	fmt.Println(domain.Max())

	// Output:
	// 0 true
	// 2 true
}

func ExampleNewDomain_multipleRange_ninth() {
	domain := model.NewDomain(
		model.NewRange(0, 2),
		model.NewRange(1, 1),
	)

	fmt.Println(domain.Min())

	fmt.Println(domain.Max())

	// Output:
	// 0 true
	// 2 true
}

func ExampleDomain_Remove_first() {
	domain := model.NewDomain()

	domain = domain.Remove(0)

	fmt.Println(domain.Empty())

	// Output:
	// true
}

func ExampleDomain_Remove_second() {
	domain := model.NewDomain(model.NewRange(0, 0))

	domain = domain.Remove(0)

	fmt.Println(domain.Empty())

	// Output:
	// true
}

func ExampleDomain_Remove_third() {
	domain := model.NewDomain(model.NewRange(0, 1))

	domain = domain.Remove(0)

	fmt.Println(domain.Min())

	fmt.Println(domain.Max())

	// Output:
	// 1 true
	// 1 true
}

func ExampleDomain_Remove_fourth() {
	domain := model.NewDomain(model.NewRange(0, 1))

	domain = domain.Remove(1)

	fmt.Println(domain.Min())

	fmt.Println(domain.Max())

	// Output:
	// 0 true
	// 0 true
}

func ExampleDomain_Add_first() {
	domain := model.NewDomain(model.NewRange(0, 1))

	domain = domain.Add(1)

	fmt.Println(domain.Min())

	fmt.Println(domain.Max())

	// Output:
	// 0 true
	// 1 true
}

func ExampleDomain_Add_second() {
	domain := model.NewDomain(model.NewRange(0, 1))

	domain = domain.Add(-1)

	fmt.Println(domain.Min())

	fmt.Println(domain.Max())

	// Output:
	// -1 true
	// 1 true
}

func ExampleDomain_Add_thirth() {
	domain := model.NewDomain(model.NewRange(0, 1))

	domain = domain.Add(2)

	fmt.Println(domain.Min())

	fmt.Println(domain.Max())

	// Output:
	// 0 true
	// 2 true
}

func ExampleDomain_Add_fourth() {
	domain := model.NewDomain(model.NewRange(0, 1))

	domain = domain.Add(3)

	fmt.Println(domain.Min())

	fmt.Println(domain.Max())

	fmt.Println(domain.Contains(2))

	// Output:
	// 0 true
	// 3 true
	// false
}

func ExampleDomain_AtLeast_first() {
	domain := model.NewDomain(model.NewRange(0, 1))

	domain = domain.AtLeast(1)

	fmt.Println(domain.Min())

	fmt.Println(domain.Max())

	// Output:
	// 1 true
	// 1 true
}

func ExampleDomain_AtLeast_second() {
	domain := model.NewDomain(model.NewRange(0, 1))

	domain = domain.AtLeast(2)

	fmt.Println(domain.Empty())

	// Output:
	// true
}

func ExampleDomain_AtLeast_third() {
	domain := model.NewDomain(model.NewRange(0, 1))

	domain = domain.AtLeast(0)

	fmt.Println(domain.Min())

	fmt.Println(domain.Max())

	// Output:
	// 0 true
	// 1 true
}

func ExampleDomain_AtLeast_fourth() {
	domain := model.NewDomain(model.NewRange(0, 1))

	domain = domain.AtLeast(-1)

	fmt.Println(domain.Min())

	fmt.Println(domain.Max())

	// Output:
	// 0 true
	// 1 true
}

func ExampleDomain_AtMost_first() {
	domain := model.NewDomain(model.NewRange(0, 1))

	domain = domain.AtMost(0)

	fmt.Println(domain.Min())

	fmt.Println(domain.Max())

	// Output:
	// 0 true
	// 0 true
}

func ExampleDomain_AtMost_second() {
	domain := model.NewDomain(model.NewRange(0, 1))

	domain = domain.AtMost(-1)

	fmt.Println(domain.Empty())

	// Output:
	// true
}

func ExampleDomain_AtMost_third() {
	domain := model.NewDomain(model.NewRange(0, 1))

	domain = domain.AtMost(1)

	fmt.Println(domain.Min())

	fmt.Println(domain.Max())

	// Output:
	// 0 true
	// 1 true
}

func ExampleDomain_AtMost_fourth() {
	domain := model.NewDomain(model.NewRange(0, 1))

	domain = domain.AtMost(2)

	fmt.Println(domain.Min())

	fmt.Println(domain.Max())

	// Output:
	// 0 true
	// 1 true
}

func ExampleDomain_Contains() {
	domain := model.NewDomain(model.NewRange(0, 1))

	fmt.Println(domain.Contains(-1))

	fmt.Println(domain.Contains(0))

	fmt.Println(domain.Contains(1))

	fmt.Println(domain.Contains(2))

	// Output:
	// false
	// true
	// true
	// false
}

func ExampleDomain_Len_first() {
	domain := model.NewDomain()

	fmt.Println(domain.Len())

	// Output:
	// 0
}

func ExampleDomain_Len_second() {
	domain := model.NewDomain(model.NewRange(0, 1))

	fmt.Println(domain.Len())

	// Output:
	// 2
}

func ExampleDomain_Len_third() {
	domain := model.NewDomain(
		model.NewRange(0, 1),
		model.NewRange(4, 4),
	)

	fmt.Println(domain.Len())

	// Output:
	// 3
}

func ExampleDomain_Slice_first() {
	domain := model.NewDomain(model.NewRange(0, 1))

	fmt.Println(len(domain.Slice()))

	// Output:
	// 2
}

func ExampleDomain_Slice_second() {
	domain := model.NewDomain()

	fmt.Println(len(domain.Slice()))

	// Output:
	// 0
}

func ExampleDomain_Slice_third() {
	domain := model.NewDomain(
		model.NewRange(0, 1),
		model.NewRange(4, 4),
	)

	fmt.Println(len(domain.Slice()))

	// Output:
	// 3
}

func ExampleDomain_Value_first() {
	domain := model.NewDomain()

	_, singleton := domain.Value()

	fmt.Println(singleton)

	// Output:
	// false
}

func ExampleDomain_Value_second() {
	domain := model.NewDomain(model.NewRange(0, 0))

	fmt.Println(domain.Value())

	// Output:
	// 0 true
}

func ExampleDomain_Value_third() {
	domain := model.NewDomain(model.NewRange(0, 1))

	_, singleton := domain.Value()

	fmt.Println(singleton)

	// Output:
	// false
}
