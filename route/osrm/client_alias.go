package osrm

import (
	"net/http"

	o "github.com/nextmv-io/sdk/measure/osrm"
)

// Endpoint defines the OSRM endpoint to be used.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/osrm] instead.
type Endpoint = o.Endpoint

const (
	// TableEndpoint is used to retrieve distance and duration matrices.
	TableEndpoint = o.TableEndpoint
	// RouteEndpoint is used to retrieve polylines for a set of points.
	RouteEndpoint = o.RouteEndpoint
)

// Client represents an OSRM client.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/osrm] instead.
type Client = o.Client

// NewClient returns a new OSRM Client.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/osrm] instead.
func NewClient(host string, opts ...ClientOption) Client {
	return o.NewClient(host, opts...)
}

// DefaultClient creates a new OSRM Client.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/osrm] instead.
func DefaultClient(host string, useCache bool) Client {
	return o.DefaultClient(host, useCache)
}

// TableOptions is a function that configures a tableConfig.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/osrm] instead.
type TableOptions = o.TableOptions

// WithDuration returns a TableOptions function for composing a tableConfig with
// duration data enabled, telling the OSRM server to include duration data in
// the response table data.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/osrm] instead.
func WithDuration() TableOptions {
	return o.WithDuration()
}

// WithDistance returns a TableOptions function for composing a tableConfig with
// distance data enabled, telling the OSRM server to include distance data in
// the response table data.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/osrm] instead.
func WithDistance() TableOptions {
	return o.WithDistance()
}

// ClientOption can pass options to be used with an OSRM client.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/osrm] instead.
type ClientOption = o.ClientOption

// WithClientTransport overwrites the RoundTripper used by the internal
// http.Client.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/osrm] instead.
func WithClientTransport(rt http.RoundTripper) ClientOption {
	return o.WithClientTransport(rt)
}

// WithCache configures the maximum number of results cached.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/osrm] instead.
func WithCache(maxItems int) ClientOption {
	return o.WithCache(maxItems)
}

// ParallelRuns set the number of parallel calls to the OSRM server. If 0 is
// passed, the default value of 16 will be used.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/osrm] instead.
func ParallelRuns(runs int) TableOptions {
	return o.ParallelRuns(runs)
}

// RouteResponse holds the route response from the OSRM server.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/osrm] instead.
type RouteResponse = o.RouteResponse

// Route partially represents the OSRM Route object.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/osrm] instead.
type Route = o.Route

// Leg partially represents the OSRM Leg object.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/osrm] instead.
type Leg = o.Leg

// Step partially represents the OSRM Step object.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/osrm] instead.
type Step = o.Step
