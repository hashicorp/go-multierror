// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package multierror

// Flatten flattens the given error, merging any *Errors together into
// a single *Error.
func Flatten(err error) error {
	// A typed nil *Error satisfies the interface but carries no error.
	if e, ok := err.(*Error); ok {
		if e == nil {
			return nil
		}
	} else {
		return err
	}

	// Otherwise, make the result and flatten away!
	flatErr := new(Error)
	flatten(err, flatErr)
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
