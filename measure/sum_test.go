package measure_test

import (
	"encoding/json"
	"testing"

	"github.com/nextmv-io/sdk/measure"
)

func TestSum(t *testing.T) {
	points := []measure.Point{{1}, {5}}
	m := measure.Sum(
		measure.Constant(10),
		measure.Scale(measure.Indexed(measure.EuclideanByPoint(), points), 3),
	)

	if v := m.Cost(0, 1); v != 22 {
		t.Errorf("got %v; want 22", v)
	}

	b, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		t.Errorf("got %+v; want nil", err)
	}
	w := `{
	"measures": [
		{
			"constant": 10,
			"type": "constant"
		},
		{
			"measure": {
				"type": "euclidean"
			},
			"scale": 3,
			"type": "scale"
		}
	],
	"type": "sum"
}`
	if v := string(b); v != w {
		t.Errorf("got %s; want %s", v, w)
	}
}
