package multierror

// Call is a helper function that will call the provided callback and
// append the error to *err if not nil.
//
// This is useful to call from defers, to preserve errors while returned
// when deferring.
// e.g. defer multierror.Call(&err, fh.Close)
func Call(err *error, callback func() error) {
	*err = Append(*err, callback())
}
