package google

import (
	"github.com/nextmv-io/sdk/measure/google"
	"github.com/nextmv-io/sdk/route"
	"googlemaps.github.io/maps"
)

// DistanceDurationMatrices makes requests to the Google Distance Matrix API
// and returns route.ByIndex types to estimate distances (in meters) and
// durations (in seconds). It receives a Google Maps Client and Request. The
// coordinates passed to the request must be in the form latitude, longitude.
// The resulting distance and duration matrices are saved in memory. To find
// out how to create a client and request, please refer to the go package docs.
// This function takes into consideration the usage limits of the Distance
// Matrix API and thus may transform the request into multiple ones and handle
// them accordingly. You can find more about usage limits here in the official
// google maps documentation for the distance matrix, usage and billing.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/google] instead.
func DistanceDurationMatrices(c *maps.Client, r *maps.DistanceMatrixRequest) (
	route.ByIndex,
	route.ByIndex,
	error,
) {
	return google.DistanceDurationMatrices(c, r)
}

// Polylines requests polylines for the given points. The first parameter
// returns a polyline from start to end and the second parameter returns a list
// of polylines, one per leg.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/google] instead.
func Polylines(
	c *maps.Client,
	orgRequest *maps.DirectionsRequest,
) (string, []string, error) {
	return google.Polylines(c, orgRequest)
}
