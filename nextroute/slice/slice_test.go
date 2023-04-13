package slice_test

import (
	"testing"

	"github.com/nextmv-io/sdk/nextroute/slice"
)

func BenchmarkFilter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		numberOfValues := 100
		values := make([]int, numberOfValues)
		for i := 0; i < numberOfValues; i++ {
			values[i] = i
		}
		_ = slice.Filter(values, func(vehicle int) bool {
			return vehicle%2 == 0
		})
	}
}

func TestDefensiveCopy(t *testing.T) {
	numberOfValues := 100
	values := make([]int, numberOfValues)
	for i := 0; i < numberOfValues; i++ {
		values[i] = i
	}
	copiedValues := slice.DefensiveCopy(values)
	for i := 0; i < numberOfValues; i++ {
		if values[i] != copiedValues[i] {
			t.Errorf("Expected %v, got %v", values[i], copiedValues[i])
		}
	}
}