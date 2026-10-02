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
	cancel func()

	mutex sync.Mutex
	err   *Error
	wg    sync.WaitGroup
}

// GroupWithContext returns a new Group and an associated Context derived from ctx.
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or the first time Wait returns, whichever occurs first.
func GroupWithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{cancel: cancel}, ctx
}

// Go calls the given function in a new goroutine.
//
// If the function returns an error it is added to the group multierror which
// is returned by Wait.
func (g *Group) Go(f func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		if err := f(); err != nil {
			g.mutex.Lock()
			g.err = Append(g.err, err)
			if g.cancel != nil {
				g.cancel()
			}
			g.mutex.Unlock()
		}
	}()
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the multierror.
func (g *Group) Wait() *Error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	g.mutex.Lock()
	defer g.mutex.Unlock()
	return g.err
}
