package osrm

import o "github.com/nextmv-io/sdk/measure/osrm"

// NewError returns a new NewError.
func NewError(err string, inputErr bool) Error {
	return o.NewError(err, inputErr)
}

// Error is an error that reflects an error that is the user's fault.
type Error = o.Error
