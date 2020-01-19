// Tideland Go Together - Loop
//
// Copyright (C) 2017-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package loop // import "tideland.dev/go/together/loop"

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"fmt"
	"time"

	"tideland.dev/go/together/fuse"
)

//--------------------
// CONSTANTS
//--------------------

// timeout defines the time to wait for signals.
const timeout = 5 * time.Second

//--------------------
// FUNCTION TYPES
//--------------------

// Terminator describes a type signaling that the work is done.
type Terminator interface {
	Done() <-chan struct{}
}

// Worker discribes the function running the loop.
type Worker func(t Terminator) error

// Recoverer allows a goroutine to react on a panic during its
// work. If it returns nil the goroutine shall continue
// work. Otherwise it will return with an error the gouroutine
// may use for its continued processing.
type Recoverer func(reason interface{}) error

// Finalizer is called with the actors internal status when
// the backend loop terminates.
type Finalizer func(err error) error

//--------------------
// LOOP
//--------------------

// Loop manages running for-select-loops in the background as goroutines
// in a controlled way. Users can get information about status and possible
// failure as well as control how to stop, restart, or recover via
// options.
type Loop struct {
	ctx       context.Context
	cancel    func()
	signal    *fuse.Signal
	worker    Worker
	recoverer Recoverer
	finalizer Finalizer
	err       fuse.Error
}

// Go starts a loop running the given worker with the
// given options.
func Go(worker Worker, options ...Option) (*Loop, error) {
	// Init with default values.
	l := &Loop{
		signal: fuse.NewSignal(),
		worker: worker,
	}
	for _, option := range options {
		if err := option(l); err != nil {
			return nil, err
		}
	}
	// Ensure default settings.
	if l.ctx == nil {
		l.ctx, l.cancel = context.WithCancel(context.Background())
	} else {
		l.ctx, l.cancel = context.WithCancel(l.ctx)
	}
	if l.recoverer == nil {
		l.recoverer = func(reason interface{}) error {
			return fmt.Errorf("loop panic: %v", reason)
		}
	}
	if l.finalizer == nil {
		l.finalizer = func(err error) error {
			return err
		}
	}
	// Create loop with its options.
	l.signal.Notify(fuse.Starting)
	go l.backend()
	if err := l.signal.Wait(fuse.Ready, timeout); err != nil {
		return nil, err
	}
	return l, nil
}

// Err returns information if the Loop has an error.
func (l *Loop) Err() error {
	return l.err.Get()
}

// Stop terminates the Loop backend.
func (l *Loop) Stop() error {
	if !l.err.IsNil() {
		return l.err.Get()
	}
	l.cancel()
	if err := l.signal.Wait(fuse.Stopped, timeout); err != nil {
		return err
	}
	return l.err.Get()
}

// backend runs the loop worker as goroutine as long as
// the status is notifier.Working.
func (l *Loop) backend() {
	defer func() {
		l.err.Set(l.finalizer(l.err.Get()))
		l.signal.Notify(fuse.Stopped)
	}()
	l.signal.Notify(fuse.Ready)
	for l.wrapper() {
	}
}

// wrapper wraps the wrapper, handles possible failure, and
// manages panics.
func (l *Loop) wrapper() (ok bool) {
	defer func() {
		if reason := recover(); reason != nil {
			// Panic!
			err := l.recoverer(reason)
			if err != nil {
				l.err.Set(err)
				ok = false
			} else {
				ok = true
			}
		} else {
			// Regular ending.
			ok = false
		}
		if !ok {
			l.signal.Notify(fuse.Stopping)
		}
	}()
	l.err.Set(l.worker(l.ctx))
	return
}

// EOF
