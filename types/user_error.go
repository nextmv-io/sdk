package types

// NewUserError returns a new UserError.
func NewUserError(err string) UserError {
	return UserError{err}
}

// UserError is an error that reflects an error that is the user's fault.
type UserError struct {
	err string
}

func (e UserError) Error() string {
	return e.err
}
