// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package multierror

// Len implements sort.Interface function for length
func (e *Error) Len() int {
	if e == nil {
		return 0
	}

	return len(e.Errors)
}

// Swap implements sort.Interface function for swapping elements
func (e *Error) Swap(i, j int) {
	e.Errors[i], e.Errors[j] = e.Errors[j], e.Errors[i]
}

// Less implements sort.Interface function for determining order
func (e *Error) Less(i, j int) bool {
	return e.Errors[i].Error() < e.Errors[j].Error()
}
