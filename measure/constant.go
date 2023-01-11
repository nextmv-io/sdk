package measure

import "encoding/json"

// ConstantByPoint measure always estimates the same cost.
func ConstantByPoint(c float64) ByPoint {
	return constantByPoint(c)
}

type constantByPoint float64

func (c constantByPoint) Cost(_, _ Point) float64 {
	return float64(c)
}

func (c constantByPoint) Triangular() bool {
	return true
}

func (c constantByPoint) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":     "constant",
		"constant": float64(c),
	})
}

// Constant measure always estimates the same cost.
func Constant(c float64) ByIndex {
	return constant(c)
}

type constant float64

func (c constant) Cost(_, _ int) float64 {
	return float64(c)
}

func (c constant) Triangular() bool {
	return true
}

func (c constant) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":     "constant",
		"constant": float64(c),
	})
}
