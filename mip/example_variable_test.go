package mip_test

import (
	"fmt"
	"testing"

	"github.com/nextmv-io/sdk/mip"
)

func ExampleVariable_continuous() {
	definition := mip.NewDefinition()

	v, _ := definition.NewContinuousVariable(-1.0, 1.0)

	fmt.Println(v.LowerBound())
	fmt.Println(v.UpperBound())
	// Output:
	// -1
	// 1
}

func ExampleVariable_integer() {
	definition := mip.NewDefinition()

	v, _ := definition.NewIntegerVariable(-1, 1)

	fmt.Println(v.LowerBound())
	fmt.Println(v.UpperBound())
	// Output:
	// -1
	// 1
}

func ExampleVariable_binary() {
	definition := mip.NewDefinition()

	v, _ := definition.NewBinaryVariable()

	fmt.Println(v.LowerBound())
	fmt.Println(v.UpperBound())
	// Output:
	// 0
	// 1
}

func ExampleVariable_variables() {
	definition := mip.NewDefinition()

	v0, _ := definition.NewBinaryVariable()
	v1, _ := definition.NewIntegerVariable(-1, 1)
	v2, _ := definition.NewContinuousVariable(-1.0, 1.0)

	fmt.Println(v0.Index())
	fmt.Println(v1.Index())
	fmt.Println(v2.Index())
	fmt.Println(len(definition.Variables()))
	// Output:
	// 0
	// 1
	// 2
	// 3
}

func BenchmarkNewBinaryVariable(b *testing.B) {
	definition := mip.NewDefinition()
	for i := 0; i < b.N; i++ {
		definition.NewBinaryVariable()
	}
}

func BenchmarkNewContinuousVariable(b *testing.B) {
	definition := mip.NewDefinition()
	for i := 0; i < b.N; i++ {
		definition.NewContinuousVariable(-1.0, 1.0)
	}
}

func BenchmarkNewIntegerVariable(b *testing.B) {
	definition := mip.NewDefinition()
	for i := 0; i < b.N; i++ {
		definition.NewIntegerVariable(-1, 1)
	}
}
