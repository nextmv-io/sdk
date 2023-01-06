package measure_test

import (
	"encoding/json"
	"testing"

	"github.com/nextmv-io/sdk/measure"
)

func TestHaversineByPointCost(t *testing.T) {
	p1 := measure.Point{-77.035382, 38.898269}
	p2 := measure.Point{-77.011250, 38.889789}

	m := measure.HaversineByPoint()

	if v := int(m.Cost(p1, p2)); v != 2291 {
		t.Errorf("got %v; want 2291", v)
	}
	if v := int(m.Cost(p2, p1)); v != 2291 {
		t.Errorf("got %v; want 2291", v)
	}
	if v := int(m.Cost(nil, nil)); v != 0 {
		t.Errorf("got %v; want 0", v)
	}
}

func TestHaversineByPointARMVsAMDCost(t *testing.T) {
	p1 := measure.Point{7.6048, 51.9492}
	p2 := measure.Point{7.5921, 51.9573}

	m := measure.HaversineByPoint()

	// on an amd, the untruncated value would be 1252.4761188966866
	// on an arm, the untruncated value would be 1252.4761188966863
	if v := m.Cost(p1, p2); v != 1252.476 {
		t.Errorf("got %v; want 1252.476", v)
	}
}

func TestHaversineByPointMarshalJSON(t *testing.T) {
	b, err := json.Marshal(measure.HaversineByPoint())
	if err != nil {
		t.Errorf("got %+v; want nil", err)
	}
	w := `{"type":"haversine"}`
	if v := string(b); v != w {
		t.Errorf("got %q; want %q", v, w)
	}
}

func TestHaversineByIndexCost(t *testing.T) {
	// All bad configs return zero and don't panic
	for i, tc := range []struct {
		m        measure.ByIndex
		f        func(m measure.ByIndex) float64
		expected float64
	}{
		// Missing Y coords in points uses 0
		{
			expected: 0,
			m: measure.Indexed(
				measure.HaversineByPoint(), []measure.Point{{1}, {1}}),
			f: func(m measure.ByIndex) float64 {
				return m.Cost(0, 1)
			},
		},
		{
			expected: 111194.926,
			m: measure.Indexed(
				measure.HaversineByPoint(), []measure.Point{{1}, {2}}),
			f: func(m measure.ByIndex) float64 {
				return m.Cost(0, 1)
			},
		},
	} {
		if v := tc.f(tc.m); v != tc.expected {
			t.Errorf("test %d: want: %v got: %v", i, tc.expected, v)
		}
	}
}

func TestHaversineByIndexMarshalJSON(t *testing.T) {
	b, err := json.Marshal(measure.Indexed(measure.HaversineByPoint(), nil))
	if err != nil {
		t.Errorf("got %+v; want nil", err)
	}
	w := `{"type":"haversine"}`
	if v := string(b); v != w {
		t.Errorf("got %q; want %q", v, w)
	}
}

func BenchmarkHaversineByPoint(b *testing.B) {
	p1 := measure.Point{-77.035382, 38.898269}
	p2 := measure.Point{-77.011250, 38.889789}
	m := measure.HaversineByPoint()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Cost(p1, p2)
	}
}
