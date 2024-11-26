// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package multierror

import (
	"fmt"
	"sync"
)

type token struct{}

// Group is a collection of goroutines which return errors that need to be
// coalesced.
type Group struct {
	mutex sync.Mutex
	err   *Error
	wg    sync.WaitGroup
	sem   chan token
}

// SetLimit limits the number of goroutines that can be concurrently active.
// A negative value indicates no limit.
func (g *Group) SetLimit(n int) {
	if n < 0 {
		g.sem = nil
		return
	}

	if g.sem != nil && len(g.sem) > 0 {
		panic(fmt.Errorf("multierror: modify limit while %v goroutines in the group are still active", len(g.sem)))
	}

	g.sem = make(chan token, n)
}

// Go calls the given function in a new goroutine.
//
// If the function returns an error it is added to the group multierror which
// is returned by Wait.
func (g *Group) Go(f func() error) {
	if g.sem != nil {
		g.sem <- token{}
	}

	g.wg.Add(1)

	go func() {
		defer func() {
			if g.sem != nil {
				<-g.sem
			}
			g.wg.Done()
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
