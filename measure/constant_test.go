package measure_test

import (
	"encoding/json"
	"testing"

	"github.com/nextmv-io/sdk/measure"
)

func TestConstant(t *testing.T) {
	m := measure.Constant(50)

	if v := m.Cost(100, 42); v != 50 {
		t.Errorf("got %v; want 50", v)
	}

	b, err := json.Marshal(m)
	if err != nil {
		t.Errorf("got %+v; want nil", err)
	}
	w := `{"constant":50,"type":"constant"}`
	if v := string(b); v != w {
		t.Errorf("got %q; want %q", v, w)
	}
}

func TestConstantByPointCost(t *testing.T) {
	p1 := measure.Point{-77.035382, 38.898269}
	p2 := measure.Point{-77.011250, 38.889789}

	m := measure.ConstantByPoint(10)

	if v := m.Cost(p1, p2); v != 10 {
		t.Errorf("got %v; want 10", v)
	}

	b, err := json.Marshal(m)
	if err != nil {
		t.Errorf("got %+v; want nil", err)
	}
	w := `{"constant":10,"type":"constant"}`
	if v := string(b); v != w {
		t.Errorf("got %q; want %q", v, w)
	}
}
