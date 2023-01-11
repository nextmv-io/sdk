package measure

import "sort"

// ByClockwise implements sort.Interface for sorting points clockwise around a
// central point.
func ByClockwise(center Point, points []Point) sort.Interface {
	return byClockwise{center: center, points: points}
}

type byClockwise struct {
	center Point
	points []Point
}

func (b byClockwise) Len() int {
	return len(b.points)
}

func (b byClockwise) Swap(i, j int) {
	b.points[i], b.points[j] = b.points[j], b.points[i]
}

func (b byClockwise) Less(i, j int) bool {
	return LessClockwise(b.center, b.points[i], b.points[j])
}

// LessClockwise returns true if a is closer to a central point than b, and
// false if it is not.
func LessClockwise(center, a, b Point) bool {
	aDX := a[0] - center[0]
	aDY := a[1] - center[1]
	bDX := b[0] - center[0]
	bDY := b[1] - center[1]

	if aDX >= 0 && bDX < 0 {
		return true
	}
	if aDX < 0 && bDX >= 0 {
		return false
	}
	if aDX == 0 && bDX == 0 {
		if aDY >= 0 || bDY >= 0 {
			return a[1] > b[1]
		}
		return b[1] > a[1]
	}

	// compute the cross product of vectors (depot -> a) * (depot -> b)
	det := (aDX * bDY) - (bDX * aDY)
	if det < 0 {
		return true
	}
	if det > 0 {
		return false
	}

	// points a and b are on the same line from the center
	// check which point is closer to the center
	d1 := (aDX * aDX) + (aDY * aDY)
	d2 := (bDX * bDX) + (bDY * bDY)
	return d1 < d2
}
