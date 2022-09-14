package mip_test

import (
	"fmt"
	"testing"

	"github.com/nextmv-io/sdk/mip"
)

func ExampleConstraint_greaterThanEqual() {
	model := mip.NewModel()

	c, _ := model.NewConstraint(mip.GreaterThanOrEqual, 1.0)

	fmt.Println(c.Sense())
	fmt.Println(c.RightHandSide())
	fmt.Println(c)
	// Output:
	// 2
	// 1
	// >=  1
}

func ExampleConstraint_equal() {
	model := mip.NewModel()

	c, _ := model.NewConstraint(mip.Equal, 1.0)

	fmt.Println(c.Sense())
	fmt.Println(c.RightHandSide())
	fmt.Println(c)
	// Output:
	// 1
	// 1
	// =  1
}

func ExampleConstraint_lessThanOrEqual() {
	model := mip.NewModel()

	c, _ := model.NewConstraint(mip.LessThanOrEqual, 1.0)

	fmt.Println(c.Sense())
	fmt.Println(c.RightHandSide())
	fmt.Println(c)
	// Output:
	// 0
	// 1
	// <=  1
}

func ExampleConstraint_terms() {
	model := mip.NewModel()

	v, _ := model.NewBinaryVar()
	c, _ := model.NewConstraint(mip.Equal, 1.0)

	t1 := c.NewTerm(1.0, v)
	t2 := c.NewTerm(2.0, v)

	fmt.Println(t1.Var().Index())
	fmt.Println(t1.Coefficient())
	fmt.Println(t2.Coefficient())
	fmt.Println(len(c.Terms()))
	fmt.Println(c.Terms()[0].Coefficient())
	fmt.Println(c)
	// Output:
	// 0
	// 1
	// 2
	// 1
	// 3
	// 3 B0 = 1
}

func benchmarkNewConstraintNewTerms(nrTerms int, b *testing.B) {
	model := mip.NewModel()
	v, _ := model.NewContinuousVar(1.0, 2.0)

	for i := 0; i < b.N; i++ {
		c, _ := model.NewConstraint(mip.Equal, 1.0)
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
