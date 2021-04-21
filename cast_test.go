package multierror

import (
	"errors"
	"testing"
)

func TestCast_Nil(t *testing.T) {
	const expected = "*multierror.Error{Errors:[]error(nil), ErrorFormat:(multierror.ErrorFormatFunc)(nil)}"
	result := Cast(nil)
	if result.GoString() != expected {
		t.Fatalf("wrong result: %s", result.GoString())
	}
}

func TestCast_Error(t *testing.T) {
	result := Cast(errors.New("test message"))
	if len(result.Errors) != 1 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}
