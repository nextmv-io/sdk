package measure

import (
	"encoding/json"
)

// Sum adds other measures together.
func Sum(m ...ByIndex) ByIndex {
	return sum(m)
}

type sum []ByIndex

func (s sum) Cost(from, to int) float64 {
	c := 0.0
	for _, m := range s {
		c += m.Cost(from, to)
	}
	return c
}

func (s sum) Triangular() bool {
	for _, m := range s {
		if !IsTriangular(m) {
			return false
		}
	}
	return true
}

func (s sum) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"measures": []ByIndex(s),
		"type":     "sum",
	})
}
