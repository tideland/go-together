// Tideland Go Together - Notifier
//
// Copyright (C) 2017-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package notifier // import "tideland.dev/go/together/notifier"

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"reflect"
	"sync"
)

//--------------------
// CLOSER
//--------------------

// Closer signals a typical for-select-loop to terminate. It listens
// to multiple structs channels itself, e.g. from a context or other
// termination signalling functions.
type Closer struct {
	mu        sync.RWMutex
	startOnce sync.Once
	closeOnce sync.Once
	cases     []reflect.SelectCase
	done      bool
	doneC     chan struct{}
}

// NewCloser creates a new Closer instance.
func NewCloser(closeCs ...<-chan struct{}) *Closer {
	c := &Closer{
		doneC: make(chan struct{}),
	}
	for _, closeC := range closeCs {
		c.cases = append(c.cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(closeC),
		})
	}
	c.startOnce.Do(c.goWaiting)
	return c
}

// Add appends more channels to the closer.
func (c *Closer) Add(closeCs ...<-chan struct{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, closeC := range closeCs {
		c.cases = append(c.cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(closeC),
		})
	}
}

// Close allows to directly close the closer.
func (c *Closer) Close() {
	c.closeOnce.Do(c.notify)
}

// Done returns a channel that closes the Closer user has to end.
func (c *Closer) Done() <-chan struct{} {
	dC := c.doneC
	return dC
}

// String implements fmt.Stringer.
func (c *Closer) String() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.done {
		return fmt.Sprintf("Closer{%d cases; done}", len(c.cases))
	}
	return fmt.Sprintf("Closer{%d cases; waiting}", len(c.cases))
}

// goWaiting starts wait() as goroutine.
func (c *Closer) goWaiting() {
	go c.wait()
}

// wait is the backend goroutine waiting for closing of
// one of the channels.
func (c *Closer) wait() {
	reflect.Select(c.cases)
	c.closeOnce.Do(c.notify)
}

// notify tells the waiting users that the closer has done.
func (c *Closer) notify() {
	c.done = true
	close(c.doneC)
}

// EOF
