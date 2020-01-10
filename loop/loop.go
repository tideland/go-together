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
	"sync"
	"time"

	"tideland.dev/go/together/notifier"
	"tideland.dev/go/trace/failure"
)

//--------------------
// RECOVERER
//--------------------

// Recoverer allows a goroutine to react on a panic during its
// work. If it returns nil the goroutine shall continue
// work. Otherwise it will return with an error the gouroutine
// may use for its continued processing.
type Recoverer func(reason interface{}) error

// DefaultRecoverer simply re-panics.
func DefaultRecoverer(reason interface{}) error {
	panic(reason)
}

//--------------------
// LOOP
//--------------------

const (
	// timeout for loop action waitings.
	timeout = 5 * time.Second
)

// Worker is a managed Loop function performing the work.
type Worker func(clsr *notifier.Closer) error

// Finalizer allows to perform some steps to clean-up when
// the worker terminates. The passed error is the state of
// the loop.
type Finalizer func(err error) error

// Loop manages running for-select-loops in the background as goroutines
// in a controlled way. Users can get information about status and possible
// failure as well as control how to stop, restart, or recover via
// options.
type Loop struct {
	mu        sync.RWMutex
	worker    Worker
	finalizer Finalizer
	closer    *notifier.Closer
	bundle    *notifier.Bundle
	recoverer Recoverer
	err       error
}

// Go starts a loop running the given worker with the
// given options.
func Go(worker Worker, options ...Option) (*Loop, error) {
	// Init with default values.
	loop := &Loop{
		worker:    worker,
		closer:    notifier.NewCloser(),
		bundle:    notifier.NewBundle(),
		recoverer: DefaultRecoverer,
	}
	// Apply options.
	for _, option := range options {
		if err := option(loop); err != nil {
			// One of the options made troubles.
			loop.err = failure.First(loop.err, err)
			loop.bundle.Notify(notifier.Stopped)
			return nil, loop.err
		}
	}
	// Start goroutine and wait until it's working.
	go loop.backend()
	loop.err = failure.First(loop.err, loop.bundle.Wait(notifier.Working, timeout))
	if loop.err != nil {
		loop.closer.Close()
		return nil, loop.err
	}
	return loop, nil
}

// Stop terminates the Loop with the passed error. That or
// a potential earlier error will be returned.
func (loop *Loop) Stop(err error) error {
	loop.mu.Lock()
	defer loop.mu.Unlock()
	loop.closer.Close()
	if werr := loop.bundle.Wait(notifier.Stopped, timeout); werr != nil {
		loop.err = failure.First(loop.err, werr)
		return loop.err
	}
	loop.err = failure.First(loop.err, err)
	return loop.err
}

// Signaler allows getting information about the status of the loop.
func (loop *Loop) Signaler() notifier.Signaler {
	loop.mu.RLock()
	defer loop.mu.RUnlock()
	return loop.bundle
}

// Err returns information if the Loop has an error.
func (loop *Loop) Err() error {
	loop.mu.RLock()
	defer loop.mu.RUnlock()
	return loop.err
}

// backend runs the loop worker as goroutine as long as
// the status is notifier.Working.
func (loop *Loop) backend() {
	defer loop.bundle.Notify(notifier.Stopped)
	loop.bundle.Notify(notifier.Working)
	for loop.bundle.Status() == notifier.Working {
		loop.container()
	}
	if loop.finalizer != nil {
		loop.err = failure.First(loop.err, loop.finalizer(loop.err))
	}
}

// container wraps the worker, handles possible failure, and
// manages panics.
func (loop *Loop) container() {
	defer func() {
		if reason := recover(); reason != nil {
			// Panic, try to recover.
			if err := loop.recoverer(reason); err != nil {
				loop.err = failure.First(loop.err, err)
				loop.bundle.Notify(notifier.Stopping)
			}
		} else {
			// Regular ending.
			loop.bundle.Notify(notifier.Stopping)
		}
	}()
	loop.err = failure.First(loop.err, loop.worker(loop.closer))
}

// EOF
