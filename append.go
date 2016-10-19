package multierror

import (
	"reflect"
)

// removeNils preserves the non-nil elements
// of a slice. This is a destructive operation
// the contents of the original slice are modified.
func removeNils(errs []error) []error {
	view := errs[:0]
	for _, err := range errs {
		if err == nil {
			continue
		}
		add := true
		switch reflect.TypeOf(err).Kind() {
		case reflect.Chan, reflect.Func,
			reflect.Interface, reflect.Map,
			reflect.Ptr, reflect.Slice:
			add = !reflect.ValueOf(err).IsNil()
		}
		if add {
			view = append(view, err)
		}
	}
	return view
}

// AppendNonNil is a helper function that will append more errors
// onto an Error in order to create a larger multi-error.
//
// If err is not a multierror.Error, then it will be turned into
// one. If any of the errs are multierr.Error, they will be flattened
// one level into err.
//
// nil values in errs are filtered out. If the err is nil and
// the length of filtered errs is zero then the function returns nil.
func AppendNonNil(err error, errs ...error) error {

	errs = removeNils(errs)
	// Preserve input value when no errors have occurred
	// Preserve output value when only one error is produced
	if len(errs) == 0 {
		return err
	} else if (err == nil) && len(errs) == 1 {
		return errs[0]
	}
	return Append(err, errs...)
}

// Append is a helper function that will append more errors
// onto an Error in order to create a larger multi-error.
//
// If err is not a multierror.Error, then it will be turned into
// one. If any of the errs are multierr.Error, they will be flattened
// one level into err.
func Append(err error, errs ...error) *Error {
	switch err := err.(type) {
	case *Error:
		// Typed nils can reach here, so initialize if we are nil
		if err == nil {
			err = new(Error)
		}

		// Go through each error and flatten
		for _, e := range errs {
			switch e := e.(type) {
			case *Error:
				err.Errors = append(err.Errors, e.Errors...)
			default:
				err.Errors = append(err.Errors, e)
			}
		}

		return err
	default:
		newErrs := make([]error, 0, len(errs)+1)
		if err != nil {
			newErrs = append(newErrs, err)
		}
		newErrs = append(newErrs, errs...)

		return Append(&Error{}, newErrs...)
	}
}
