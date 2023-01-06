package measure_test

import (
	"encoding/json"
	"testing"

	"github.com/nextmv-io/sdk/measure"
)

func TestSparse(t *testing.T) {
	matrix := map[int]map[int]float64{
		10: {20: 3.14, 24: 100},
		12: {10: 2},
	}
	m := measure.DebugSparse(measure.Constant(-1), matrix)

	for i, row := range matrix {
		for j, w := range row {
			v := m.Cost(i, j)
			if v != w {
				t.Errorf("got %f; want %f", v, w)
			}
		}
	}

	if v := m.Cost(24, 12); v != -1 {
		t.Errorf("got %f; want -1", v)
	}

	b, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		t.Errorf("got %+v; want nil", err)
	}
	w := `{
	"arcs": {
		"10": {
			"20": 3.14,
			"24": 100
		},
		"12": {
			"10": 2
		}
	},
	"counts": {
		"hit": 3,
		"miss": 1
	},
	"measure": {
		"constant": -1,
		"type": "constant"
	},
	"type": "sparse"
}`
	if v := string(b); v != w {
		t.Errorf("got %s\nwant %s\n", v, w)
	}
}
