// Tideland Go Together - Loop - Unit Tests
//
// Copyright (C) 2017-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package loop_test

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/loop"
)

//--------------------
// TESTS
//--------------------

// TestPureOK tests a loop without any options, stopping without an error.
func TestPureOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	started := make(chan struct{})
	stopped := make(chan struct{})
	beenThereDoneThat := false
	worker := func(ctx context.Context) error {
		defer close(stopped)
		close(started)
		for {
			select {
			case <-ctx.Done():
				beenThereDoneThat = true
				return nil
			case <-time.Tick(time.Minute):
				// Just for linter.
			}
		}
	}
	l, err := loop.Go(worker)
	assert.NoError(err)

	// Test.
	assert.NoError(l.Err())
	<-started
	l.Stop()
	<-stopped
	assert.NoError(l.Err())
	assert.True(beenThereDoneThat)
}

// TestPureError tests a loop without any options, stopping with an error.
func TestPureError(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	started := make(chan struct{})
	stopped := make(chan struct{})
	worker := func(ctx context.Context) error {
		defer close(stopped)
		close(started)
		for {
			select {
			case <-ctx.Done():
				return errors.New("ouch")
			case <-time.Tick(50 * time.Millisecond):
				// Just for linter.
			}
		}
	}
	l, err := loop.Go(worker)
	assert.NoError(err)

	// Test.
	assert.NoError(l.Err())
	<-started
	l.Stop()
	<-stopped
	assert.ErrorMatch(l.Err(), "ouch")
}

// TestContextCancelOK tests the stopping after a context cancel w/o error.
func TestContextCancelOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	ctx, cancel := context.WithCancel(context.Background())
	stopped := make(chan struct{})
	worker := func(ctx context.Context) error {
		defer close(stopped)
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-time.Tick(time.Minute):
				// Just for linter.
			}
		}
	}
	l, err := loop.Go(
		worker,
		loop.WithContext(ctx),
	)
	assert.NoError(err)

	// Test.
	cancel()
	<-stopped
	assert.NoError(l.Err())
}

// TestContextCancelError tests the stopping after a context cancel w/ error.
func TestContextCancelError(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	ctx, cancel := context.WithCancel(context.Background())
	stopped := make(chan struct{})
	worker := func(ctx context.Context) error {
		defer close(stopped)
		for {
			select {
			case <-ctx.Done():
				return errors.New("oh, no")
			case <-time.Tick(time.Minute):
				// Just for linter.
			}
		}
	}
	l, err := loop.Go(
		worker,
		loop.WithContext(ctx),
	)
	assert.NoError(err)

	// Test.
	cancel()
	<-stopped
	assert.ErrorMatch(l.Err(), "oh, no")
}

// TestFinalizerOK tests calling a finalizer returning an own error.
func TestFinalizerOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	stopped := make(chan struct{})
	finalized := false
	worker := func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-time.Tick(time.Minute):
				// Just for linter.
			}
		}
	}
	finalizer := func(err error) error {
		defer close(stopped)
		assert.NoError(err)
		finalized = true
		return err
	}
	l, err := loop.Go(
		worker,
		loop.WithFinalizer(finalizer),
	)
	assert.NoError(err)

	// Test.
	l.Stop()
	<-stopped
	assert.NoError(l.Err())
	assert.True(finalized)
}

// TestFinalizerError tests the stopping with an error but
// finalizer returns an own one.
func TestFinalizerError(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	stopped := make(chan struct{})
	worker := func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return errors.New("don't want to stop")
			case <-time.Tick(time.Minute):
				// Just for linter.
			}
		}
	}
	finalizer := func(err error) error {
		defer close(stopped)
		assert.ErrorMatch(err, "don't want to stop")
		return errors.New("don't care")
	}
	l, err := loop.Go(
		worker,
		loop.WithFinalizer(finalizer),
	)
	assert.NoError(err)

	// Test.
	l.Stop()
	<-stopped
	assert.ErrorMatch(l.Err(), "don't care")
}

// TestInternalError tests the stopping after an internal error.
func TestInternalError(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	stopped := make(chan struct{})
	worker := func(ctx context.Context) error {
		defer close(stopped)
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(50 * time.Millisecond):
				return errors.New("time over")
			}
		}
	}
	l, err := loop.Go(worker)
	assert.NoError(err)

	// Test.
	<-stopped
	l.Stop()
	assert.ErrorMatch(l.Err(), "time over")
}

// TestRepairerOK tests the stopping without an error if Loop has a repairer.
// Repairer must never been called.
func TestRepairerOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	repaired := make(chan struct{})
	worker := func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-time.Tick(time.Minute):
				// Just for linter.
			}
		}
	}
	repairer := func(reason interface{}) error {
		defer close(repaired)
		return nil
	}
	l, err := loop.Go(
		worker,
		loop.WithRepairer(repairer),
	)
	assert.NoError(err)

	// Test.
	l.Stop()

	select {
	case <-repaired:
		assert.Fail("repairer called")
	case <-time.After(100 * time.Millisecond):
		assert.OK(true)
	}
}

// TestRepairerErrorOK tests the stopping with an error if Loop has a repairer.
// Repairer must never been called.
func TestRepairerErrorOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	stopped := make(chan struct{})
	worker := func(ctx context.Context) error {
		defer close(stopped)
		for {
			select {
			case <-ctx.Done():
				return errors.New("oh, no")
			case <-time.Tick(time.Minute):
				// Just for linter.
			}
		}
	}
	repairer := func(reason interface{}) error {
		return fmt.Errorf("unexpected")
	}
	l, err := loop.Go(
		worker,
		loop.WithRepairer(repairer),
	)
	assert.NoError(err)

	// Test.
	l.Stop()
	<-stopped
	assert.ErrorMatch(l.Err(), "oh, no")
}

// TestRecoverPanics tests the stopping handling and later stopping
// after panics.
func TestRecoverPanicsOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	panics := 0
	stopped := make(chan struct{})
	worker := func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(50 * time.Millisecond):
				panic("bam")
			}
		}
	}
	finalizer := func(err error) error {
		defer close(stopped)
		return err
	}
	repairer := func(reason interface{}) error {
		panics++
		if panics > 10 {
			return fmt.Errorf("too many panics: %v", reason)
		}
		return nil
	}
	l, err := loop.Go(
		worker,
		loop.WithFinalizer(finalizer),
		loop.WithRepairer(repairer),
	)
	assert.NoError(err)

	// Test.
	<-stopped
	assert.ErrorMatch(l.Err(), "too many panics: bam")
}

//--------------------
// EXAMPLES
//--------------------

// ExampleWorker shows the usage of Loop with no repairer. The inner loop
// contains a select listening to the channel returned by Closer.Done().
// Other channels are for the standard communication with the Loop.
func ExampleWorker() {
	prints := make(chan string)
	ctx, cancel := context.WithCancel(context.Background())
	// Sample loop worker.
	worker := func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				// We shall stop.
				return nil
			case str := <-prints:
				// Standard work of example loop.
				if str == "panic" {
					return errors.New("panic")
				}
				println(str)
			}
		}
	}
	l, err := loop.Go(worker, loop.WithContext(ctx))
	if err != nil {
		panic(err)
	}

	prints <- "Hello"
	prints <- "World"

	// cancel() terminates the loop via the context.
	cancel()

	// Returned error must be nil in this example.
	if l.Err() != nil {
		panic(l.Err())
	}
}

// ExampleRepairer demonstrates the usage of a repairer.
func ExampleRepairer() {
	panics := make(chan string)
	// Sample loop worker.
	worker := func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case str := <-panics:
				panic(str)
			}
		}
	}
	// Repairer function checks the reasion. "never mind" will
	// be repaired, all others lead to an error. The repairer
	// is also responsable for fixing the owners state crashed
	// during panic.
	repairer := func(reason interface{}) error {
		why := reason.(string)
		if why == "never mind" {
			return nil
		}
		return fmt.Errorf("worker panic: %v", why)
	}
	l, err := loop.Go(worker, loop.WithRepairer(repairer))
	if err != nil {
		panic(err)
	}
	l.Stop()
}

// EOF
