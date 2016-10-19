package multierror

import (
	"errors"
	"testing"
)

type TestError struct{}

func (e TestError) Error() string { return "TestError" }

func TestRemoveNils(t *testing.T) {
	errs := []error{errors.New("foo"), nil, nil, errors.New("foo"), nil}
	errs = removeNils(errs)
	if len(errs) != 2 {
		t.Fatalf("wrong len: %d", len(errs))
	}
}

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

func TestAppend_NilError(t *testing.T) {
	var err error
	result := Append(err, errors.New("bar"))
	if len(result.Errors) != 1 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}

func TestAppend_NilNil(t *testing.T) {
	var err error
	result := AppendNonNil(err, nil)
	if result != nil {
		t.Fatalf("non-nil errors: %s", result.Error())
	}
}

func TestAppendNonNil(t *testing.T) {
	var err1 error
	var err2 *Error
	result := AppendNonNil(err1, err2, nil, nil)
	if result != nil {
		t.Fatalf("non-nil errors: %s", result.Error())
	}
	err1 = errors.New("foo")
	result = AppendNonNil(err1, err2, nil, nil)
	if result != err1 {
		t.Fatalf("input error modified: %s", result.Error())
	}
	err1 = errors.New("foo")
	result = AppendNonNil(nil, err1, nil, nil)
	if result != err1 {
		t.Fatalf("input error modified: %s", result.Error())
	}
}

func TestAppendNonNilStruct(t *testing.T) {
	var err1 error
	var err3 TestError
	result := AppendNonNil(err1, err3, nil)
	if result == nil {
		t.Fatalf("TestError was not appended")
	}
}

func TestAppend_NonError(t *testing.T) {
	original := errors.New("foo")
	result := Append(original, errors.New("bar"))
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
