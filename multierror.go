package multierror

import (
	"fmt"
)

// Error is an error type to track multiple errors. This is used to
// accumulate errors in cases and return them as a single "error".
type Error struct {
	Errors []error
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

// ErrorAppend is a helper function that will append more errors
// onto an Error in order to create a larger multi-error.
//
// If err is not a multierror.Error, then it will be turned into
// one. If any of the errs are multierr.Error, they will be flattened
// one level into err.
func ErrorAppend(err error, errs ...error) *Error {
	if err == nil {
		err = new(Error)
	}

	switch err := err.(type) {
	case *Error:
		if err == nil {
			err = new(Error)
		}

		err.Errors = append(err.Errors, errs...)
		return err
	default:
		newErrs := make([]error, len(errs)+1)
		newErrs[0] = err
		copy(newErrs[1:], errs)
		return &Error{
			Errors: newErrs,
		}
	}
}
