package multierror

import (
	"errors"
	"testing"
)

func TestSuffix_Error(t *testing.T) {
	original := &Error{
		Errors: []error{errors.New("foo")},
	}

	result := Suffix(original, "bar")
	if result.(*Error).Errors[0].Error() != "foo bar" {
		t.Fatalf("bad: %s", result)
	}
}

func TestSuffix_NilError(t *testing.T) {
	var err error
	result := Suffix(err, "bar")
	if result != nil {
		t.Fatalf("bad: %#v", result)
	}
}

func TestSuffix_NonError(t *testing.T) {
	original := errors.New("foo")
	result := Suffix(original, "bar")
	if result.Error() != "foo bar" {
		t.Fatalf("bad: %s", result)
	}
}
