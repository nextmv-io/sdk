package measure

import (
	"encoding/json"
	"fmt"
)

// Indexed creates a ByIndex measure from the given ByPoint measure
// and wrapping the provided points.
func Indexed(m ByPoint, points []Point) ByIndex {
	return byIndex{
		points:  points,
		ByPoint: m,
	}
}

type byIndex struct {
	points []Point
	ByPoint
}

func (m byIndex) Cost(from, to int) float64 {
	if from < 0 {
		panic(
			fmt.Sprintf("byIndex.Cost: invalid `from` index %d - must be >= 0", from),
		)
	}
	if to < 0 {
		panic(fmt.Sprintf("byIndex.Cost: invalid `to` index %d - must be >= 0", to))
	}
	if from >= len(m.points) {
		panic(fmt.Sprintf(
			"byIndex.Cost: invalid `from` index %d - must be < number of points (%d)",
			from, len(m.points),
		))
	}
	if to >= len(m.points) {
		panic(fmt.Sprintf(
			"byIndex.Cost: invalid `to` index %d - must be < number of points (%d)",
			to, len(m.points)),
		)
	}
	return m.ByPoint.Cost(m.points[from], m.points[to])
}

func (m byIndex) Triangular() bool {
	return IsTriangular(m.ByPoint)
}

func (m byIndex) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.ByPoint)
}
