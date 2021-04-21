package multierror

import (
	"errors"
	"testing"
)

func TestAppend_Error(t *testing.T) {
	original := &Error{
		Errors: []error{errors.New("foo")},
	}

	err := Append(original, errors.New("bar"))
	result := Cast(err)
	if len(result.Errors) != 2 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}

	original = &Error{}
	err = Append(original, errors.New("bar"))
	result = Cast(err)
	if len(result.Errors) != 1 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}

	// Test when a typed nil is passed
	var e *Error
	err = Append(e, errors.New("baz"))
	result = Cast(err)
	if len(result.Errors) != 1 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}

	// Test flattening
	original = &Error{
		Errors: []error{errors.New("foo")},
	}

	err = Append(original, Append(nil, errors.New("foo"), errors.New("bar")))
	result = Cast(err)
	if len(result.Errors) != 3 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}

func TestAppend_NilError(t *testing.T) {
	var err error
	err = Append(err, errors.New("bar"))
	result := Cast(err)
	if len(result.Errors) != 1 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}

func TestAppend_NilErrorArg(t *testing.T) {
	var err error
	var nilErr *Error
	err = Append(err, nilErr)
	result := Cast(err)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}

func TestAppend_NilErrorIfaceArg(t *testing.T) {
	var err error
	var nilErr error
	err = Append(err, nilErr)
	result := Cast(err)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}

func TestAppend_NonError(t *testing.T) {
	original := errors.New("foo")
	err := Append(original, errors.New("bar"))
	result := Cast(err)
	if len(result.Errors) != 2 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}

func TestAppend_NonError_Error(t *testing.T) {
	original := errors.New("foo")
	err := Append(original, Append(nil, errors.New("bar")))
	result := Cast(err)
	if len(result.Errors) != 2 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}
