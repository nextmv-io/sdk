package here

import (
	h "github.com/nextmv-io/sdk/measure/here"
)

// NewClient returns a new OSRM Client.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use github.com/nextmv-io/sdk/measure/here instead.
func NewClient(apiKey string, opts ...h.ClientOption) h.Client {
	return h.NewClient(apiKey, opts...)
}
