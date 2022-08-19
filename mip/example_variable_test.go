package mip_test

import (
	"fmt"
	"testing"

	"github.com/nextmv-io/sdk/mip"
)

func ExampleVariable_continuous() {
	model := mip.NewModel()

	v, _ := model.NewContinuousVariable(-1.0, 1.0)

	fmt.Println(v.LowerBound())
	fmt.Println(v.UpperBound())
	// Output:
	// -1
	// 1
}

func ExampleVariable_integer() {
	model := mip.NewModel()

	v, _ := model.NewIntegerVariable(-1, 1)

	fmt.Println(v.LowerBound())
	fmt.Println(v.UpperBound())
	// Output:
	// -1
	// 1
}

func ExampleVariable_binary() {
	model := mip.NewModel()

	v, _ := model.NewBinaryVariable()

	fmt.Println(v.LowerBound())
	fmt.Println(v.UpperBound())
	// Output:
	// 0
	// 1
}

func ExampleVariable_variables() {
	model := mip.NewModel()

	v0, _ := model.NewBinaryVariable()
	v1, _ := model.NewIntegerVariable(-1, 1)
	v2, _ := model.NewContinuousVariable(-1.0, 1.0)

	fmt.Println(v0.Index())
	fmt.Println(v1.Index())
	fmt.Println(v2.Index())
	fmt.Println(len(model.Variables()))
	// Output:
	// 0
	// 1
	// 2
	// 3
}

func BenchmarkNewBinaryVariable(b *testing.B) {
	model := mip.NewModel()
	for i := 0; i < b.N; i++ {
		model.NewBinaryVariable()
	}
}

func BenchmarkNewContinuousVariable(b *testing.B) {
	model := mip.NewModel()
	for i := 0; i < b.N; i++ {
		model.NewContinuousVariable(-1.0, 1.0)
	}
}

func BenchmarkNewIntegerVariable(b *testing.B) {
	model := mip.NewModel()
	for i := 0; i < b.N; i++ {
		model.NewIntegerVariable(-1, 1)
	}
}
