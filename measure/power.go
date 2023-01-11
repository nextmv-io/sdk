package measure

import (
	"encoding/json"
	"math"
)

// Power raises the cost of some other measure to an exponent.
func Power(m ByIndex, exponent float64) ByIndex {
	return power{m: m, exponent: exponent}
}

type power struct {
	m        ByIndex
	exponent float64
}

func (p power) Cost(from, to int) float64 {
	return math.Pow(p.m.Cost(from, to), p.exponent)
}

func (p power) Triangular() bool {
	return IsTriangular(p.m)
}

func (p power) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"measure":  p.m,
		"type":     "power",
		"exponent": p.exponent,
	})
}
