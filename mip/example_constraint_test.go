package mip_test

import (
	"fmt"

	"github.com/nextmv-io/sdk/mip"
)

func ExampleConstraint_greaterThanEqual() {
	definition := mip.NewDefinition()

	c, _ := definition.AddConstraint(mip.GreaterThanOrEqual, 1.0)

	fmt.Println(c.Sense())
	fmt.Println(c.RightHandSide())
	// Output:
	// 2
	// 1
}

func ExampleConstraint_equal() {
	definition := mip.NewDefinition()

	c, _ := definition.AddConstraint(mip.Equal, 1.0)

	fmt.Println(c.Sense())
	fmt.Println(c.RightHandSide())
	// Output:
	// 1
	// 1
}

func ExampleConstraint_lessThanOrEqual() {
	definition := mip.NewDefinition()

	c, _ := definition.AddConstraint(mip.LessThanOrEqual, 1.0)

	fmt.Println(c.Sense())
	fmt.Println(c.RightHandSide())
	// Output:
	// 0
	// 1
}

func ExampleConstraint_terms() {
	definition := mip.NewDefinition()

	v, _ := definition.AddBinaryVariable()
	c, _ := definition.AddConstraint(mip.Equal, 1.0)

	t1 := c.AddTerm(1.0, v)
	t2 := c.AddTerm(2.0, v)

	fmt.Println(t1.Variable().Index())
	fmt.Println(t1.Coefficient())
	fmt.Println(t2.Coefficient())
	fmt.Println(len(c.Terms()))
	fmt.Println(c.Terms()[0].Coefficient())
	// Output:
	// 0
	// 1
	// 2
	// 1
	// 3
}
