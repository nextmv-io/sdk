package measure

import (
	"encoding/json"
	"math"
)

// TaxicabByPoint adds absolute distances between two points in all dimensions.
func TaxicabByPoint() ByPoint {
	return taxicabByPoint{}
}

type taxicabByPoint struct{}

func (t taxicabByPoint) Cost(p1, p2 Point) float64 {
	len1, len2 := len(p1), len(p2)
	longest := math.Max(float64(len1), float64(len2))

	var sum float64
	for i := 0; float64(i) < longest; i++ {
		var n1, n2 float64
		if i >= len1 {
			n1 = 0
		} else {
			n1 = p1[i]
		}
		if i >= len2 {
			n2 = 0
		} else {
			n2 = p2[i]
		}

		sum += math.Abs(n1 - n2)
	}

	return sum
}

func (t taxicabByPoint) Triangular() bool {
	return true
}

func (t taxicabByPoint) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"type": "taxicab"})
}
