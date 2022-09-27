package mip_test

import (
	"fmt"

	"github.com/nextmv-io/sdk/mip"
)

func ExampleModel_empty() {
	model := mip.NewModel()

	fmt.Println(len(model.Constraints()))
	fmt.Println(len(model.Vars()))
	fmt.Println(len(model.Objective().Terms()))
	fmt.Println(model.Objective().IsMaximize())
	fmt.Println(model)
	// Output:
	// 0
	// 0
	// 0
	// false
	// minimize
}

func ExampleModel_queries() {
	model := mip.NewModel()

	model.NewBool()
	model.NewFloat(1.0, 2.0)
	model.NewBool()
	model.NewConstraint(mip.Equal, 0.0)

	fmt.Println(len(model.Vars()))
	fmt.Println(len(model.Constraints()))
	fmt.Println(model)
	// Output:
	// 3
	// 1
	// minimize
	//       0: = 0
	//       0: B0 [0, 1]
	//       1: C1 [1, 2]
	//       2: B2 [0, 1]
}
