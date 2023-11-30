package here

import (
	"net/http"
	"time"

	h "github.com/nextmv-io/sdk/measure/here"
)

// ClientOption can pass options to be used with a HERE client.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/here] instead.
type ClientOption = h.ClientOption

// MatrixOption is passed to functions on the Client that create matrices,
// configuring the HERE request the client will make.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/here] instead.
type MatrixOption = h.MatrixOption

// WithDepartureTime sets departure time to be used in the request. This will
// take traffic data into account for the given time. If no departure time is
// given, "any" will be used in the request and no traffic data is included,
// see official documentation for HERE matrix routing, concepts traffic.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/here] instead.
func WithDepartureTime(t time.Time) MatrixOption {
	return h.WithDepartureTime(t)
}

// WithTransportMode sets the transport mode for the request.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/here] instead.
func WithTransportMode(mode TransportMode) MatrixOption {
	return h.WithTransportMode(mode)
}

// WithAvoidFeatures sets features that will be avoided in the calculated
// routes.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/here] instead.
func WithAvoidFeatures(features []Feature) MatrixOption {
	return h.WithAvoidFeatures(features)
}

// WithAvoidAreas sets bounding boxes that will be avoided in the calculated
// routes.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/here] instead.
func WithAvoidAreas(areas []BoundingBox) MatrixOption {
	return h.WithAvoidAreas(areas)
}

// WithTruckProfile sets a Truck profile on the matrix request. The following
// attributes are required by HERE:
// * TunnelCategory: if this is an empty string, the Client will automatically
// set it to TunnelCategoryNone
// * Type
// * AxleCount.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/here] instead.
func WithTruckProfile(t Truck) MatrixOption {
	return h.WithTruckProfile(t)
}

// WithScooterProfile sets a Scooter profile on the request.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/here] instead.
func WithScooterProfile(scooter Scooter) MatrixOption {
	return h.WithScooterProfile(scooter)
}

// WithTaxiProfile sets a Taxi profile on the request.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/here] instead.
func WithTaxiProfile(taxi Taxi) MatrixOption {
	return h.WithTaxiProfile(taxi)
}

// WithClientTransport overwrites the RoundTripper used by the internal
// http.Client.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/here] instead.
func WithClientTransport(rt http.RoundTripper) ClientOption {
	return h.WithClientTransport(rt)
}

// WithDenyRedirectPolicy block redirected requests to specified hostnames.
// Matches hostname greedily e.g. google.com will match api.google.com,
// file.api.google.com, ...
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/here] instead.
func WithDenyRedirectPolicy(hostnames ...string) ClientOption {
	return h.WithDenyRedirectPolicy(hostnames...)
}
