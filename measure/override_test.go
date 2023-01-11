package measure_test

import (
	"encoding/json"
	"testing"

	"github.com/nextmv-io/sdk/measure"
)

func TestOverride(t *testing.T) { //nolint:dupl
	points := []measure.Point{{1, 2}, {4, 6}, {6, 6}}
	m := measure.Override(
		measure.Indexed(measure.EuclideanByPoint(), points),
		measure.Constant(100),
		func(from, to int) bool {
			return from == 1 || to == 2
		},
	)

	for i, test := range []struct {
		l1   int
		l2   int
		cost float64
	}{
		{l1: 0, l2: 1, cost: 5},
		{l1: 2, l2: 1, cost: 2},
		{l1: 0, l2: 2, cost: 100},
		{l1: 1, l2: 0, cost: 100},
		{l1: 1, l2: 1, cost: 100},
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
	"default": {
		"type": "euclidean"
	},
	"override": {
		"constant": 100,
		"type": "constant"
	},
	"type": "override"
}`
	if v := string(b); v != w {
		t.Errorf("got %s\nwant %s\n", v, w)
	}
}

func TestDebugOverride(t *testing.T) { //nolint:dupl
	points := []measure.Point{{1, 2}, {4, 6}, {6, 6}}
	m := measure.DebugOverride(
		measure.Indexed(measure.EuclideanByPoint(), points),
		measure.Constant(100),
		func(from, to int) bool {
			return from == 1 || to == 2
		},
	)

	for i, test := range []struct {
		l1   int
		l2   int
		cost float64
	}{
		{l1: 0, l2: 1, cost: 5},
		{l1: 2, l2: 1, cost: 2},
		{l1: 0, l2: 2, cost: 100},
		{l1: 1, l2: 0, cost: 100},
		{l1: 1, l2: 1, cost: 100},
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
	"counts": {
		"default": 2,
		"override": 3
	},
	"default": {
		"type": "euclidean"
	},
	"override": {
		"constant": 100,
		"type": "constant"
	},
	"type": "override"
}`
	if v := string(b); v != w {
		t.Errorf("got %s\nwant %s\n", v, w)
	}
}
