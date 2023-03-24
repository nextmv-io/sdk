package common

import (
	"math"
)

// Haversine calculates the distance between two locations using the
// Haversine formula. Haversine is a good approximation for short
// distances (up to a few hundred kilometers).
func Haversine(from, to Location) Distance {
	if !from.valid || !to.valid {
		return NewDistance(0, Meters)
	}

	x1 := degreesToRadian(from.Longitude())
	y1 := degreesToRadian(from.Latitude())
	x2 := degreesToRadian(to.Longitude())
	y2 := degreesToRadian(to.Latitude())

	dx := x1 - x2
	dy := y1 - y2

	sdy := math.Sin(dy / 2)
	sdx := math.Sin(dx / 2)
	a := (sdy * sdy) + math.Cos(y1)*math.Cos(y2)*sdx*sdx

	return NewDistance(
		2*radius*math.Atan2(math.Sqrt(a), math.Sqrt(1-a)),
		Meters,
	)
}

func degreesToRadian(d float64) float64 {
	return d * math.Pi / 180.0
}

const radius = 6371 * 1000
