// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package multierror

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestFlatten(t *testing.T) {
	original := &Error{
		Errors: []error{
			errors.New("one"),
			&Error{
				Errors: []error{
					errors.New("two"),
					&Error{
						Errors: []error{
							errors.New("three"),
						},
					},
				},
			},
		},
	}

	expected := `3 errors occurred:
	* one
	* two
	* three

`
	actual := fmt.Sprintf("%s", Flatten(original))

	if expected != actual {
		t.Fatalf("expected: %s, got: %s", expected, actual)
	}
}

func TestFlattenTypedNil(t *testing.T) {
	var result *Error
	if got := Flatten(result); got != nil {
		t.Fatalf("Flatten(typed nil) = %#v, want nil", got)
	}

	flat := Flatten(&Error{Errors: []error{errors.New("one"), (*Error)(nil)}}).(*Error)
	if len(flat.Errors) != 1 || flat.Errors[0].Error() != "one" {
		t.Fatalf("got %#v", flat.Errors)
	}
}

func TestFlatten_nonError(t *testing.T) {
	err := errors.New("foo")
	actual := Flatten(err)
	if !reflect.DeepEqual(actual, err) {
		t.Fatalf("bad: %#v", actual)
	}
}
