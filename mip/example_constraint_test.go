package mip_test

import (
	"fmt"
	"testing"

	"github.com/nextmv-io/sdk/mip"
)

func ExampleConstraint_greaterThanEqual() {
	definition := mip.NewDefinition()

	c, _ := definition.NewConstraint(mip.GreaterThanOrEqual, 1.0)

	fmt.Println(c.Sense())
	fmt.Println(c.RightHandSide())
	// Output:
	// 2
	// 1
}

func ExampleConstraint_equal() {
	definition := mip.NewDefinition()

	c, _ := definition.NewConstraint(mip.Equal, 1.0)

	fmt.Println(c.Sense())
	fmt.Println(c.RightHandSide())
	// Output:
	// 1
	// 1
}

func ExampleConstraint_lessThanOrEqual() {
	definition := mip.NewDefinition()

	c, _ := definition.NewConstraint(mip.LessThanOrEqual, 1.0)

	fmt.Println(c.Sense())
	fmt.Println(c.RightHandSide())
	// Output:
	// 0
	// 1
}

func ExampleConstraint_terms() {
	definition := mip.NewDefinition()

	v, _ := definition.NewBinaryVariable()
	c, _ := definition.NewConstraint(mip.Equal, 1.0)

	t1 := c.NewTerm(1.0, v)
	t2 := c.NewTerm(2.0, v)

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

func benchmarkNewConstraintNewTerms(nrTerms int, b *testing.B) {
	definition := mip.NewDefinition()
	v, _ := definition.NewContinuousVariable(1.0, 2.0)

	for i := 0; i < b.N; i++ {
		c, _ := definition.NewConstraint(mip.Equal, 1.0)
		for i := 0; i < nrTerms; i++ {
			c.NewTerm(1.0, v)
		}
	}
}

func BenchmarkNewConstraintNewTerms0(b *testing.B) {
	benchmarkNewConstraintNewTerms(0, b)
}

func BenchmarkNewConstraintNewTerms1(b *testing.B) {
	benchmarkNewConstraintNewTerms(1, b)
}

func BenchmarkNewConstraintNewTerms2(b *testing.B) {
	benchmarkNewConstraintNewTerms(2, b)
}

func BenchmarkNewConstraintNewTerms4(b *testing.B) {
	benchmarkNewConstraintNewTerms(4, b)
}

func BenchmarkNewConstraintNewTerms8(b *testing.B) {
	benchmarkNewConstraintNewTerms(8, b)
}

func BenchmarkNewConstraintNewTerms16(b *testing.B) {
	benchmarkNewConstraintNewTerms(16, b)
}

func BenchmarkNewConstraintNewTerms32(b *testing.B) {
	benchmarkNewConstraintNewTerms(32, b)
}
