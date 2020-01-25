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
	"time"

	"tideland.dev/go/together/fuse"
	"tideland.dev/go/trace/failure"
)

//--------------------
// CONSTANTS
//--------------------

const (
	// DefaultTimeout is used in a DoSync() call.
	DefaultTimeout = 5 * time.Second

	// DefaultQueueCap is the minimum and default capacity
	// of the async actions queue.
	DefaultQueueCap = 64
)

//--------------------
// FUNCTION TYPES
//--------------------

// Action defines the signature of an actor action.
type Action func()

// Recoverer allows the actor to react on a panic during its
// work. If it returns nil the backend shall continue
// work. Otherwise the error is stored and the backend
// terminated.
type Recoverer func(reason interface{}) error

// Finalizer is called with the actors internal status when
// the backend loop terminates.
type Finalizer func(err error) error

//--------------------
// ACTOR
//--------------------

// Actor allows to simply use and control a goroutine.
type Actor struct {
	ctx          context.Context
	cancel       func()
	signal       *fuse.Signal
	asyncActions chan Action
	syncActions  chan Action
	recoverer    Recoverer
	finalizer    Finalizer
	err          failure.Error
}

// Go starts an Actor with the passed options.
func Go(options ...Option) (*Actor, error) {
	// Init with options.
	act := &Actor{
		signal:      fuse.NewSignal(),
		syncActions: make(chan Action),
	}
	for _, option := range options {
		if err := option(act); err != nil {
			act.err.Set(err)
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
	if act.recoverer == nil {
		act.recoverer = func(reason interface{}) error {
			return fmt.Errorf("actor panic: %v", reason)
		}
	}
	if act.finalizer == nil {
		act.finalizer = func(err error) error {
			return err
		}
	}
	// Create loop with its options.
	act.signal.Notify(fuse.Starting)
	go act.backend()
	if err := act.signal.Wait(fuse.Ready, DefaultTimeout); err != nil {
		return nil, err
	}
	return act, nil
}

// DoAsync send the actor function to the backend and returns
// when it's queued.
func (act *Actor) DoAsync(action Action) error {
	if !act.err.IsNil() {
		return act.err.Get()
	}
	act.asyncActions <- action
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
	if !act.err.IsNil() {
		return act.err.Get()
	}
	done := make(chan struct{})
	act.syncActions <- func() {
		action()
		close(done)
	}
	select {
	case <-done:
	case <-time.After(timeout):
		if !act.err.IsNil() {
			return act.err.Get()
		}
		return failure.New("synchronous actor do: timeout")
	}
	return nil
}

// Err returns information if the Actor has an error.
func (act *Actor) Err() error {
	return act.err.Get()
}

// Kill terminates the Actor backend with a given external error.
func (act *Actor) Kill(err error) error {
	if !act.err.IsNil() {
		return act.err.Get()
	}
	act.err.Set(err)
	act.cancel()
	if err := act.signal.Wait(fuse.Stopped, DefaultTimeout); err != nil {
		return err
	}
	return act.err.Get()
}

// Stop terminates the Actor backend.
func (act *Actor) Stop() error {
	return act.Kill(nil)
}

// backend runs the goroutine of the Actor.
func (act *Actor) backend() {
	defer func() {
		act.err.Set(act.finalizer(act.err.Get()))
		act.signal.Notify(fuse.Stopped)
	}()
	act.signal.Notify(fuse.Ready)
	for act.loop() {
	}
}

// loop runs the select in a loop, including
// a possible recoverer.
func (act *Actor) loop() (ok bool) {
	defer func() {
		if reason := recover(); reason != nil {
			// Panic!
			err := act.recoverer(reason)
			if err != nil {
				act.err.Set(err)
				ok = false
			} else {
				ok = true
			}
		} else {
			// Regular ending.
			ok = false
		}
		if !ok {
			act.signal.Notify(fuse.Stopping)
		}
	}()
	runs := 0
	for {
		runs++
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

// EOF
