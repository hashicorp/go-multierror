// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package multierror

import (
	"errors"
	"strings"
	"testing"
)

func TestGroup(t *testing.T) {
	err1 := errors.New("group_test: 1")
	err2 := errors.New("group_test: 2")

	cases := []struct {
		errs      []error
		nilResult bool
	}{
		{errs: []error{}, nilResult: true},
		{errs: []error{nil}, nilResult: true},
		{errs: []error{err1}},
		{errs: []error{err1, nil}},
		{errs: []error{err1, nil, err2}},
	}

	for _, tc := range cases {
		var g Group

		for _, err := range tc.errs {
			err := err
			g.Go(func() error { return err })

		}

		gErr := g.Wait()
		if gErr != nil {
			for i := range tc.errs {
				if tc.errs[i] != nil && !strings.Contains(gErr.Error(), tc.errs[i].Error()) {
					t.Fatalf("expected error to contain %q, actual: %v", tc.errs[i].Error(), gErr)
				}
			}
		} else if !tc.nilResult {
			t.Fatalf("Group.Wait() should not have returned nil for errs: %v", tc.errs)
		}
	}
}

func TestGroupWait_ErrorNil(t *testing.T) {
	g := new(Group)
	g.Go(func() error { return nil })
	err := g.Wait()
	if err != nil {
		t.Fatalf("expected error to be nil, but was %v", err)
	}
}

func TestGroupWait_ErrorNotNil(t *testing.T) {
	g := new(Group)
	msg := "test error"
	g.Go(func() error { return errors.New(msg) })
	err := g.Wait()
	if err == nil {
		t.Fatalf("expected error to be nil, but was %v", err)
	}

	// err is a *Error, and e is set to the error's value
	var e *Error
	if !errors.As(err, &e) {
		t.Fatalf("expected err to be type *Error, but was type %T, value %v", err, err)
	}

	errs := e.WrappedErrors()
	if len(errs) != 1 {
		t.Fatalf("expected one wrapped error, but found %d", len(errs))
	}

	wrapped := errs[0]
	if wrapped.Error() != "test error" {
		t.Fatalf("expected wrap error message to be '%s', but was '%s'", msg, wrapped.Error())
	}
}
