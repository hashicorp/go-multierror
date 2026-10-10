// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package multierror

// Flatten flattens the given error, merging any *Errors together into
// a single *Error.
func Flatten(err error) error {
	// If it isn't an *Error, just return the error as-is.
	// A typed nil is an *Error and would panic below. Len treats that as empty.
	me, ok := err.(*Error)
	if !ok {
		return err
	}
	if me == nil {
		return nil
	}

	// Otherwise, make the result and flatten away!
	flatErr := new(Error)
	flatten(me, flatErr)
	return flatErr
}

func flatten(err error, flatErr *Error) {
	switch err := err.(type) {
	case *Error:
		if err == nil {
			return
		}
		for _, e := range err.Errors {
			flatten(e, flatErr)
		}
	default:
		flatErr.Errors = append(flatErr.Errors, err)
	}
}
