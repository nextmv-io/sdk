package measure

import (
	"encoding/json"
)

// Scale the cost of some other measure by a constant.
func Scale(m ByIndex, constant float64) ByIndex {
	return scale{m: m, constant: constant}
}

type scale struct {
	m        ByIndex
	constant float64
}

func (s scale) Cost(from, to int) float64 {
	return s.m.Cost(from, to) * s.constant
}

func (s scale) Triangular() bool {
	return IsTriangular(s.m)
}

func (s scale) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"measure": s.m,
		"type":    "scale",
		"scale":   s.constant,
	})
}

// ScaleByPoint scales an underlying measure by a constant.
func ScaleByPoint(m ByPoint, constant float64) ByPoint {
	return scaleByPoint{m: m, constant: constant}
}

type scaleByPoint struct {
	m        ByPoint
	constant float64
}

func (s scaleByPoint) Cost(p1, p2 Point) float64 {
	return s.m.Cost(p1, p2) * s.constant
}

func (s scaleByPoint) Triangular() bool {
	return IsTriangular(s.m)
}

func (s scaleByPoint) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"measure": s.m,
		"type":    "scale",
		"scale":   s.constant,
	})
}
