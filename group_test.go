// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package multierror

import (
	"context"
	"errors"
	"strings"
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

func TestGroupWithContextCancelOnFirstError(t *testing.T) {
	ctx, done := context.WithCancel(context.Background())
	defer done()
	g, gctx := WithContext(ctx)

	firstErr := errors.New("group_test: first failure")

	started := make(chan struct{})
	released := make(chan struct{})

	// The first goroutine fails after a small delay; its error must cancel
	// the derived context.
	g.Go(func() error {
		close(started)
		<-released
		return firstErr
	})

	// The second goroutine respects the derived context and stops early.
	g.Go(func() error {
		<-gctx.Done()
		return nil
	})

	<-started
	close(released)

	gErr := g.Wait()
	if gErr == nil {
		t.Fatal("Group.Wait() should have returned the first goroutine error")
	}
	if !strings.Contains(gErr.Error(), firstErr.Error()) {
		t.Fatalf("expected aggregated error to contain %q, got: %v", firstErr.Error(), gErr)
	}
	if err := gctx.Err(); err == nil {
		t.Fatal("derived context should be canceled after the first error")
	}
}

func TestGroupWithContextCollectsAllErrors(t *testing.T) {
	ctx, done := context.WithCancel(context.Background())
	defer done()
	g, gctx := WithContext(ctx)

	err1 := errors.New("group_test: error one")
	err2 := errors.New("group_test: error two")

	released := make(chan struct{})

	started := make(chan struct{}, 2)
	for _, err := range []error{err1, err2} {
		err := err
		g.Go(func() error {
			started <- struct{}{}
			<-released
			return err
		})
	}

	// Wait until both goroutines are running, then let them all fail.
	<-started
	<-started
	close(released)

	gErr := g.Wait()
	if gErr == nil {
		t.Fatal("Group.Wait() should have returned errors")
	}
	for _, want := range []string{err1.Error(), err2.Error()} {
		if !strings.Contains(gErr.Error(), want) {
			t.Fatalf("expected aggregated error to contain %q, got: %v", want, gErr)
		}
	}
	if err := gctx.Err(); err == nil {
		t.Fatal("derived context should be canceled after errors occur")
	}
}

func TestGroupWithContextParentCancel(t *testing.T) {
	parent, cancelParent := context.WithCancel(context.Background())
	defer cancelParent()

	g, gctx := WithContext(parent)

	aborted := make(chan struct{})
	g.Go(func() error {
		<-gctx.Done()
		close(aborted)
		return nil
	})

	select {
	case <-gctx.Done():
		t.Fatal("derived context should not be canceled before the parent")
	default:
	}

	cancelParent()

	select {
	case <-aborted:
	case <-time.After(5 * time.Second):
		t.Fatal("goroutine did not observe the parent cancellation")
	}

	if err := gctx.Err(); err == nil {
		t.Fatal("derived context should be canceled when the parent is canceled")
	}
	if gErr := g.Wait(); gErr != nil {
		t.Fatalf("expected no aggregated error, got: %v", gErr)
	}
}

func TestGroupWithContextWaitCancels(t *testing.T) {
	ctx, done := context.WithCancel(context.Background())
	defer done()
	g, gctx := WithContext(ctx)

	if g.Wait() != nil {
		t.Fatal("expected nil error from group with no work")
	}
	if err := gctx.Err(); err == nil {
		t.Fatal("derived context should be canceled when Wait returns")
	}
}

func TestGroupWithContextReturnsContextError(t *testing.T) {
	parent, cancelParent := context.WithCancel(context.Background())
	defer cancelParent()
	g, gctx := WithContext(parent)

	// A goroutine that is canceled mid-work may choose to surface the
	// context error; it is aggregated like any other error.
	g.Go(func() error {
		<-gctx.Done()
		return gctx.Err()
	})

	cancelParent()

	gErr := g.Wait()
	if gErr == nil {
		t.Fatal("Group.Wait() should have returned the context error")
	}
	if !strings.Contains(gErr.Error(), context.Canceled.Error()) {
		t.Fatalf("expected aggregated error to contain %q, got: %v", context.Canceled.Error(), gErr)
	}
}
