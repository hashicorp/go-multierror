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

func TestGroupWithContext(t *testing.T) {
	t.Run("cancels on error and accumulates all errors", func(t *testing.T) {
		g, ctx := GroupWithContext(context.Background())

		err1 := errors.New("group_test: 1")
		err2 := errors.New("group_test: 2")

		g.Go(func() error {
			return err1
		})

		g.Go(func() error {
			select {
			case <-ctx.Done():
				return err2
			case <-time.After(2 * time.Second):
				t.Error("context was not canceled in time")
				return errors.New("timeout")
			}
		})

		gErr := g.Wait()

		select {
		case <-ctx.Done():
			if !errors.Is(ctx.Err(), context.Canceled) {
				t.Fatalf("expected context.Canceled, got %v", ctx.Err())
			}
		default:
			t.Fatal("expected context to be canceled")
		}

		if gErr == nil {
			t.Fatal("expected non-nil multierror")
		}
		if len(gErr.Errors) != 2 {
			t.Fatalf("expected 2 errors, got %d: %v", len(gErr.Errors), gErr.Errors)
		}
		if !strings.Contains(gErr.Error(), err1.Error()) {
			t.Fatalf("expected error to contain %q, got: %v", err1.Error(), gErr)
		}
		if !strings.Contains(gErr.Error(), err2.Error()) {
			t.Fatalf("expected error to contain %q, got: %v", err2.Error(), gErr)
		}
	})

	t.Run("completes when all functions succeed", func(t *testing.T) {
		g, ctx := GroupWithContext(context.Background())

		g.Go(func() error {
			return nil
		})
		g.Go(func() error {
			return nil
		})

		gErr := g.Wait()
		if gErr != nil {
			t.Fatalf("expected nil error, got %v", gErr)
		}

		select {
		case <-ctx.Done():
			if !errors.Is(ctx.Err(), context.Canceled) {
				t.Fatalf("expected context.Canceled, got %v", ctx.Err())
			}
		default:
			t.Fatal("expected context to be canceled after Wait")
		}
	})
}
