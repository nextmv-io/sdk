package mip_test

import (
	"fmt"
	"testing"

	"github.com/nextmv-io/sdk/mip"
)

func ExampleVar_continuous() {
	model := mip.NewModel()

	v, err := model.NewContinuousVar(-1.0, 1.0)
	if err != nil {
		panic(err)
	}

	fmt.Println(v.LowerBound())
	fmt.Println(v.UpperBound())
	fmt.Println(v.IsContinuous())
	fmt.Println(v.IsInteger())
	fmt.Println(v.IsBinary())
	// Output:
	// -1
	// 1
	// true
	// false
	// false
}

func ExampleVar_integer() {
	model := mip.NewModel()

	v, err := model.NewIntegerVar(-1, 1)
	if err != nil {
		panic(err)
	}

	fmt.Println(v.LowerBound())
	fmt.Println(v.UpperBound())
	fmt.Println(v.IsContinuous())
	fmt.Println(v.IsInteger())
	fmt.Println(v.IsBinary())
	// Output:
	// -1
	// 1
	// false
	// true
	// false
}

func ExampleVar_binary() {
	model := mip.NewModel()

	v, err := model.NewBinaryVar()
	if err != nil {
		panic(err)
	}

	fmt.Println(v.LowerBound())
	fmt.Println(v.UpperBound())
	fmt.Println(v.IsContinuous())
	fmt.Println(v.IsInteger())
	fmt.Println(v.IsBinary())
	// Output:
	// 0
	// 1
	// false
	// true
	// true
}

func ExampleVar_vars() {
	model := mip.NewModel()

	v0, err := model.NewBinaryVar()
	if err != nil {
		panic(err)
	}
	v1, err := model.NewIntegerVar(-1, 1)
	if err != nil {
		panic(err)
	}
	v2, err := model.NewContinuousVar(-1.0, 1.0)
	if err != nil {
		panic(err)
	}

	fmt.Println(v0.Index())
	fmt.Println(v1.Index())
	fmt.Println(v2.Index())
	fmt.Println(len(model.Vars()))
	// Output:
	// 0
	// 1
	// 2
	// 3
}

func BenchmarkNewBinaryVar(b *testing.B) {
	model := mip.NewModel()
	for i := 0; i < b.N; i++ {
		_, err := model.NewBinaryVar()
		if err != nil {
			b.Errorf(err.Error())
		}
	}
}

func BenchmarkNewContinuousVar(b *testing.B) {
	model := mip.NewModel()
	for i := 0; i < b.N; i++ {
		_, err := model.NewContinuousVar(-1.0, 1.0)
		if err != nil {
			b.Errorf(err.Error())
		}
	}
}

func BenchmarkNewIntegerVar(b *testing.B) {
	model := mip.NewModel()
	for i := 0; i < b.N; i++ {
		_, err := model.NewIntegerVar(-1, 1)
		if err != nil {
			b.Errorf(err.Error())
		}
	}
}
