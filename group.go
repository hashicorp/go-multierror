// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package multierror

import (
	"context"
	"sync"
)

// Group is a collection of goroutines which return errors that need to be
// coalesced.
type Group struct {
	mutex  sync.Mutex
	err    *Error
	wg     sync.WaitGroup
	cancel context.CancelFunc
}

// Go calls the given function in a new goroutine.
//
// If the function returns an error it is added to the group multierror which
// is returned by Wait. If the group was created by WithContext, the associated
// context is canceled the first time a function returns a non-nil error, so
// that the remaining goroutines can stop early. Errors returned by the
// goroutines that already completed are still collected and returned by Wait.
func (g *Group) Go(f func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		if err := f(); err != nil {
			g.mutex.Lock()
			g.err = Append(g.err, err)
			g.mutex.Unlock()
			if g.cancel != nil {
				g.cancel()
			}
		}
	}()
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the multierror. If the group was created by WithContext, the
// associated context is canceled when Wait returns.
func (g *Group) Wait() *Error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	g.mutex.Lock()
	defer g.mutex.Unlock()
	return g.err
}

// WithContext returns a new Group and an associated Context derived from ctx.
//
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or the first time Wait returns, whichever occurs
// first. If the caller cancels ctx, the derived Context is canceled as well.
//
// Functions passed to Go should check the derived Context (typically with a
// select on ctx.Done()) and return as soon as it is canceled, in order to stop
// early. Unlike errgroup, which returns only the first error, all errors
// returned by the goroutines that completed before the cancellation are
// collected and returned by Wait.
func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{cancel: cancel}, ctx
}
