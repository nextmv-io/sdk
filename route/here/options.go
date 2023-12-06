package here

import (
	"net/http"
	"time"
)

// ClientOption can pass options to be used with a HERE client.
//
// Deprecated: This package is deprecated and will be removed in the next major release.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
type ClientOption func(*client)

// MatrixOption is passed to functions on the Client that create matrices,
// configuring the HERE request the client will make.
//
// Deprecated: This package is deprecated and will be removed in the next major release.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
type MatrixOption func(req *matrixRequest)

// WithDepartureTime sets departure time to be used in the request. This will
// take traffic data into account for the given time. If no departure time is
// given, "any" will be used in the request and no traffic data is included,
// see official documentation for HERE matrix routing, concepts traffic.
//
// Deprecated: This package is deprecated and will be removed in the next major release.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
func WithDepartureTime(t time.Time) MatrixOption {
	return func(req *matrixRequest) {
		depTime := "any"
		if !t.IsZero() {
			depTime = t.Format(time.RFC3339)
		}
		req.DepartureTime = depTime
	}
}

// WithTransportMode sets the transport mode for the request.
//
// Deprecated: This package is deprecated and will be removed in the next major release.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
func WithTransportMode(mode TransportMode) MatrixOption {
	return func(req *matrixRequest) {
		req.TransportMode = mode
	}
}

// WithAvoidFeatures sets features that will be avoided in the calculated
// routes.
//
// Deprecated: This package is deprecated and will be removed in the next major release.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
func WithAvoidFeatures(features []Feature) MatrixOption {
	return func(req *matrixRequest) {
		featureStrs := make([]string, len(features))
		for i, f := range features {
			featureStrs[i] = string(f)
		}

		if req.Avoid == nil {
			req.Avoid = &avoid{
				Features: features,
			}
		} else {
			req.Avoid.Features = features
		}
	}
}

// WithAvoidAreas sets bounding boxes that will be avoided in the calculated
// routes.
//
// Deprecated: This package is deprecated and will be removed in the next major release.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
func WithAvoidAreas(areas []BoundingBox) MatrixOption {
	return func(req *matrixRequest) {
		as := make([]area, len(areas))
		for i, a := range areas {
			as[i] = area{
				Type:  "boundingBox",
				West:  a.West,
				South: a.South,
				East:  a.East,
				North: a.North,
			}
		}
		if req.Avoid == nil {
			req.Avoid = &avoid{
				Areas: as,
			}
		} else {
			req.Avoid.Areas = append(req.Avoid.Areas, as...)
		}
	}
}

// WithTruckProfile sets a Truck profile on the matrix request. The following
// attributes are required by HERE:
// * TunnelCategory: if this is an empty string, the Client will automatically
// set it to TunnelCategoryNone
// * Type
// * AxleCount.
//
// Deprecated: This package is deprecated and will be removed in the next major release.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
func WithTruckProfile(t Truck) MatrixOption {
	return func(req *matrixRequest) {
		if t.TunnelCategory == "" {
			t.TunnelCategory = TunnelCategoryNone
		}
		req.Truck = &t
	}
}

// WithScooterProfile sets a Scooter profile on the request.
//
// Deprecated: This package is deprecated and will be removed in the next major release.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
func WithScooterProfile(scooter Scooter) MatrixOption {
	return func(req *matrixRequest) {
		req.Scooter = &scooter
	}
}

// WithTaxiProfile sets a Taxi profile on the request.
//
// Deprecated: This package is deprecated and will be removed in the next major release.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
func WithTaxiProfile(taxi Taxi) MatrixOption {
	return func(req *matrixRequest) {
		req.Taxi = &taxi
	}
}

// WithClientTransport overwrites the RoundTripper used by the internal
// http.Client.
//
// Deprecated: This package is deprecated and will be removed in the next major release.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
func WithClientTransport(rt http.RoundTripper) ClientOption {
	if rt == nil {
		rt = http.DefaultTransport
	}

	return func(c *client) {
		c.httpClient.Transport = rt
	}
}

// WithDenyRedirectPolicy block redirected requests to specified hostnames.
// Matches hostname greedily e.g. google.com will match api.google.com,
// file.api.google.com, ...
//
// Deprecated: This package is deprecated and will be removed in the next major release.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
func WithDenyRedirectPolicy(hostnames ...string) ClientOption {
	return func(c *client) {
		c.denyRedirectedRequests = hostnames
	}
}
