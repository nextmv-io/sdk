package osrm

// NewError returns a new NewError.
//
// Deprecated: This package is deprecated and will be removed in the future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/osrm].
func NewError(err string, inputErr bool) Error {
	return Error{err, inputErr}
}

// Error is an error that reflects an error that is the user's fault.
//
// Deprecated: This package is deprecated and will be removed in the future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/osrm].
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
