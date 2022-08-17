package mip_test

import (
	"fmt"
	"testing"

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

func benchmarkAddConstraintAddTerms(nrTerms int, b *testing.B) {
	definition := mip.NewDefinition()
	v, _ := definition.AddContinuousVariable(1.0, 2.0)

	for i := 0; i < b.N; i++ {
		c, _ := definition.AddConstraint(mip.Equal, 1.0)
		for i := 0; i < nrTerms; i++ {
			c.AddTerm(1.0, v)
		}
	}
}

func BenchmarkAddConstraintAddTerms0(b *testing.B) {
	benchmarkAddConstraintAddTerms(0, b)
}

func BenchmarkAddConstraintAddTerms1(b *testing.B) {
	benchmarkAddConstraintAddTerms(1, b)
}

func BenchmarkAddConstraintAddTerms2(b *testing.B) {
	benchmarkAddConstraintAddTerms(2, b)
}

func BenchmarkAddConstraintAddTerms4(b *testing.B) {
	benchmarkAddConstraintAddTerms(4, b)
}

func BenchmarkAddConstraintAddTerms8(b *testing.B) {
	benchmarkAddConstraintAddTerms(8, b)
}

func BenchmarkAddConstraintAddTerms16(b *testing.B) {
	benchmarkAddConstraintAddTerms(16, b)
}

func BenchmarkAddConstraintAddTerms32(b *testing.B) {
	benchmarkAddConstraintAddTerms(32, b)
}
