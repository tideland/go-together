// Tideland Go Together - Actor
//
// Copyright (C) 2017-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package actor // import "tideland.dev/go/together/actor"

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

//--------------------
// CONSTANTS
//--------------------

const (
	// DefaultTimeout is used in a DoSync() call.
	DefaultTimeout = 5 * time.Second

	// DefaultQueueCap is the minimum and default capacity
	// of the async actions queue.
	DefaultQueueCap = 256
)

//--------------------
// FUNCTION TYPES
//--------------------

// Action defines the signature of an actor action.
type Action func()

// Repairer allows the Actor to react on a panic during its
// work. If it returns nil the backend shall continue
// work. Otherwise the error is stored and the backend
// terminated.
type Repairer func(reason interface{}) error

// Finalizer is called with the Actors internal status when
// the backend loop terminates.
type Finalizer func(err error) error

//--------------------
// ACTOR
//--------------------

// Actor allows to simply use and control a goroutine.
type Actor struct {
	mu           sync.Mutex
	ctx          context.Context
	cancel       func()
	asyncActions chan Action
	syncActions  chan Action
	repairer     Repairer
	finalizer    Finalizer
	works        atomic.Value
	err          error
}

// Go starts an Actor with the passed options.
func Go(options ...Option) (*Actor, error) {
	// Init with options.
	act := &Actor{
		syncActions: make(chan Action),
	}
	act.works.Store(true)
	for _, option := range options {
		if err := option(act); err != nil {
			return nil, err
		}
	}
	// Ensure default settings.
	if act.ctx == nil {
		act.ctx, act.cancel = context.WithCancel(context.Background())
	} else {
		act.ctx, act.cancel = context.WithCancel(act.ctx)
	}
	if act.asyncActions == nil {
		act.asyncActions = make(chan Action, DefaultQueueCap)
	}
	// Create loop with its options.
	started := make(chan struct{})
	go act.backend(started)
	select {
	case <-started:
		return act, nil
	case <-time.After(DefaultTimeout):
		return nil, fmt.Errorf("actor starting timeout after %.1f seconds", DefaultTimeout.Seconds())
	}
}

// DoAsync send the actor function to the backend and returns
// when it's queued.
func (act *Actor) DoAsync(action Action) error {
	return act.DoAsyncTimeout(action, DefaultTimeout)
}

// DoAsyncTimeout send the actor function to the backend and returns
// when it's queued.
func (act *Actor) DoAsyncTimeout(action Action, timeout time.Duration) error {
	act.mu.Lock()
	if act.err != nil {
		act.mu.Unlock()
		return act.err
	}
	if !act.works.Load().(bool) {
		act.mu.Unlock()
		return fmt.Errorf("actor doesn't work anymore")
	}
	act.mu.Unlock()
	select {
	case act.asyncActions <- action:
	case <-time.After(timeout):
		return fmt.Errorf("timeout")
	}
	return nil
}

// DoSync executes the actor function and returns when it's done
// or it has the default timeout.
func (act *Actor) DoSync(action Action) error {
	return act.DoSyncTimeout(action, DefaultTimeout)
}

// DoSyncTimeout executes the action and returns when it's done
// or it has a timeout.
func (act *Actor) DoSyncTimeout(action Action, timeout time.Duration) error {
	act.mu.Lock()
	if act.err != nil {
		act.mu.Unlock()
		return act.err
	}
	if !act.works.Load().(bool) {
		act.mu.Unlock()
		return fmt.Errorf("actor doesn't work anymore")
	}
	act.mu.Unlock()
	done := make(chan struct{})
	syncAction := func() {
		action()
		close(done)
	}
	select {
	case act.syncActions <- syncAction:
	case <-time.After(timeout):
		return fmt.Errorf("timeout")
	}
	select {
	case <-done:
	case <-time.After(timeout):
		if !act.works.Load().(bool) {
			return act.err
		}
		return fmt.Errorf("timeout")
	}
	return nil
}

// Err returns information if the Actor has an error.
func (act *Actor) Err() error {
	act.mu.Lock()
	defer act.mu.Unlock()
	return act.err
}

// Stop terminates the Actor backend.
func (act *Actor) Stop() {
	act.mu.Lock()
	defer act.mu.Unlock()
	if !act.works.Load().(bool) {
		// Already stopped.
		return
	}
	act.works.Store(false)
	act.cancel()
}

// backend runs the goroutine of the Actor.
func (act *Actor) backend(started chan struct{}) {
	defer act.finalize()
	close(started)
	for act.works.Load().(bool) {
		act.work()
	}
}

// work runs the select in a loop, including
// a possible repairer.
func (act *Actor) work() {
	defer func() {
		// Check and handle panics!
		reason := recover()
		switch {
		case reason != nil && act.repairer != nil:
			// Try to repair.
			err := act.repairer(reason)
			act.mu.Lock()
			act.err = err
			act.works.Store(act.err == nil)
			act.mu.Unlock()
		case reason != nil && act.repairer == nil:
			// Accept panic.
			act.mu.Lock()
			act.err = fmt.Errorf("actor panic: %v", reason)
			act.works.Store(false)
			act.mu.Unlock()
		}
	}()
	// Select in loop.
	for {
		select {
		case <-act.ctx.Done():
			return
		case action := <-act.asyncActions:
			action()
		case action := <-act.syncActions:
			action()
		}
	}
}

// finalize takes care for a clean loop finalization.
func (act *Actor) finalize() {
	act.mu.Lock()
	defer act.mu.Unlock()
	if act.finalizer != nil {
		act.err = act.finalizer(act.err)
	}
}

// EOF
