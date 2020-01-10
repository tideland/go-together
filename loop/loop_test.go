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
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/loop"
	"tideland.dev/go/together/notifier"
)

//--------------------
// TESTS
//--------------------

// TestPure tests a loop without any options, stopping without an error.
func TestPure(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	waitC := make(chan struct{})
	beenThereDoneThat := false
	worker := func(c *notifier.Closer) error {
		close(waitC)
		for {
			select {
			case <-c.Done():
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
	<-waitC
	assert.Equal(l.Signaler().Status(), notifier.Working)
	assert.NoError(l.Stop(nil))
	assert.Equal(l.Signaler().Status(), notifier.Stopped)
	assert.True(beenThereDoneThat)
}

// TestPureError tests a loop without any options, stopping with an error.
func TestPureError(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	waitC := make(chan struct{})
	worker := func(c *notifier.Closer) error {
		close(waitC)
		for {
			select {
			case <-c.Done():
				return nil
			case <-time.Tick(time.Minute):
				// Just for linter.
			}
		}
	}
	l, err := loop.Go(worker)
	assert.NoError(err)

	// Test.
	<-waitC
	err = l.Stop(errors.New("ouch"))
	assert.ErrorMatch(err, "ouch")
}

// TestPureInternalError tests a loop without any options, stopping leads
// to an internal error.
func TestPureInternalError(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	waitC := make(chan struct{})
	worker := func(c *notifier.Closer) error {
		close(waitC)
		for {
			select {
			case <-c.Done():
				return errors.New("ouch")
			case <-time.Tick(time.Minute):
				// Just for linter.
			}
		}
	}
	l, err := loop.Go(worker)
	assert.NoError(err)

	// Test.
	<-waitC
	err = l.Stop(nil)
	assert.ErrorMatch(err, "ouch")
}

// TestContextCancelOK tests the stopping after a context cancel w/o error.
func TestContextCancelOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	ctx, cancel := context.WithCancel(context.Background())
	worker := func(c *notifier.Closer) error {
		for {
			select {
			case <-c.Done():
				return nil
			case <-time.Tick(time.Minute):
				// Just for linter.
			}
		}
	}
	signalbox := notifier.NewSignalbox()
	l, err := loop.Go(
		worker,
		loop.WithContext(ctx),
		loop.WithSignalbox(signalbox),
	)
	assert.NoError(err)

	// Test.
	<-signalbox.Done(notifier.Working)
	cancel()
	<-signalbox.Done(notifier.Stopped)
	assert.NoError(l.Err())
}

// TestContextCancelError tests the stopping after a context cancel w/ error.
func TestContextCancelError(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	ctx, cancel := context.WithCancel(context.Background())
	worker := func(c *notifier.Closer) error {
		for {
			select {
			case <-c.Done():
				return errors.New("oh, no")
			case <-time.Tick(time.Minute):
				// Just for linter.
			}
		}
	}
	signalbox := notifier.NewSignalbox()
	l, err := loop.Go(
		worker,
		loop.WithContext(ctx),
		loop.WithSignalbox(signalbox),
	)
	assert.NoError(err)

	// Test.
	<-signalbox.Done(notifier.Working)
	cancel()
	<-signalbox.Done(notifier.Stopped)
	assert.ErrorMatch(l.Err(), "oh, no")
}

// TestMultipleNotifier tests the usage of multiple notifiers.
func TestMultipleNotifier(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	worker := func(c *notifier.Closer) error {
		for {
			select {
			case <-c.Done():
				return nil
			case <-time.Tick(time.Minute):
				// Just for linter.
			}
		}
	}
	signalboxA := notifier.NewSignalbox()
	signalboxB := notifier.NewSignalbox()
	signalboxC := notifier.NewSignalbox()
	l, err := loop.Go(
		worker,
		loop.WithSignalbox(signalboxA),
		loop.WithSignalbox(signalboxB),
		loop.WithSignalbox(signalboxC),
	)
	assert.NoError(err)

	// Test.
	<-signalboxC.Done(notifier.Working)
	l.Stop(nil)

	x := 0

timeout:
	for x != 7 {
		select {
		case <-signalboxA.Done(notifier.Stopped):
			x |= 1
		case <-signalboxB.Done(notifier.Stopped):
			x |= 2
		case <-signalboxC.Done(notifier.Stopped):
			x |= 4
		case <-time.After(time.Second):
			break timeout
		}
	}

	assert.Equal(x, 7)
	assert.NoError(l.Err())
}

// TestFinalizerOK tests calling a finalizer returning an own error.
func TestFinalizerOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	finalized := false
	worker := func(c *notifier.Closer) error {
		for {
			select {
			case <-c.Done():
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
	signalbox := notifier.NewSignalbox()
	l, err := loop.Go(
		worker,
		loop.WithFinalizer(finalizer),
		loop.WithSignalbox(signalbox),
	)
	assert.NoError(err)

	// Test.
	<-signalbox.Done(notifier.Working)
	l.Stop(nil)
	<-signalbox.Done(notifier.Stopped)
	assert.ErrorMatch(l.Err(), "finalization error")
	assert.True(finalized)
}

// TestFinalizerError tests the stopping with an error, is kept
// even if finalizer returns an error.
func TestFinalizerError(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	worker := func(c *notifier.Closer) error {
		for {
			select {
			case <-c.Done():
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
	signalbox := notifier.NewSignalbox()
	l, err := loop.Go(
		worker,
		loop.WithFinalizer(finalizer),
		loop.WithSignalbox(signalbox),
	)
	assert.NoError(err)

	// Test.
	<-signalbox.Done(notifier.Working)
	l.Stop(nil)
	<-signalbox.Done(notifier.Stopped)
	assert.ErrorMatch(l.Err(), "don't want to stop")
}

// TestInternalOK tests the stopping w/o an error.
func TestInternalOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	worker := func(c *notifier.Closer) error {
		for {
			select {
			case <-c.Done():
				return nil
			case <-time.After(50 * time.Millisecond):
				return nil
			}
		}
	}
	signalbox := notifier.NewSignalbox()
	l, err := loop.Go(
		worker,
		loop.WithSignalbox(signalbox),
	)
	assert.NoError(err)

	// Test.
	<-signalbox.Done(notifier.Stopped)
	assert.NoError(l.Err())
}

// TestInternalError tests the stopping after an internal error.
func TestInternalError(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	worker := func(c *notifier.Closer) error {
		for {
			select {
			case <-c.Done():
				return nil
			case <-time.After(50 * time.Millisecond):
				return errors.New("time over")
			}
		}
	}
	signalbox := notifier.NewSignalbox()
	l, err := loop.Go(
		worker,
		loop.WithSignalbox(signalbox),
	)
	assert.NoError(err)

	// Test.
	<-signalbox.Done(notifier.Stopped)
	assert.ErrorMatch(l.Stop(nil), "time over")
	assert.ErrorMatch(l.Err(), "time over")
}

// TestRecoveredOK tests the stopping without an error if Loop has a recoverer.
// Recoverer must never been called.
func TestRecoveredOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	beenThereDoneThat := false
	worker := func(c *notifier.Closer) error {
		for {
			select {
			case <-c.Done():
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
	signalbox := notifier.NewSignalbox()
	l, err := loop.Go(
		worker,
		loop.WithRecoverer(recoverer),
		loop.WithSignalbox(signalbox),
	)
	assert.NoError(err)

	// Test.
	<-signalbox.Done(notifier.Working)
	l.Stop(nil)
	<-signalbox.Done(notifier.Stopped)
	assert.Nil(l.Err())
	assert.False(beenThereDoneThat)
}

// TestRecovererError tests the stopping with an error if Loop has a recoverer.
// Recoverer must never been called.
func TestRecovererError(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	beenThereDoneThat := false
	worker := func(c *notifier.Closer) error {
		for {
			select {
			case <-c.Done():
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
	signalbox := notifier.NewSignalbox()
	l, err := loop.Go(
		worker,
		loop.WithRecoverer(recoverer),
		loop.WithSignalbox(signalbox),
	)
	assert.NoError(err)

	// Test.
	<-signalbox.Done(notifier.Working)
	l.Stop(nil)
	<-signalbox.Done(notifier.Stopped)
	assert.ErrorMatch(l.Err(), "oh, no")
	assert.False(beenThereDoneThat)
}

// TestRecoverPanicsOK tests the stopping w/o an error.
func TestRecoverPanicsOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	panics := 0
	doneC := make(chan struct{})
	worker := func(c *notifier.Closer) error {
		for {
			select {
			case <-c.Done():
				return nil
			case <-doneC:
				return nil
			case <-time.After(10 * time.Millisecond):
				panic("bam")
			}
		}
	}
	recoverer := func(reason interface{}) error {
		panics++
		if panics > 10 {
			close(doneC)
		}
		return nil
	}
	signalbox := notifier.NewSignalbox()
	l, err := loop.Go(
		worker,
		loop.WithRecoverer(recoverer),
		loop.WithSignalbox(signalbox),
	)
	assert.NoError(err)

	// Test.
	<-signalbox.Done(notifier.Stopped)
	assert.NoError(l.Err())
}

// TestRecoverPanicsError tests the stopping w/o an error.
func TestRecoverPanicsError(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	panics := 0
	worker := func(c *notifier.Closer) error {
		for {
			select {
			case <-c.Done():
				return nil
			case <-time.After(10 * time.Millisecond):
				panic("bam")
			}
		}
	}
	recoverer := func(reason interface{}) error {
		panics++
		if panics > 10 {
			return errors.New("superbam")
		}
		return nil
	}
	signalbox := notifier.NewSignalbox()
	l, err := loop.Go(
		worker,
		loop.WithRecoverer(recoverer),
		loop.WithSignalbox(signalbox),
	)
	assert.NoError(err)

	// Test.
	<-signalbox.Done(notifier.Stopped)
	assert.ErrorMatch(l.Err(), "superbam")
}

// TestReasons tests collecting and analysing loop recovery reasons.
func TestReasons(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	rins := []error{
		errors.New("error a"),
		errors.New("error b"),
		errors.New("error c"),
		errors.New("error d"),
		errors.New("error e"),
	}
	rs := loop.MakeReasons()
	for _, rin := range rins {
		time.Sleep(100 * time.Millisecond)
		rs = rs.Append(rin)
	}

	// Test.
	assert.Length(rs, 5)
	assert.Equal(rs.Last().Reason, rins[4])
	assert.True(rs.Frequency(5, time.Second))
	assert.False(rs.Frequency(5, 10*time.Millisecond))

	rs = rs.Trim(3)

	assert.Length(rs, 3)
	assert.Match(rs.String(), `\[\['error c' @ .*\] / \['error d' @ .*\] / \['error e' @ .*\]\]`)
}

//--------------------
// EXAMPLES
//--------------------

// ExampleWorker shows the usage of Loo with no recoverer. The inner loop
// contains a select listening to the channel returned by Closer.Done().
// Other channels are for the standard communication with the Loop.
func ExampleWorker() {
	printC := make(chan string)
	ctx, cancel := context.WithCancel(context.Background())
	// Sample loop worker.
	worker := func(c *notifier.Closer) error {
		for {
			select {
			case <-c.Done():
				// We shall stop.
				return nil
			case str := <-printC:
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

	printC <- "Hello"
	printC <- "World"

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
	panicC := make(chan string)
	// Sample loop worker.
	worker := func(c *notifier.Closer) error {
		for {
			select {
			case <-c.Done():
				return nil
			case str := <-panicC:
				panic(str)
			}
		}
	}
	// Recovery function checking frequency and total number.
	rs := loop.MakeReasons()
	recoverer := func(reason interface{}) error {
		rs = rs.Append(reason)
		if rs.Frequency(5, 10*time.Millisecond) {
			return errors.New("too high error frequency")
		}
		if rs.Len() >= 10 {
			return errors.New("too many errors")
		}
		return nil
	}
	loop.Go(worker, loop.WithRecoverer(recoverer))
}

// EOF
