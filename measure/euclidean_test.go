package measure_test

import (
	"encoding/json"
	"testing"

	"github.com/nextmv-io/sdk/measure"
)

func TestEuclideanByPointCost(t *testing.T) {
	p1 := measure.Point{1, 2}
	p2 := measure.Point{4, 6}

	m := measure.EuclideanByPoint()

	if v := m.Cost(p1, p2); v != 5 {
		t.Errorf("got %v; want 5", v)
	}
	if v := m.Cost(p2, p1); v != 5 {
		t.Errorf("got %v; want 5", v)
	}
}

func TestEuclideanByPointMarshalJSON(t *testing.T) {
	b, err := json.Marshal(measure.EuclideanByPoint())
	if err != nil {
		t.Errorf("got %+v; want nil", err)
	}
	w := `{"type":"euclidean"}`
	if v := string(b); v != w {
		t.Errorf("got %q; want %q", v, w)
	}
}

func TestEuclideanByIndexCost(t *testing.T) {
	// All bad configs return zero and don't panic
	for i, tc := range []struct {
		m           measure.ByIndex
		f           func(m measure.ByIndex) float64
		expected    float64
		shouldPanic bool
	}{
		// No point slice
		{
			m: measure.Indexed(measure.EuclideanByPoint(), nil),
			f: func(m measure.ByIndex) float64 {
				return m.Cost(1, 2)
			},
			shouldPanic: true,
		},
		// No points in slice
		{
			m: measure.Indexed(measure.EuclideanByPoint(), []measure.Point{}),
			f: func(m measure.ByIndex) float64 {
				return m.Cost(1, 2)
			},
			shouldPanic: true,
		},
		// Missing Y coords in points uses 0
		{
			expected: 0,
			m: measure.Indexed(
				measure.EuclideanByPoint(), []measure.Point{{1}, {1}}),
			f: func(m measure.ByIndex) float64 {
				return m.Cost(0, 1)
			},
		},
		{
			expected: 1,
			m: measure.Indexed(
				measure.EuclideanByPoint(), []measure.Point{{1}, {2}}),
			f: func(m measure.ByIndex) float64 {
				return m.Cost(0, 1)
			},
		},
		// Only 1 point and out of bounds panics
		{
			m: measure.Indexed(
				measure.EuclideanByPoint(), []measure.Point{{1, 1}}),
			f: func(m measure.ByIndex) float64 {
				return m.Cost(0, 1)
			},
			shouldPanic: true,
		},
		// Negative indices
		{
			m: measure.Indexed(
				measure.EuclideanByPoint(), []measure.Point{{1, 2}, {2, 3}}),
			f: func(m measure.ByIndex) float64 {
				return m.Cost(-1, -12)
			},
			shouldPanic: true,
		},
		// Out of bounds
		{
			m: measure.Indexed(
				measure.EuclideanByPoint(), []measure.Point{{1, 2}, {2, 3}}),
			f: func(m measure.ByIndex) float64 {
				return m.Cost(2, 4)
			},
			shouldPanic: true,
		},
		// Mismatched dimensions
		{
			expected: 3.1622776601683795,
			m: measure.Indexed(
				measure.EuclideanByPoint(), []measure.Point{{1}, {2, 3}}),
			f: func(m measure.ByIndex) float64 {
				return m.Cost(0, 1)
			},
		},
		{
			expected: 3.1622776601683795,
			m: measure.Indexed(
				measure.EuclideanByPoint(), []measure.Point{{2, 3}, {1}}),
			f: func(m measure.ByIndex) float64 {
				return m.Cost(0, 1)
			},
		},
		{
			expected: 5.0990195135927845,
			m: measure.Indexed(
				measure.EuclideanByPoint(), []measure.Point{{2, 3, 4}, {1}}),
			f: func(m measure.ByIndex) float64 {
				return m.Cost(0, 1)
			},
		},
	} {
		func() {
			defer func() {
				if r := recover(); r != nil {
					if !tc.shouldPanic {
						t.Errorf("test %d unexpectedly panicked", i)
					}
				}
			}()
			if v := tc.f(tc.m); v != tc.expected {
				t.Errorf("test %d: want: %v got: %v", i, tc.expected, v)
			}
		}()
	}
}

func TestEuclideanByIndexMarshalJSON(t *testing.T) {
	b, err := json.Marshal(measure.Indexed(measure.EuclideanByPoint(), nil))
	if err != nil {
		t.Errorf("got %+v; want nil", err)
	}
	w := `{"type":"euclidean"}`
	if v := string(b); v != w {
		t.Errorf("got %q; want %q", v, w)
	}
}
