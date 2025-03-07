// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package multierror

import (
	"fmt"
)

// Prefix is a helper function that will prefix some text
// to the given error. If the error is a multierror.Error, then
// it will be prefixed to each wrapped error.
//
// This is useful to use when appending multiple multierrors
// together in order to give better scoping.
func Prefix(err error, prefix string) error {
	if err == nil {
		return nil
	}

	switch err := err.(type) {
	case *Error:
		// Typed nils can reach here, so initialize if we are nil
		if err == nil {
			err = new(Error)
		}

		// Wrap each of the errors
		for i, e := range err.Errors {
			err.Errors[i] = fmt.Errorf("%s %s", prefix, e)
		}

		return err
	default:
		return fmt.Errorf("%s %s", prefix, err)
	}
}
