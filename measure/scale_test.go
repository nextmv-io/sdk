package measure_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/nextmv-io/sdk/measure"
)

func TestScale(t *testing.T) {
	m := measure.Scale(measure.Constant(50), 0.1)

	if v := m.Cost(100, 42); v != 5 {
		t.Errorf("got %v; want 5", v)
	}

	b, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		t.Errorf("got %+v; want nil", err)
	}
	w := `{
	"measure": {
		"constant": 50,
		"type": "constant"
	},
	"scale": 0.1,
	"type": "scale"
}`
	if v := string(b); v != w {
		t.Errorf("got %s; want %s", v, w)
	}
}

func TestScaleByPoint(t *testing.T) {
	for i, tc := range []struct {
		inner    measure.ByPoint // measure to scale
		from     measure.Point   // distance from this point
		to       measure.Point   // to this point
		scale    float64         // scale constant to apply
		expected float64         // expected distance
	}{
		// Taxicab
		{
			inner:    measure.TaxicabByPoint(),
			from:     measure.Point{0, 0},
			to:       measure.Point{100, 100},
			scale:    0.5,
			expected: 100,
		},
		// Euclidean
		{
			inner:    measure.EuclideanByPoint(),
			from:     measure.Point{0, 0},
			to:       measure.Point{100, 100},
			scale:    1.5,
			expected: 212.13203435596427,
		},
		// Haversine
		{
			inner:    measure.HaversineByPoint(),
			from:     measure.Point{8.7536, 51.7173},   // Paderborn
			to:       measure.Point{-75.1647, 39.9525}, // Philadelphia
			scale:    0.1,
			expected: 626380.7847000001,
		},
		// Double scale (cancel each other)
		{
			inner:    measure.ScaleByPoint(measure.TaxicabByPoint(), 2),
			from:     measure.Point{0, 0},
			to:       measure.Point{100, 100},
			scale:    0.5,
			expected: 200,
		},
	} {
		m := measure.ScaleByPoint(tc.inner, tc.scale)
		if v := m.Cost(tc.from, tc.to); v != tc.expected {
			t.Errorf("test %d: got %v; want %v", i, v, tc.expected)
		}

		b, err := json.Marshal(m)
		if err != nil {
			t.Errorf("test %d: got %+v; want nil", i, err)
		}
		innerB, err := json.Marshal(tc.inner)
		if err != nil {
			t.Errorf("test %d: got %+v; want nil", i, err)
		}
		w := fmt.Sprintf(`{"measure":%s,"scale":%v,"type":"scale"}`,
			string(innerB), tc.scale)
		if v := string(b); v != w {
			t.Errorf("test %d: got %s; want %s", i, v, w)
		}
	}
}
