package measure_test

import (
	"encoding/json"
	"testing"

	"github.com/nextmv-io/sdk/measure"
	"github.com/nextmv-io/sdk/model"
)

func TestLocation(t *testing.T) {
	matrix := [][]float64{{10, 20}, {5, 2}}
	matrixM := measure.Matrix(matrix)
	costs := []float64{10, 5}
	dg := measure.DurationGroup{
		Group:    model.NewDomain(model.NewRange(0, 0)),
		Duration: 300,
	}
	durationGroups := []measure.DurationGroup{
		dg,
	}
	m, _ := measure.Location(matrixM, costs, durationGroups)

	for i, row := range matrix {
		for j, w := range row {
			v := m.Cost(i, j)
			dg := 0
			if j == 0 && i != 0 {
				dg = 300
			}
			if v != w+costs[j]+float64(dg) {
				t.Errorf("got %f; want %f", v, w)
			}
		}
	}

	b, err := json.Marshal(m)
	if err != nil {
		t.Errorf("got %+v; want nil", err)
	}
	w := `{"costs":[10,5],"measure":{"matrix":[[10,20],[5,2]],` +
		`"type":"matrix"},"type":"location"}`
	if v := string(b); v != w {
		t.Errorf("got %q; want %q", v, w)
	}
}
