package here

import (
	h "github.com/nextmv-io/sdk/measure/here"
)

// NewClient returns a new OSRM Client.
func NewClient(apiKey string, opts ...h.ClientOption) h.Client {
	return h.NewClient(apiKey, opts...)
}
