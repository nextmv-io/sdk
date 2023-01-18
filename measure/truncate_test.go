package measure_test

import (
	"encoding/json"
	"testing"

	"github.com/nextmv-io/sdk/measure"
)

func TestTruncate(t *testing.T) {
	// Inside the bounds.
	m := measure.Truncate(measure.Constant(25), 10, 45)
	if v := m.Cost(100, 45); v != 25 {
		t.Errorf("got %v; want 25", v)
	}

	// Over the upper bound.
	m = measure.Truncate(measure.Constant(50), 10, 45)
	if v := m.Cost(100, 45); v != 45 {
		t.Errorf("got %v; want 45", v)
	}

	// Below the lower bound.
	m = measure.Truncate(measure.Constant(10), 15, 20)
	if v := m.Cost(100, 15); v != 15 {
		t.Errorf("got %v; want 15", v)
	}

	b, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		t.Errorf("got %+v; want nil", err)
	}
	w := `{
	"lower": 15,
	"measure": {
		"constant": 10,
		"type": "constant"
	},
	"type": "truncate",
	"upper": 20
}`
	if v := string(b); v != w {
		t.Errorf("got %s; want %s", v, w)
	}
}
