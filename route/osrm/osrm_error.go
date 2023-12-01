package osrm

// NewError returns a new NewError.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/osrm] instead.
func NewError(err string, inputErr bool) Error {
	return Error{err, inputErr}
}

// Error is an error that reflects an error that is the user's fault.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/osrm] instead.
type Error struct {
	err      string
	inputErr bool
}

func (e Error) Error() string {
	return e.err
}

// IsInputError returns true if the error is the user's fault.
func (e Error) IsInputError() bool {
	return e.inputErr
}
