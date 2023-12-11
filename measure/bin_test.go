package measure_test

import (
	"encoding/json"
	"testing"

	"github.com/nextmv-io/sdk/measure"
)

func TestBin(t *testing.T) {
	points := []measure.Point{{1, 2}, {4, 6}, {6, 6}}
	measures := []measure.ByIndex{
		measure.Constant(0),
		measure.Indexed(measure.EuclideanByPoint(), points),
		measure.Constant(10),
		measure.Constant(30),
	}

	m := measure.Bin(
		measures,
		func(from, to int) int {
			if from == 0 && to == 1 { //nolint:gocritic
				return 1
			} else if from == 1 && to == 2 {
				return 2
			} else if from == 2 && to == 3 {
				return 3
			}
			return 0
		},
	)

	for i, test := range []struct {
		l1   int
		l2   int
		cost float64
	}{
		{l1: 0, l2: 1, cost: 5},
		{l1: 1, l2: 2, cost: 10},
		{l1: 2, l2: 3, cost: 30},
		{l1: 3, l2: 2, cost: 0},
		{l1: 0, l2: 2, cost: 0},
		{l1: 5, l2: 5, cost: 0},
	} {
		if got := m.Cost(test.l1, test.l2); got != test.cost {
			t.Errorf("test %v: got %v; want %v", i, got, test.cost)
		}
	}

	b, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		t.Errorf("got %+v; want nil", err)
	}
	w := `{
	"measures": [
		{
			"count": 3,
			"measure": {
				"constant": 0,
				"type": "constant"
			}
		},
		{
			"count": 1,
			"measure": {
				"type": "euclidean"
			}
		},
		{
			"count": 1,
			"measure": {
				"constant": 10,
				"type": "constant"
			}
		},
		{
			"count": 1,
			"measure": {
				"constant": 30,
				"type": "constant"
			}
		}
	],
	"type": "bin"
}`
	if v := string(b); v != w {
		t.Errorf("got %s\nwant %s\n", v, w)
	}
}
