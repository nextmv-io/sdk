package google_test

import (
	"fmt"

	"github.com/nextmv-io/sdk/route/google"
	"googlemaps.github.io/maps"
)

// Minimal example of how to create a client and matrix request, assuming the
// points are in the form longitude, latitude.
// Please note that this example does not define an output as it requires API
// authentication and is here for illustrative purposes only.
func Example() {
	points := [][2]float64{
		{-74.028297, 4.875835},
		{-74.046965, 4.872842},
		{-74.041763, 4.885648},
	}
	coords := make([]string, len(points))
	for p, point := range points {
		coords[p] = fmt.Sprintf("%f,%f", point[1], point[0])
	}
	r := &maps.DistanceMatrixRequest{
		Origins:      coords,
		Destinations: coords,
	}
	c, err := maps.NewClient(maps.WithAPIKey("<YOUR-API-KEY>"))
	if err != nil {
		panic(err)
	}

	// Distance and duration matrices can be constructed with the functions
	// provided in the package.
	dist, dur, err := google.DistanceDurationMatrices(c, r)
	if err != nil {
		panic(err)
	}

	// Once the measures have been created, you may estimate the distances and
	// durations by calling the Cost function.
	for p1 := range points {
		for p2 := range points {
			fmt.Printf(
				"(%d, %d) = [%f, %f]\n",
				p1, p2, dist.Cost(p1, p2), dur.Cost(p1, p2),
			)
		}
	}
	// This is the expected output.
	// (0, 0) = [0.000000, 0.000000]
	// (0, 1) = [6526.000000, 899.000000]
	// (0, 2) = [4889.000000, 669.000000]
	// (1, 0) = [5211.000000, 861.000000]
	// (1, 1) = [0.000000, 0.000000]
	// (1, 2) = [2260.000000, 302.000000]
	// (2, 0) = [3799.000000, 638.000000]
	// (2, 1) = [2260.000000, 311.000000]
	// (2, 2) = [0.000000, 0.000000]
}

// Making a request to retrieve polylines works similar. In this example we
// make a request to retrieve polylines by creating a DirectionsRequest. The
// polylines function returns a polyline from start to end and a slice of
// polylines for each leg, given through the waypoints. All polylines are
// encoded in Google's polyline format.
// Please note that this example does not define an output as it requires API
// authentication and is here for illustrative purposes only.
func Example_polylines() {
	points := [][2]float64{
		{-74.028297, 4.875835},
		{-74.046965, 4.872842},
		{-74.041763, 4.885648},
	}
	coords := make([]string, len(points))
	rPoly := &maps.DirectionsRequest{
		Origin:      coords[0],
		Destination: coords[len(coords)-1],
		Waypoints:   coords[1 : len(coords)-1],
	}
	c, err := maps.NewClient(maps.WithAPIKey("<YOUR-API-KEY>"))
	if err != nil {
		panic(err)
	}
	fullPoly, polyLegs, err := google.Polylines(c, rPoly)
	if err != nil {
		panic(err)
	}
	_, _ = fullPoly, polyLegs
}
