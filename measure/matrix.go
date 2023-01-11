package measure

import "encoding/json"

// Matrix measure returns pre-computed cost between two locations. Cost is
// assumed to be asymmetric.
func Matrix(arcs [][]float64) ByIndex {
	return matrix(arcs)
}

type matrix [][]float64

func (m matrix) Cost(from, to int) float64 {
	return m[from][to]
}

func (m matrix) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{"type": "matrix", "matrix": [][]float64(m)})
}
