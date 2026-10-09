// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package multierror

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"testing"
	"time"
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

func TestGroupWithContext_ErrorCancellation(t *testing.T) {
	g, ctx := GroupWithContext(context.Background())

	errExpected := errors.New("first task failure")

	g.Go(func() error {
		return errExpected
	})

	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Second):
			return errors.New("timed out waiting for context cancellation")
		}
	})

	gErr := g.Wait()
	if gErr == nil {
		t.Fatal("expected error, got nil")
	}

	if len(gErr.Errors) != 2 {
		t.Fatalf("expected 2 errors, got %d: %v", len(gErr.Errors), gErr.Errors)
	}

	if !errors.Is(gErr, errExpected) {
		t.Fatalf("expected multierror to contain %v", errExpected)
	}

	if !errors.Is(gErr, context.Canceled) {
		t.Fatalf("expected multierror to contain context.Canceled")
	}
}

func TestGroupWithContext_Success(t *testing.T) {
	g, ctx := GroupWithContext(context.Background())

	var counter int64
	for i := 0; i < 5; i++ {
		g.Go(func() error {
			atomic.AddInt64(&counter, 1)
			return nil
		})
	}

	gErr := g.Wait()
	if gErr != nil {
		t.Fatalf("expected nil error, got %v", gErr)
	}

	if atomic.LoadInt64(&counter) != 5 {
		t.Fatalf("expected counter 5, got %d", atomic.LoadInt64(&counter))
	}

	select {
	case <-ctx.Done():
	default:
		t.Fatal("expected context to be canceled after Wait returns")
	}
}

func TestGroupWithContext_MultipleErrors(t *testing.T) {
	g, _ := GroupWithContext(context.Background())

	taskCount := 10
	for i := 0; i < taskCount; i++ {
		i := i
		g.Go(func() error {
			return fmt.Errorf("task error %d", i)
		})
	}

	gErr := g.Wait()
	if gErr == nil {
		t.Fatal("expected errors, got nil")
	}

	if len(gErr.Errors) != taskCount {
		t.Fatalf("expected %d errors, got %d", taskCount, len(gErr.Errors))
	}

	for i := 0; i < taskCount; i++ {
		expectedSubstr := fmt.Sprintf("task error %d", i)
		if !strings.Contains(gErr.Error(), expectedSubstr) {
			t.Fatalf("expected error output to contain %q, got: %s", expectedSubstr, gErr.Error())
		}
	}
}

func TestGroupWithContext_PreCancelledContext(t *testing.T) {
	parentCtx, parentCancel := context.WithCancel(context.Background())
	parentCancel()

	g, ctx := GroupWithContext(parentCtx)

	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return errors.New("context was expected to be cancelled immediately")
		}
	})

	gErr := g.Wait()
	if gErr == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(gErr, context.Canceled) {
		t.Fatalf("expected context.Canceled error, got: %v", gErr)
	}
}

func TestWithContext_Alias(t *testing.T) {
	g, ctx := WithContext(context.Background())

	errAlias := errors.New("alias failure")
	g.Go(func() error {
		return errAlias
	})

	gErr := g.Wait()
	if gErr == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(gErr, errAlias) {
		t.Fatalf("expected %v, got %v", errAlias, gErr)
	}

	select {
	case <-ctx.Done():
	default:
		t.Fatal("expected context to be canceled")
	}
}

func TestGroupWithContext_WaitWithoutGo(t *testing.T) {
	g, ctx := GroupWithContext(context.Background())

	gErr := g.Wait()
	if gErr != nil {
		t.Fatalf("expected nil error, got %v", gErr)
	}

	select {
	case <-ctx.Done():
	default:
		t.Fatal("expected context to be canceled after Wait")
	}
}
