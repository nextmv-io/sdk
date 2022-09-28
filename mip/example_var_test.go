package mip_test

import (
	"fmt"
	"testing"

	"github.com/nextmv-io/sdk/mip"
)

func ExampleVar_float() {
	model := mip.NewModel()

	v := model.NewFloat(-1.0, 1.0)

	fmt.Println(v.LowerBound())
	fmt.Println(v.UpperBound())
	fmt.Println(v.IsFloat())
	fmt.Println(v.IsInt())
	fmt.Println(v.IsBool())
	fmt.Println(v)
	// Output:
	// -1
	// 1
	// true
	// false
	// false
	// F0
}

func ExampleVar_int() {
	model := mip.NewModel()

	v := model.NewInt(-1, 1)

	fmt.Println(v.LowerBound())
	fmt.Println(v.UpperBound())
	fmt.Println(v.IsFloat())
	fmt.Println(v.IsInt())
	fmt.Println(v.IsBool())
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

func ExampleVar_bool() {
	model := mip.NewModel()

	v := model.NewBool()

	fmt.Println(v.LowerBound())
	fmt.Println(v.UpperBound())
	fmt.Println(v.IsFloat())
	fmt.Println(v.IsInt())
	fmt.Println(v.IsBool())
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

	v0 := model.NewBool()
	v1 := model.NewInt(-1, 1)
	v2 := model.NewFloat(-1.0, 1.0)

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

func BenchmarkNewBool(b *testing.B) {
	model := mip.NewModel()
	for i := 0; i < b.N; i++ {
		model.NewBool()
	}
}

func BenchmarkNewFloat(b *testing.B) {
	model := mip.NewModel()
	for i := 0; i < b.N; i++ {
		model.NewFloat(-1.0, 1.0)
	}
}

func BenchmarkNewInt(b *testing.B) {
	model := mip.NewModel()
	for i := 0; i < b.N; i++ {
		model.NewInt(-1, 1)
	}
}
