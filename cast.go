package multierror

// Returns any given error (including nil) as a *Error object for compatibility
// which expects to always handle *Error rather than the stdlib error interface.
func Cast(err error) *Error {
	switch err := err.(type) {
	case *Error:
		return err
	default:
		if err == nil {
			return &Error{}
		}
		return &Error{
			Errors: []error{err},
		}
	}
}
