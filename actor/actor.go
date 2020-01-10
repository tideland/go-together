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
	"sync"
	"time"

	"tideland.dev/go/together/loop"
	"tideland.dev/go/together/notifier"
	"tideland.dev/go/trace/failure"
)

//--------------------
// CONSTANTS
//--------------------

const (
	// DefaultTimeout is used in a DoSync() call.
	DefaultTimeout = 5 * time.Second
)

//--------------------
// RECOVERER
//--------------------

// Recoverer allows a goroutine to react on a panic during its
// work. If it returns nil the goroutine shall continue
// work. Otherwise it will return with an error the gouroutine
// may use for its continued processing.
type Recoverer func(reason interface{}) error

//--------------------
// ACTOR
//--------------------

// Action defines the signature of an actor action.
type Action func() error

// Actor allows to simply use and control a goroutine.
type Actor struct {
	mu      sync.RWMutex
	actionC chan Action
	options []loop.Option
	loop    *loop.Loop
	err     error
}

// Go starts an Actor with the passed options.
func Go(options ...Option) (*Actor, error) {
	// Init with options.
	act := &Actor{}
	for _, option := range options {
		if err := option(act); err != nil {
			// One of the options made troubles.
			act.err = failure.First(act.err, err)
			return nil, act.err
		}
	}
	// Ensure default settings.
	if act.actionC == nil {
		act.actionC = make(chan Action, 1)
	}
	// Create loop with its options.
	l, err := loop.Go(act.worker, act.options...)
	if err != nil {
		act.err = failure.First(act.err, err)
		return nil, act.err
	}
	act.loop = l
	return act, nil
}

// DoSync executes the actor function and returns when it's done
// or it has the default timeout.
func (act *Actor) DoSync(action Action) error {
	return act.DoSyncTimeout(action, DefaultTimeout)
}

// DoSyncTimeout executes the action and returns when it's done
// or it has a timeout.
func (act *Actor) DoSyncTimeout(action Action, timeout time.Duration) error {
	waitC := make(chan struct{})
	if err := act.DoAsync(func() error {
		err := action()
		close(waitC)
		return err
	}); err != nil {
		return err
	}
	select {
	case <-waitC:
	case <-time.After(timeout):
		return failure.New("synchronous actor do: timed out")
	}
	return nil
}

// DoAsync executes the actor function and returns immediately
func (act *Actor) DoAsync(action Action) error {
	act.mu.Lock()
	defer act.mu.Unlock()
	if act.err != nil {
		return act.err
	}
	act.actionC <- action
	return nil
}

// Stop terminates the Actor with the passed error. That or
// a potential earlier error will be returned.
func (act *Actor) Stop(err error) error {
	if act.loop != nil {
		return act.loop.Stop(err)
	}
	return act.err
}

// Signaler allows getting information about the status of the actor.
func (act *Actor) Signaler() notifier.Signaler {
	act.mu.RLock()
	defer act.mu.RUnlock()
	return act.loop.Signaler()
}

// Err returns information if the Actor has an error.
func (act *Actor) Err() error {
	act.mu.Lock()
	defer act.mu.Unlock()
	if act.err != nil {
		return act.err
	}
	return act.loop.Err()
}

// QueueCap returns the capacity of the action queue.
func (act *Actor) QueueCap() int {
	act.mu.Lock()
	defer act.mu.Unlock()
	return cap(act.actionC)
}

// worker is the Loop worker of the Actor.
func (act *Actor) worker(closer *notifier.Closer) error {
	for {
		select {
		case <-closer.Done():
			return nil
		case action := <-act.actionC:
			if err := action(); err != nil {
				return err
			}
		}
	}
}

// EOF
