package multierror

import "testing"

func TestFlattenTypedNilError(t *testing.T) {
	var nilErr *Error
	if got := Flatten(nilErr); got != nil {
		t.Errorf("Flatten(typed nil *Error) = %v, want nil", got)
	}

	// A typed nil nested inside a multi-error should be dropped, not panic.
	e := &Error{Errors: []error{nilErr, errTest("real")}}
	got := Flatten(e).(*Error)
	if len(got.Errors) != 1 || got.Errors[0].Error() != "real" {
		t.Errorf("Flatten with nested typed nil = %v, want just [real]", got.Errors)
	}
}

type errTest string

func (e errTest) Error() string { return string(e) }
