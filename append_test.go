package multierror

import (
	"errors"
	"testing"
)

func TestAppend_Error(t *testing.T) {
	original := &Error{
		Errors: []error{errors.New("foo")},
	}

	result := Append(original, errors.New("bar"))
	if len(result.Errors) != 2 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}

	original = &Error{}
	result = Append(original, errors.New("bar"))
	if len(result.Errors) != 1 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}

	// Test when a typed nil is passed
	var e *Error
	result = Append(e, errors.New("baz"))
	if len(result.Errors) != 1 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}

	// Test flattening
	original = &Error{
		Errors: []error{errors.New("foo")},
	}

	result = Append(original, Append(nil, errors.New("foo"), errors.New("bar")))
	if len(result.Errors) != 3 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}

func TestAppendNoFlatten_Error(t *testing.T) {
	original := &Error{
		Errors: []error{errors.New("foo")},
	}

	result := AppendNoFlatten(original, errors.New("bar"))
	if len(result.Errors) != 2 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}

	original = &Error{}
	result = AppendNoFlatten(original, errors.New("bar"))
	if len(result.Errors) != 1 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}

	// Test when a typed nil is passed
	var e *Error
	result = AppendNoFlatten(e, errors.New("baz"))
	if len(result.Errors) != 1 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}

	// Test flattening
	original = &Error{
		Errors: []error{errors.New("foo")},
	}

	result = AppendNoFlatten(original, AppendNoFlatten(nil, errors.New("foo"), errors.New("bar")))
	if len(result.Errors) != 2 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}

func TestAppend_NilError(t *testing.T) {
	var err error
	result := Append(err, errors.New("bar"))
	if len(result.Errors) != 1 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}

func TestAppendNoFlatten_NilError(t *testing.T) {
	var err error
	result := AppendNoFlatten(err, errors.New("bar"))
	if len(result.Errors) != 1 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}

	// Test flattening...
	inner := Append(errors.New("foo"), errors.New("bar"))

	// ...with nil error:
	result = AppendNoFlatten(nil, inner)
	if len(result.Errors) != 1 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}

	// ...with nil *Error:
	var nilErr *Error
	result = AppendNoFlatten(nilErr, inner)
	if len(result.Errors) != 1 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}

func TestAppend_NilErrorArg(t *testing.T) {
	var err error
	var nilErr *Error
	result := Append(err, nilErr)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}

func TestAppendNoFlatten_NilErrorArg(t *testing.T) {
	var err error
	var nilErr *Error
	result := AppendNoFlatten(err, nilErr)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}

func TestAppend_NilErrorIfaceArg(t *testing.T) {
	var err error
	var nilErr error
	result := Append(err, nilErr)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}

func TestAppendNoFlatten_NilErrorIfaceArg(t *testing.T) {
	var err error
	var nilErr error
	result := AppendNoFlatten(err, nilErr)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}

func TestAppend_NonError(t *testing.T) {
	original := errors.New("foo")
	result := Append(original, errors.New("bar"))
	if len(result.Errors) != 2 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}

func TestAppendNoFlatten_NonError(t *testing.T) {
	original := errors.New("foo")
	result := AppendNoFlatten(original, errors.New("bar"))
	if len(result.Errors) != 2 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}

func TestAppend_NonError_Error(t *testing.T) {
	original := errors.New("foo")
	result := Append(original, Append(nil, errors.New("bar")))
	if len(result.Errors) != 2 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}

func TestAppendNoFlatten_NonError_Error(t *testing.T) {
	original := errors.New("foo")
	result := AppendNoFlatten(original, AppendNoFlatten(nil, errors.New("bar")))
	if len(result.Errors) != 2 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}
