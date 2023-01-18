package measure_test

import (
	"encoding/json"
	"testing"

	"github.com/nextmv-io/sdk/measure"
)

func TestMatrix(t *testing.T) {
	matrix := [][]float64{{10, 20}, {5, 2}}
	m := measure.Matrix(matrix)

	for i, row := range matrix {
		for j, w := range row {
			v := m.Cost(i, j)
			if v != w {
				t.Errorf("got %f; want %f", v, w)
			}
		}
	}

	b, err := json.Marshal(m)
	if err != nil {
		t.Errorf("got %+v; want nil", err)
	}
	w := `{"matrix":[[10,20],[5,2]],"type":"matrix"}`
	if v := string(b); v != w {
		t.Errorf("got %q; want %q", v, w)
	}
}
