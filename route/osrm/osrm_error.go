package osrm

import o "github.com/nextmv-io/sdk/measure/osrm"

// NewError returns a new NewError.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/osrm] instead.
func NewError(err string, inputErr bool) Error {
	return o.NewError(err, inputErr)
}

// Error is an error that reflects an error that is the user's fault.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/osrm] instead.
type Error = o.Error
