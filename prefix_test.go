// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package multierror

import (
	"errors"
	"testing"
)

func TestPrefix_Error(t *testing.T) {
	original := &Error{
		Errors: []error{errors.New("foo")},
	}

	result := Prefix(original, "bar")
	if result.(*Error).Errors[0].Error() != "bar foo" {
		t.Fatalf("bad: %s", result)
	}
}

func TestPrefix_NilError(t *testing.T) {
	var err error
	result := Prefix(err, "bar")
	if result != nil {
		t.Fatalf("bad: %#v", result)
	}
}

func TestPrefix_NonError(t *testing.T) {
	original := errors.New("foo")
	result := Prefix(original, "bar")
	if result == nil {
		t.Fatal("error result was nil")
	}
	if result.Error() != "bar foo" {
		t.Fatalf("bad: %s", result)
	}
}

func TestPrefix_ErrorsIs(t *testing.T) {
	// errors.Is should work through errors returned by Prefix.
	sentinel := errors.New("sentinel")

	// Non-multierror case
	result := Prefix(sentinel, "context:")
	if !errors.Is(result, sentinel) {
		t.Fatalf("errors.Is should find sentinel through Prefix result, got %s", result)
	}

	// multierror case
	multi := &Error{Errors: []error{sentinel}}
	prefixed := Prefix(multi, "context:")
	if !errors.Is(prefixed, sentinel) {
		t.Fatalf("errors.Is should find sentinel through Prefix multierror result, got %s", prefixed)
	}
}
