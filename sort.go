package multierror

// Len implements sort.Interface function for length
func (err *Error) Len() int {
	if err == nil {
		return 0
	}
	return len(err.Errors)
}

// Swap implements sort.Interface function for swapping elements
func (err *Error) Swap(i, j int) {
	if err == nil {
		return
	}
	err.Errors[i], err.Errors[j] = err.Errors[j], err.Errors[i]
}

// Less implements sort.Interface function for determining order
func (err *Error) Less(i, j int) bool {
	if err == nil {
		return false
	}
	return err.Errors[i].Error() < err.Errors[j].Error()
}
