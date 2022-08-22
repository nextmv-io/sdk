package mip_test

import (
	"fmt"
	"testing"

	"github.com/nextmv-io/sdk/mip"
)

func ExampleVariable_continuous() {
	model := mip.NewModel()

	v, err := model.NewContinuousVariable(-1.0, 1.0)
	if err != nil {
		panic(err)
	}

	fmt.Println(v.LowerBound())
	fmt.Println(v.UpperBound())
	// Output:
	// -1
	// 1
}

func ExampleVariable_integer() {
	model := mip.NewModel()

	v, err := model.NewIntegerVariable(-1, 1)
	if err != nil {
		panic(err)
	}

	fmt.Println(v.LowerBound())
	fmt.Println(v.UpperBound())
	// Output:
	// -1
	// 1
}

func ExampleVariable_binary() {
	model := mip.NewModel()

	v, err := model.NewBinaryVariable()
	if err != nil {
		panic(err)
	}

	fmt.Println(v.LowerBound())
	fmt.Println(v.UpperBound())
	// Output:
	// 0
	// 1
}

func ExampleVariable_variables() {
	model := mip.NewModel()

	v0, err := model.NewBinaryVariable()
	if err != nil {
		panic(err)
	}
	v1, err := model.NewIntegerVariable(-1, 1)
	if err != nil {
		panic(err)
	}
	v2, err := model.NewContinuousVariable(-1.0, 1.0)
	if err != nil {
		panic(err)
	}

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
		_, err := model.NewBinaryVariable()
		if err != nil {
			b.Errorf(err.Error())
		}
	}
}

func BenchmarkNewContinuousVariable(b *testing.B) {
	model := mip.NewModel()
	for i := 0; i < b.N; i++ {
		_, err := model.NewContinuousVariable(-1.0, 1.0)
		if err != nil {
			b.Errorf(err.Error())
		}
	}
}

func BenchmarkNewIntegerVariable(b *testing.B) {
	model := mip.NewModel()
	for i := 0; i < b.N; i++ {
		_, err := model.NewIntegerVariable(-1, 1)
		if err != nil {
			b.Errorf(err.Error())
		}
	}
}
