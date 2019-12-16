package multierror

import (
	"fmt"
)

// Error is an error type to track multiple errors. This is used to
// accumulate errors in cases and return them as a single "error".
type Error struct {
	Errors      []error
	ErrorFormat ErrorFormatFunc
}

func (e *Error) Error() string {
	fn := e.ErrorFormat
	if fn == nil {
		fn = ListFormatFunc
	}

	return fn(e.Errors)
}

func (e *Error) GoString() string {
	return fmt.Sprintf("*%#v", *e)
}

// WrappedErrors returns the list of errors that this Error is wrapping.
// It is an implementation of the errwrap.Wrapper interface so that
// multierror.Error can be used with that library.
//
// This method is not safe to be called concurrently and is no different
// than accessing the Errors field directly. It is implemented only to
// satisfy the errwrap.Wrapper interface.
func (e *Error) WrappedErrors() []error {
	return e.Errors
}
