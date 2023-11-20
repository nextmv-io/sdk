package measure_test

import (
	"testing"

	"github.com/nextmv-io/sdk/measure"
)

func TestUnique(t *testing.T) {
	uniqueMatrix := measure.Matrix([][]float64{
		{0, 1, 2},
		{1, 0, 3},
		{2, 3, 0},
	})
	references := []int{1, 2, 0, 1, 2, 0}
	m := measure.Unique(uniqueMatrix, references)
	c := m.Cost(2, 5)
	if c != 0 {
		t.Errorf("expected cost to be 0 but was %f", c)
	}
	c = m.Cost(0, 4)
	if c != 3 {
		t.Errorf("expected cost to be 1 but was %f", c)
	}
}
