package measure_test

import (
	"encoding/json"
	"testing"

	"github.com/nextmv-io/sdk/measure"
)

func TestTaxicabByPointCost(t *testing.T) {
	p1 := measure.Point{1, 2}
	p2 := measure.Point{4, 6}

	m := measure.TaxicabByPoint()

	if v := m.Cost(p1, p2); v != 7 {
		t.Errorf("got %v; want 7", v)
	}
	if v := m.Cost(p2, p1); v != 7 {
		t.Errorf("got %v; want 7", v)
	}
}

func TestTaxicabByPointMarshalJSON(t *testing.T) {
	b, err := json.Marshal(measure.TaxicabByPoint())
	if err != nil {
		t.Errorf("got %+v; want nil", err)
	}
	w := `{"type":"taxicab"}`
	if v := string(b); v != w {
		t.Errorf("got %q; want %q", v, w)
	}
}

func TestTaxicabByIndexCost(t *testing.T) {
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
				measure.TaxicabByPoint(), []measure.Point{{1}, {1}}),
			f: func(m measure.ByIndex) float64 {
				return m.Cost(0, 1)
			},
		},
		{
			expected: 1,
			m: measure.Indexed(
				measure.TaxicabByPoint(), []measure.Point{{1}, {2}}),
			f: func(m measure.ByIndex) float64 {
				return m.Cost(0, 1)
			},
		},
		// Negative values
		{
			expected: 10,
			m: measure.Indexed(
				measure.TaxicabByPoint(), []measure.Point{{-5, 2}, {2, -1}}),
			f: func(m measure.ByIndex) float64 {
				return m.Cost(1, 0)
			},
		},
		{
			expected: 4,
			m: measure.Indexed(
				measure.TaxicabByPoint(), []measure.Point{{-5, -2}, {-2, -1}}),
			f: func(m measure.ByIndex) float64 {
				return m.Cost(1, 0)
			},
		},
		// Mismatched dimensions
		{
			expected: 4,
			m: measure.Indexed(
				measure.TaxicabByPoint(), []measure.Point{{1}, {2, 3}}),
			f: func(m measure.ByIndex) float64 {
				return m.Cost(0, 1)
			},
		},
		{
			expected: 4,
			m: measure.Indexed(
				measure.TaxicabByPoint(), []measure.Point{{2, 3}, {1}}),
			f: func(m measure.ByIndex) float64 {
				return m.Cost(0, 1)
			},
		},
		{
			expected: 8,
			m: measure.Indexed(
				measure.TaxicabByPoint(), []measure.Point{{2, 3, 4}, {1}}),
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

func TestTaxicabByIndexMarshalJSON(t *testing.T) {
	b, err := json.Marshal(measure.Indexed(measure.TaxicabByPoint(), nil))
	if err != nil {
		t.Errorf("got %+v; want nil", err)
	}
	w := `{"type":"taxicab"}`
	if v := string(b); v != w {
		t.Errorf("got %q; want %q", v, w)
	}
}
