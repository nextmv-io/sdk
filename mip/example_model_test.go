package mip_test

import (
	"fmt"

	"github.com/nextmv-io/sdk/mip"
)

func ExampleModel_empty() {
	model := mip.NewModel()

	fmt.Println(len(model.Constraints()))
	fmt.Println(len(model.Variables()))
	fmt.Println(len(model.Objective().Terms()))
	fmt.Println(model.Objective().IsMaximize())
	// Output:
	// 0
	// 0
	// 0
	// false
}

func ExampleModel_queries() {
	model := mip.NewModel()

	model.NewBinaryVariable()
	model.NewContinuousVariable(1.0, 2.0)
	model.NewBinaryVariable()

	model.NewConstraint(mip.Equal, 0.0)

	fmt.Println(len(model.Variables()))
	fmt.Println(len(model.Constraints()))
	// Output:
	// 3
	// 1
}
