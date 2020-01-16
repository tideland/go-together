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
	beenThereDoneThat := false
	worker := func(lt loop.Terminator) error {
		close(started)
		for {
			select {
			case <-lt.Done():
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
	assert.NoError(l.Stop())
	assert.True(beenThereDoneThat)
}

// TestPureError tests a loop without any options, stopping with an error.
func TestPureError(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	started := make(chan struct{})
	worker := func(lt loop.Terminator) error {
		close(started)
		for {
			select {
			case <-lt.Done():
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
	assert.ErrorMatch(l.Stop(), "ouch")
	assert.ErrorMatch(l.Err(), "ouch")
}

// TestContextCancelOK tests the stopping after a context cancel w/o error.
func TestContextCancelOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	worker := func(lt loop.Terminator) error {
		defer close(done)
		for {
			select {
			case <-lt.Done():
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
	<-done
	assert.NoError(l.Err())
}

// TestContextCancelError tests the stopping after a context cancel w/ error.
func TestContextCancelError(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	worker := func(lt loop.Terminator) error {
		defer close(done)
		for {
			select {
			case <-lt.Done():
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
	<-done
	assert.ErrorMatch(l.Err(), "oh, no")
}

// TestFinalizerOK tests calling a finalizer returning an own error.
func TestFinalizerOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	finalized := false
	worker := func(lt loop.Terminator) error {
		for {
			select {
			case <-lt.Done():
				return nil
			case <-time.Tick(time.Minute):
				// Just for linter.
			}
		}
	}
	finalizer := func(err error) error {
		assert.NoError(err)
		finalized = true
		return errors.New("finalization error")
	}
	l, err := loop.Go(
		worker,
		loop.WithFinalizer(finalizer),
	)
	assert.NoError(err)

	// Test.
	assert.ErrorMatch(l.Stop(), "finalization error")
	assert.ErrorMatch(l.Err(), "finalization error")
	assert.True(finalized)
}

// TestFinalizerError tests the stopping with an error, is kept
// even if finalizer returns an error.
func TestFinalizerError(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	worker := func(lt loop.Terminator) error {
		for {
			select {
			case <-lt.Done():
				return errors.New("don't want to stop")
			case <-time.Tick(time.Minute):
				// Just for linter.
			}
		}
	}
	finalizer := func(err error) error {
		assert.ErrorMatch(err, "don't want to stop")
		return errors.New("don't care")
	}
	l, err := loop.Go(
		worker,
		loop.WithFinalizer(finalizer),
	)
	assert.NoError(err)

	// Test.
	assert.ErrorMatch(l.Stop(), "don't care")
	assert.ErrorMatch(l.Err(), "don't care")
}

// TestInternalError tests the stopping after an internal error.
func TestInternalError(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	worker := func(lt loop.Terminator) error {
		for {
			select {
			case <-lt.Done():
				return nil
			case <-time.After(50 * time.Millisecond):
				return errors.New("time over")
			}
		}
	}
	l, err := loop.Go(worker)
	assert.NoError(err)

	// Test.
	time.Sleep(100 * time.Millisecond)
	assert.ErrorMatch(l.Stop(), "time over")
	assert.ErrorMatch(l.Err(), "time over")
}

// TestRecoveredOK tests the stopping without an error if Loop has a recoverer.
// Recoverer must never been called.
func TestRecoveredOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	beenThereDoneThat := false
	worker := func(lt loop.Terminator) error {
		for {
			select {
			case <-lt.Done():
				return nil
			case <-time.Tick(time.Minute):
				// Just for linter.
			}
		}
	}
	recoverer := func(reason interface{}) error {
		beenThereDoneThat = true
		return nil
	}
	l, err := loop.Go(
		worker,
		loop.WithRecoverer(recoverer),
	)
	assert.NoError(err)

	// Test.
	assert.NoError(l.Stop())
	assert.NoError(l.Err())
	assert.False(beenThereDoneThat)
}

// TestRecovererErrorOK tests the stopping with an error if Loop has a recoverer.
// Recoverer must never been called.
func TestRecovererErrorOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	beenThereDoneThat := false
	worker := func(lt loop.Terminator) error {
		for {
			select {
			case <-lt.Done():
				return errors.New("oh, no")
			case <-time.Tick(time.Minute):
				// Just for linter.
			}
		}
	}
	recoverer := func(reason interface{}) error {
		beenThereDoneThat = true
		return nil
	}
	l, err := loop.Go(
		worker,
		loop.WithRecoverer(recoverer),
	)
	assert.NoError(err)

	// Test.
	assert.ErrorMatch(l.Stop(), "oh, no")
	assert.ErrorMatch(l.Err(), "oh, no")
	assert.False(beenThereDoneThat)
}

// TestRecoverPanics tests the stopping handling and later stopping
// after panics.
func TestRecoverPanicsOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	panics := 0
	done := make(chan struct{})
	worker := func(lt loop.Terminator) error {
		for {
			select {
			case <-lt.Done():
				return nil
			case <-time.After(50 * time.Millisecond):
				panic("bam")
			}
		}
	}
	finalizer := func(err error) error {
		defer close(done)
		return err
	}
	recoverer := func(reason interface{}) error {
		panics++
		if panics > 10 {
			return fmt.Errorf("too much: %v", reason)
		}
		return nil
	}
	l, err := loop.Go(
		worker,
		loop.WithFinalizer(finalizer),
		loop.WithRecoverer(recoverer),
	)
	assert.NoError(err)

	// Test.
	<-done
	assert.ErrorMatch(l.Err(), "too much: bam")
}

//--------------------
// EXAMPLES
//--------------------

// ExampleWorker shows the usage of Loop with no recoverer. The inner loop
// contains a select listening to the channel returned by Closer.Done().
// Other channels are for the standard communication with the Loop.
func ExampleWorker() {
	prints := make(chan string)
	ctx, cancel := context.WithCancel(context.Background())
	// Sample loop worker.
	worker := func(lt loop.Terminator) error {
		for {
			select {
			case <-lt.Done():
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

	cancel()

	if l.Err() != nil {
		panic(l.Err())
	}
}

// ExampleRecoverer demonstrates the usage of a recoverer.
// Here the frequency of the recovered reasons (more than five
// in 10 milliseconds) or the total number is checked. If the
// total number is not interesting the reasons could be
// trimmed by e.g. rs.Trim(5). The fields Time and Reason per
// recovering allow even more diagnosis.
func ExampleRecoverer() {
	panics := make(chan string)
	// Sample loop worker.
	worker := func(lt loop.Terminator) error {
		for {
			select {
			case <-lt.Done():
				return nil
			case str := <-panics:
				panic(str)
			}
		}
	}
	// Recovery function checking frequency and total number.
	count := 0
	recoverer := func(reason interface{}) error {
		count++
		if count > 10 {
			return errors.New("too many errors")
		}
		return nil
	}
	loop.Go(worker, loop.WithRecoverer(recoverer))
}

// EOF
