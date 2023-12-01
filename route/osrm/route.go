package osrm

import "github.com/nextmv-io/sdk/measure"

// Polyline requests polylines for the given points. The first parameter returns
// a polyline from start to end and the second parameter returns a list of
// polylines, one per leg.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/osrm] instead.
func Polyline(
	c Client, points []measure.Point,
) (string, []string, error) {
	polyline, legLines, err := c.Polyline(points)
	if err != nil {
		return "", []string{}, err
	}

	return polyline, legLines, nil
}
