package multierror

import (
	"fmt"

	"github.com/hashicorp/errwrap"
)

// Suffix is a helper function that will suffix some text
// to the given error. If the error is a multierror.Error, then
// it will be suffixed to each wrapped error.
//
// This is useful when Prefix would make the root error
// hard to read by putting too much text in front of it.
func Suffix(err error, suffix string) error {
	if err == nil {
		return nil
	}

	format := fmt.Sprintf("{{err}} %s", suffix)
	switch err := err.(type) {
	case *Error:
		// Typed nils can reach here, so initialize if we are nil
		if err == nil {
			err = new(Error)
		}

		// Wrap each of the errors
		for i, e := range err.Errors {
			err.Errors[i] = errwrap.Wrapf(format, e)
		}

		return err
	default:
		return errwrap.Wrapf(format, err)
	}
}
