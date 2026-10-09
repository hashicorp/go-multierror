// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package multierror

import (
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
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

func TestGroup_Panic(t *testing.T) {
	var g Group

	errPanic := errors.New("custom panic error")
	normalErr := errors.New("normal error")

	g.Go(func() error {
		panic(errPanic)
	})
	g.Go(func() error {
		panic("string panic")
	})
	g.Go(func() error {
		panic(42)
	})
	g.Go(func() error {
		return normalErr
	})
	g.Go(func() error {
		return nil
	})

	gErr := g.Wait()
	if gErr == nil {
		t.Fatal("expected non-nil error from Wait()")
	}

	if len(gErr.Errors) != 4 {
		t.Fatalf("expected 4 errors, got %d", len(gErr.Errors))
	}

	if !errors.Is(gErr, errPanic) {
		t.Fatalf("expected gErr to unwrap to errPanic, got: %v", gErr)
	}

	if !errors.Is(gErr, normalErr) {
		t.Fatalf("expected gErr to unwrap to normalErr, got: %v", gErr)
	}

	var panicErr *PanicError
	if !errors.As(gErr, &panicErr) {
		t.Fatal("expected errors.As to match *PanicError")
	}

	errStr := gErr.Error()
	if !strings.Contains(errStr, "panic: custom panic error") {
		t.Fatalf("expected error to contain %q, actual: %v", "panic: custom panic error", errStr)
	}
	if !strings.Contains(errStr, "panic: string panic") {
		t.Fatalf("expected error to contain %q, actual: %v", "panic: string panic", errStr)
	}
	if !strings.Contains(errStr, "panic: 42") {
		t.Fatalf("expected error to contain %q, actual: %v", "panic: 42", errStr)
	}
	if !strings.Contains(errStr, "normal error") {
		t.Fatalf("expected error to contain %q, actual: %v", "normal error", errStr)
	}
}

func TestGroup_PanicHandler(t *testing.T) {
	var g Group
	var handledCount int32

	g.PanicHandler = func(v interface{}) error {
		atomic.AddInt32(&handledCount, 1)
		return fmt.Errorf("custom-handler: %v", v)
	}

	g.Go(func() error {
		panic("handled panic")
	})

	gErr := g.Wait()
	if gErr == nil {
		t.Fatal("expected non-nil error from Wait()")
	}

	if got := atomic.LoadInt32(&handledCount); got != 1 {
		t.Fatalf("expected handledCount to be 1, got %d", got)
	}

	if !strings.Contains(gErr.Error(), "custom-handler: handled panic") {
		t.Fatalf("expected error to contain custom-handler output, actual: %v", gErr)
	}
}

func TestGroup_PanicHandler_NilError(t *testing.T) {
	var g Group

	g.PanicHandler = func(v interface{}) error {
		return nil
	}

	g.Go(func() error {
		panic("ignored panic")
	})
	g.Go(func() error {
		return nil
	})

	gErr := g.Wait()
	if gErr != nil {
		t.Fatalf("expected nil error when panic handler returns nil, got: %v", gErr)
	}
}

func TestGroup_ConcurrentPanicsAndErrors(t *testing.T) {
	var g Group

	numWorkers := 60
	for i := 0; i < numWorkers; i++ {
		idx := i
		g.Go(func() error {
			switch idx % 3 {
			case 0:
				return nil
			case 1:
				return fmt.Errorf("err-%d", idx)
			default:
				panic(fmt.Sprintf("panic-%d", idx))
			}
		})
	}

	gErr := g.Wait()
	if gErr == nil {
		t.Fatal("expected non-nil error from Wait()")
	}

	expectedErrors := (numWorkers / 3) * 2
	if len(gErr.Errors) != expectedErrors {
		t.Fatalf("expected %d errors, got %d", expectedErrors, len(gErr.Errors))
	}
}

func TestPanicError_Methods(t *testing.T) {
	var nilPE *PanicError
	if nilPE.Error() != "" {
		t.Fatalf("expected empty string for nil PanicError, got: %q", nilPE.Error())
	}
	if nilPE.Unwrap() != nil {
		t.Fatalf("expected nil unwrap for nil PanicError, got: %v", nilPE.Unwrap())
	}

	innerErr := errors.New("inner error")
	peWithErr := &PanicError{Value: innerErr}
	if peWithErr.Error() != "panic: inner error" {
		t.Fatalf("unexpected Error() output: %q", peWithErr.Error())
	}
	if !errors.Is(peWithErr, innerErr) {
		t.Fatal("expected peWithErr to unwrap to innerErr")
	}

	peWithString := &PanicError{Value: "raw text"}
	if peWithString.Error() != "panic: raw text" {
		t.Fatalf("unexpected Error() output: %q", peWithString.Error())
	}
	if peWithString.Unwrap() != nil {
		t.Fatalf("expected nil unwrap when Value is not error, got: %v", peWithString.Unwrap())
	}
}

