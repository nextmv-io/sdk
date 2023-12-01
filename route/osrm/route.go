package osrm

import "github.com/nextmv-io/sdk/route"

// Polyline requests polylines for the given points. The first parameter returns
// a polyline from start to end and the second parameter returns a list of
// polylines, one per leg.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/osrm].
func Polyline(
	c Client, points []route.Point,
) (string, []string, error) {
	polyline, legLines, err := c.Polyline(points)
	if err != nil {
		return "", []string{}, err
	}

	return polyline, legLines, nil
}
