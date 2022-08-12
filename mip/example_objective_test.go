package mip_test

import (
	"fmt"

	"github.com/nextmv-io/sdk/mip"
)

func ExampleObjective_sense() {
	definition := mip.NewDefinition()

	definition.Objective().SetMaximize()

	fmt.Println(definition.Objective().IsMaximize())

	definition.Objective().SetMinimize()

	fmt.Println(definition.Objective().IsMaximize())
	// Output:
	// true
	// false
}

func ExampleObjective_terms() {
	definition := mip.NewDefinition()

	v1, _ := definition.AddBinaryVariable()
	v2, _ := definition.AddBinaryVariable()

	fmt.Println(len(definition.Objective().Terms()))

	t1 := definition.Objective().AddTerm(2.0, v1)
	t2 := definition.Objective().AddTerm(1.0, v1)
	t3 := definition.Objective().AddTerm(3.0, v2)

	fmt.Println(t1.Variable().Index())
	fmt.Println(t1.Coefficient())

	fmt.Println(t2.Variable().Index())
	fmt.Println(t2.Coefficient())

	fmt.Println(t3.Variable().Index())
	fmt.Println(t3.Coefficient())

	fmt.Println(len(definition.Objective().Terms()))
	fmt.Println(definition.Objective().Terms()[0].Coefficient())
	// Output:
	// 0
	// 0
	// 2
	// 0
	// 1
	// 1
	// 3
	// 2
	// 3
}
