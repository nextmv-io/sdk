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

	_, err := model.NewBinaryVariable()
	if err != nil {
		panic(err)
	}
	_, err = model.NewContinuousVariable(1.0, 2.0)
	if err != nil {
		panic(err)
	}
	_, err = model.NewBinaryVariable()
	if err != nil {
		panic(err)
	}

	_, err = model.NewConstraint(mip.Equal, 0.0)
	if err != nil {
		panic(err)
	}

	fmt.Println(len(model.Variables()))
	fmt.Println(len(model.Constraints()))
	// Output:
	// 3
	// 1
}
