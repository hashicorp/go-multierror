package multierror

func Flatten(err error) error {
	flatErr := new(Error)
	flatten(err, flatErr)
	return error(flatErr)
}

func flatten(err error, flatErr *Error) {
	switch err := err.(type) {
	case *Error:
		for _, e := range err.Errors {
			flatten(e, flatErr)
		}
	default:
		flatErr.Errors = append(flatErr.Errors, err)
	}
}
