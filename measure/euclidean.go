package measure

import (
	"encoding/json"
	"math"
)

// EuclideanByPoint computes straight line distance connecting two indices.
func EuclideanByPoint() ByPoint {
	return euclideanByPoint{}
}

type euclideanByPoint struct{}

func (e euclideanByPoint) Cost(p1, p2 Point) float64 {
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

		sum += (n1 - n2) * (n1 - n2)
	}

	return math.Sqrt(sum)
}

func (e euclideanByPoint) Triangular() bool {
	return true
}

func (e euclideanByPoint) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"type": "euclidean"})
}
