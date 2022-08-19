package mip_test

import (
	"fmt"
	"testing"

	"github.com/nextmv-io/sdk/mip"
)

func ExampleObjective_sense() {
	model := mip.NewModel()

	model.Objective().SetMaximize()

	fmt.Println(model.Objective().IsMaximize())

	model.Objective().SetMinimize()

	fmt.Println(model.Objective().IsMaximize())
	// Output:
	// true
	// false
}

func ExampleObjective_terms() {
	model := mip.NewModel()

	v1, _ := model.NewBinaryVariable()
	v2, _ := model.NewBinaryVariable()

	fmt.Println(len(model.Objective().Terms()))

	t1 := model.Objective().NewTerm(2.0, v1)
	t2 := model.Objective().NewTerm(1.0, v1)
	t3 := model.Objective().NewTerm(3.0, v2)

	fmt.Println(t1.Variable().Index())
	fmt.Println(t1.Coefficient())

	fmt.Println(t2.Variable().Index())
	fmt.Println(t2.Coefficient())

	fmt.Println(t3.Variable().Index())
	fmt.Println(t3.Coefficient())

	fmt.Println(len(model.Objective().Terms()))
	fmt.Println(model.Objective().Terms()[0].Coefficient())
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

func benchmarkObjectiveNewTerms(nrTerms int, b *testing.B) {
	model := mip.NewModel()
	v, _ := model.NewContinuousVariable(1.0, 2.0)

	for i := 0; i < b.N; i++ {
		for i := 0; i < nrTerms; i++ {
			model.Objective().NewTerm(1.0, v)
		}
	}
}

func BenchmarkObjectiveNewTerms1(b *testing.B) {
	benchmarkObjectiveNewTerms(1, b)
}

func BenchmarkObjectiveNewTerms2(b *testing.B) {
	benchmarkObjectiveNewTerms(2, b)
}

func BenchmarkObjectiveNewTerms4(b *testing.B) {
	benchmarkObjectiveNewTerms(4, b)
}

func BenchmarkObjectiveNewTerms8(b *testing.B) {
	benchmarkObjectiveNewTerms(8, b)
}

func BenchmarkObjectiveNewTerms16(b *testing.B) {
	benchmarkObjectiveNewTerms(16, b)
}

func BenchmarkObjectiveNewTerms32(b *testing.B) {
	benchmarkObjectiveNewTerms(32, b)
}
