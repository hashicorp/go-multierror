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

func TestGroup_GoWithRecover_Panic(t *testing.T) {
	var g Group

	errPanic := errors.New("custom panic error")

	g.GoWithRecover(func() error {
		panic(errPanic)
	})
	g.GoWithRecover(func() error {
		panic("string panic")
	})
	g.GoWithRecover(func() error {
		return nil
	})

	gErr := g.Wait()
	if gErr == nil {
		t.Fatal("expected non-nil error from Wait()")
	}

	if len(gErr.Errors) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(gErr.Errors))
	}

	if !errors.Is(gErr, errPanic) {
		t.Fatalf("expected gErr to unwrap to errPanic, got: %v", gErr)
	}

	errStr := gErr.Error()
	if !strings.Contains(errStr, "panic: custom panic error") {
		t.Fatalf("expected error to contain %q, actual: %v", "panic: custom panic error", errStr)
	}
	if !strings.Contains(errStr, "panic: string panic") {
		t.Fatalf("expected error to contain %q, actual: %v", "panic: string panic", errStr)
	}
}

func TestGroup_GoWithRecover_Error(t *testing.T) {
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
			g.GoWithRecover(func() error { return err })
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

