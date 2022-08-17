package mip_test

import (
	"fmt"
	"testing"

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

func benchmarkObjectiveAddTerms(nrTerms int, b *testing.B) {
	definition := mip.NewDefinition()
	v, _ := definition.AddContinuousVariable(1.0, 2.0)

	for i := 0; i < b.N; i++ {
		for i := 0; i < nrTerms; i++ {
			definition.Objective().AddTerm(1.0, v)
		}
	}
}

func BenchmarkObjectiveAddTerms1(b *testing.B) {
	benchmarkObjectiveAddTerms(1, b)
}

func BenchmarkObjectiveAddTerms2(b *testing.B) {
	benchmarkObjectiveAddTerms(2, b)
}

func BenchmarkObjectiveAddTerms4(b *testing.B) {
	benchmarkObjectiveAddTerms(4, b)
}

func BenchmarkObjectiveAddTerms8(b *testing.B) {
	benchmarkObjectiveAddTerms(8, b)
}

func BenchmarkObjectiveAddTerms16(b *testing.B) {
	benchmarkObjectiveAddTerms(16, b)
}

func BenchmarkObjectiveAddTerms32(b *testing.B) {
	benchmarkObjectiveAddTerms(32, b)
}
