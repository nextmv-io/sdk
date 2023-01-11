package measure

import (
	"encoding/json"
	"math"
)

// HaversineByPoint estimates meters connecting two points along the surface
// of the earth.
func HaversineByPoint() ByPoint {
	return haversineByPoint{}
}

type haversineByPoint struct{}

func (h haversineByPoint) Cost(p1, p2 Point) float64 {
	x1, y1 := toRadians(p1)
	x2, y2 := toRadians(p2)

	sdy := math.Sin((y1 - y2) / 2.0)
	sdx := math.Sin((x1 - x2) / 2.0)

	// to help the compiler generate efficient code we square sdy and sdx by
	// using multiplication instead of a call to math.Pow(x, 2.0)
	a := sdy*sdy + math.Cos(y1)*math.Cos(y2)*sdx*sdx

	x := 2.0 * radius * math.Atan2(math.Sqrt(a), math.Sqrt(1.0-a))

	// on different architectures x is slightly different. we don't need perfect
	// precision here, but we do care about reproducability, so we are fine with
	// a value that is precise up to millimeters
	return math.Floor(x*1000.0) / 1000.0
}

func (h haversineByPoint) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"type": "haversine"})
}

const radius = 6371.0 * 1000.0 // radius of the earth in meters

func degToRad(d float64) float64 {
	return d * math.Pi / 180.0
}

// toRadians converts points as a float slice of degrees to separate
// discrete values as radians.
//
// ([]{1, 2}, []{3, 4}) => (0.017..., 0.034...)
func toRadians(p []float64) (x, y float64) {
	switch len(p) {
	case 0:
		x, y = 0, 0
	case 1:
		x, y = degToRad(p[0]), 0
	default:
		x, y = degToRad(p[0]), degToRad(p[1])
	}

	return x, y
}

func (h haversineByPoint) Triangular() bool {
	return true
}
