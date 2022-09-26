package mip_test

import (
	"fmt"
	"testing"

	"github.com/nextmv-io/sdk/mip"
)

func ExampleVar_continuous() {
	model := mip.NewModel()

	v := model.NewContinuousVar(-1.0, 1.0)

	fmt.Println(v.LowerBound())
	fmt.Println(v.UpperBound())
	fmt.Println(v.IsContinuous())
	fmt.Println(v.IsInteger())
	fmt.Println(v.IsBinary())
	fmt.Println(v)
	// Output:
	// -1
	// 1
	// true
	// false
	// false
	// C0
}

func ExampleVar_integer() {
	model := mip.NewModel()

	v := model.NewIntegerVar(-1, 1)

	fmt.Println(v.LowerBound())
	fmt.Println(v.UpperBound())
	fmt.Println(v.IsContinuous())
	fmt.Println(v.IsInteger())
	fmt.Println(v.IsBinary())
	fmt.Println(v)
	v.SetName("v")
	fmt.Println(v)
	// Output:
	// -1
	// 1
	// false
	// true
	// false
	// I0
	// v
}

func ExampleVar_binary() {
	model := mip.NewModel()

	v := model.NewBinaryVar()

	fmt.Println(v.LowerBound())
	fmt.Println(v.UpperBound())
	fmt.Println(v.IsContinuous())
	fmt.Println(v.IsInteger())
	fmt.Println(v.IsBinary())
	fmt.Println(v)
	// Output:
	// 0
	// 1
	// false
	// true
	// true
	// B0
}

func ExampleVar_vars() {
	model := mip.NewModel()

	v0 := model.NewBinaryVar()
	v1 := model.NewIntegerVar(-1, 1)
	v2 := model.NewContinuousVar(-1.0, 1.0)

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
		model.NewBinaryVar()
	}
}

func BenchmarkNewContinuousVar(b *testing.B) {
	model := mip.NewModel()
	for i := 0; i < b.N; i++ {
		model.NewContinuousVar(-1.0, 1.0)
	}
}

func BenchmarkNewIntegerVar(b *testing.B) {
	model := mip.NewModel()
	for i := 0; i < b.N; i++ {
		model.NewIntegerVar(-1, 1)
	}
}
