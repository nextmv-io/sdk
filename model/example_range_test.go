package model_test

import (
	"fmt"

	"github.com/nextmv-io/sdk/model"
)

func ExampleNewRange() {
	r := model.NewRange(1, 0)
	fmt.Println(r.Min())
	fmt.Println(r.Max())

	r = model.NewRange(0, 0)
	fmt.Println(r.Min())
	fmt.Println(r.Max())

	r = model.NewRange(0, 1)
	fmt.Println(r.Min())
	fmt.Println(r.Max())

	r = model.NewRange(model.MinInt, model.MaxInt)
	fmt.Println(r.Min())
	fmt.Println(r.Max())

	// Output:
	// 0
	// 1
	// 0
	// 0
	// 0
	// 1
	// -9223372036854775808
	// 9223372036854775807
}
