// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package multierror

import (
	"fmt"
	"sync"
)

// Group is a collection of goroutines which return errors that need to be
// coalesced.
type Group struct {
	mutex sync.Mutex
	err   *Error
	wg    sync.WaitGroup

	// PanicHandler is an optional function called when a goroutine started by Go
	// recovers from a panic. If set, the returned error is added to the group
	// multierror. If nil, the recovered panic is wrapped in a *PanicError.
	PanicHandler func(interface{}) error
}

// PanicError represents a panic recovered during goroutine execution in a Group.
type PanicError struct {
	Value interface{}
}

// Error implements the error interface.
func (p *PanicError) Error() string {
	if p == nil {
		return ""
	}
	return fmt.Sprintf("panic: %v", p.Value)
}

// Unwrap returns the underlying error if Value implements error, or nil otherwise.
func (p *PanicError) Unwrap() error {
	if p == nil {
		return nil
	}
	if err, ok := p.Value.(error); ok {
		return err
	}
	return nil
}

// Go calls the given function in a new goroutine.
//
// If the function returns an error or panics, the error (or recovered panic)
// is added to the group multierror which is returned by Wait.
func (g *Group) Go(f func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()
		defer func() {
			if r := recover(); r != nil {
				var err error
				if g.PanicHandler != nil {
					err = g.PanicHandler(r)
				} else {
					err = &PanicError{Value: r}
				}

				if err != nil {
					g.mutex.Lock()
					g.err = Append(g.err, err)
					g.mutex.Unlock()
				}
			}
		}()

		if err := f(); err != nil {
			g.mutex.Lock()
			g.err = Append(g.err, err)
			g.mutex.Unlock()
		}
	}()
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the multierror.
func (g *Group) Wait() *Error {
	g.wg.Wait()
	g.mutex.Lock()
	defer g.mutex.Unlock()
	return g.err
}
