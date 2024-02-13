// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package multierror

import (
	"errors"
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

func TestGroupSetLimit(t *testing.T) {
	var (
		active    int32
		maxActive int32
	)

	g := &Group{}
	g.SetLimit(2)

	work := func() error {
		atomic.AddInt32(&active, 1)

		for {
			currentMax := atomic.LoadInt32(&maxActive)
			currentActive := atomic.LoadInt32(&active)
			if currentActive > currentMax {
				if atomic.CompareAndSwapInt32(&maxActive, currentMax, currentActive) {
					break
				}
			} else {
				break
			}
		}

		time.Sleep(200 * time.Millisecond)

		atomic.AddInt32(&active, -1)

		return nil
	}

	// Start more goroutines than the limit
	for i := 0; i < 5; i++ {
		g.Go(work)
	}

	err := g.Wait()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if maxActive != 2 {
		t.Errorf("expected max 2 active goroutines, got %d", maxActive)
	}

	g = &Group{}
	g.SetLimit(-1)

	// Test unlimited
	for i := 0; i < 10; i++ {
		g.Go(work)
	}

	err = g.Wait()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if maxActive != 10 {
		t.Errorf("expected max 2 active goroutines, got %d", maxActive)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic when modifying limit, got none")
		}
	}()

	g = &Group{}

	g.SetLimit(2)

	g.Go(work)

	g.SetLimit(3) // attempt to modify limit while goroutine is active
}
