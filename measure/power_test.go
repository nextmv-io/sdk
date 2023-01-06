package measure_test

import (
	"encoding/json"
	"testing"

	"github.com/nextmv-io/sdk/measure"
)

func TestPower(t *testing.T) {
	m := measure.Power(measure.Constant(10), 2)

	if v := m.Cost(1, 5); v != 100 {
		t.Errorf("got %v; want 100", v)
	}

	b, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		t.Errorf("got %+v; want nil", err)
	}
	w := `{
	"exponent": 2,
	"measure": {
		"constant": 10,
		"type": "constant"
	},
	"type": "power"
}`
	if v := string(b); v != w {
		t.Errorf("got %s; want %s", v, w)
	}
}
