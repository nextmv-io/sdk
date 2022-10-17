package osrm

import "github.com/nextmv-io/sdk/route"

// Polyline requests polylines for the given points. The first parameter returns
// a polyline from start to end and the second parameter returns a list of
// polylines, one per leg.
func Polyline(
	c Client, points []route.Point,
) (string, []string, error) {
	polyline, legLines, err := c.Polyline(points)
	if err != nil {
		return "", []string{}, err
	}

	return polyline, legLines, nil
}
