package mip_test

import (
	"fmt"

	"github.com/nextmv-io/sdk/mip"
)

func ExampleDefinition_empty() {
	definition := mip.NewDefinition()

	fmt.Println(len(definition.Constraints()))
	fmt.Println(len(definition.Variables()))
	fmt.Println(len(definition.Objective().Terms()))
	fmt.Println(definition.Objective().IsMaximize())
	// Output:
	// 0
	// 0
	// 0
	// false
}

func ExampleDefinition_queries() {
	definition := mip.NewDefinition()

	definition.NewBinaryVariable()
	definition.NewContinuousVariable(1.0, 2.0)
	definition.NewBinaryVariable()

	definition.NewConstraint(mip.Equal, 0.0)

	fmt.Println(len(definition.Variables()))
	fmt.Println(len(definition.Constraints()))
	// Output:
	// 3
	// 1
}
