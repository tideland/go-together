// Tideland Go Together - Loop - Unit Tests
//
// Copyright (C) 2017-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package notifier_test

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"sync"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/notifier"
)

//--------------------
// CONSTANTS
//--------------------

// timeout is the waitng time..
var timeout time.Duration = 5 * time.Second

//--------------------
// TESTS
//--------------------

// TestCloserOK tests the closing of the Closer when one or more of the
// input closer channels closes.
func TestCloserOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	ccs := []chan struct{}{
		make(chan struct{}),
		make(chan struct{}),
		make(chan struct{}),
		make(chan struct{}),
		make(chan struct{}),
	}
	closer := notifier.NewCloser(ccs[0], ccs[1], ccs[2], ccs[3], ccs[4])
	beenThereDoneThat := false

	// Test.
	go func() {
		time.Sleep(100 * time.Millisecond)
		close(ccs[1])
		time.Sleep(100 * time.Millisecond)
		close(ccs[3])
	}()

	select {
	case <-closer.Done():
		beenThereDoneThat = true
	case <-time.After(timeout):
		assert.Fail("timeout")
	}
	assert.True(beenThereDoneThat)
}

// TestCloserAddOK tests the closing the Closer by an added channel.
func TestCloserAddOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	cca := make(chan struct{})
	ccb := make(chan struct{})
	closer := notifier.NewCloser(cca)
	beenThereDoneThat := false

	closer.Add(ccb)

	// Test.
	go func() {
		time.Sleep(100 * time.Millisecond)
		close(ccb)
	}()

	select {
	case <-closer.Done():
		beenThereDoneThat = true
	case <-time.After(timeout):
		assert.Fail("timeout")
	}
	assert.True(beenThereDoneThat)
}

// TestCloserContext tests the closing of the Closer with a closing context.
func TestCloserContext(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	ctx, cancel := context.WithCancel(context.Background())
	closer := notifier.NewCloser(ctx.Done())
	beenThereDoneThat := false

	// Test.
	go func() {
		time.Sleep(time.Second)
		cancel()
	}()

	select {
	case <-closer.Done():
		beenThereDoneThat = true
	case <-time.After(timeout):
		assert.Fail("timeout")
	}
	assert.True(beenThereDoneThat)
}

// TestCloserClose tests the closing the Closer directly.
func TestCloserClose(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	closer := notifier.NewCloser()
	beenThereDoneThat := false

	// Test.
	go func() {
		time.Sleep(time.Second)
		closer.Close()
	}()

	select {
	case <-closer.Done():
		beenThereDoneThat = true
	case <-time.After(timeout):
		assert.Fail("timeout")
	}
	assert.True(beenThereDoneThat)
}

// TestCloserTimeout tests what happes if no channel signals the closing.
func TestCloserTimeout(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	closer := notifier.NewCloser()
	beenThereDoneThat := false

	// Test.
	select {
	case <-closer.Done():
		assert.Fail("invalid closing")
	case <-time.After(timeout):
		beenThereDoneThat = true
	}
	assert.True(beenThereDoneThat)
}

// TestSignalboxWaitOK tests the correct notification and waiting.
func TestSignalboxWaitOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	timeout := 10 * time.Millisecond
	signalbox := notifier.NewSignalbox()
	var wg sync.WaitGroup
	wg.Add(2)
	asserter := func() {
		defer wg.Done()
		assert.NoError(signalbox.Wait(notifier.Starting, timeout))
		assert.NoError(signalbox.Wait(notifier.Starting, timeout))
		assert.NoError(signalbox.Wait(notifier.Ready, timeout))
		assert.NoError(signalbox.Wait(notifier.Working, timeout))
		assert.NoError(signalbox.Wait(notifier.Working, timeout))
		assert.NoError(signalbox.Wait(notifier.Stopping, timeout))
		assert.NoError(signalbox.Wait(notifier.Stopped, timeout))
	}

	// Test.
	go asserter()
	go asserter()
	go func() {
		states := []notifier.Status{
			notifier.Starting,
			notifier.Ready,
			notifier.Working,
			notifier.Working,
			notifier.Ready,
			notifier.Stopping,
			notifier.Stopped,
		}
		for _, s := range states {
			signalbox.Notify(s)
		}
	}()

	wg.Wait()
}

// TestSignalboxWaitDirectOK tests the direct setting of a high status
// also raises lower signals.
func TestSignalboxWaitDirectOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	timeout := 10 * time.Millisecond
	signalbox := notifier.NewSignalbox()

	// Test.
	go func() {
		time.Sleep(5 * time.Millisecond)
		signalbox.Notify(notifier.Stopped)
	}()

	assert.NoError(signalbox.Wait(notifier.Working, timeout))
}

// TestSignalboxWaitTimeout tests the timeout when waiting for
// a signal.
func TestSignalboxWaitTimeout(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	timeout := 10 * time.Millisecond
	signalbox := notifier.NewSignalbox()

	// Test.
	go func() {
		time.Sleep(5 * time.Millisecond)
		signalbox.Notify(notifier.Working)
	}()

	assert.NoError(signalbox.Wait(notifier.Ready, timeout))
	assert.NoError(signalbox.Wait(notifier.Working, timeout))
	assert.ErrorContains(signalbox.Wait(notifier.Stopped, timeout), "waiting signalbox for stopped: timeout")
}

// TestSignalboxDoneOK tests the correct notification and signalling.
func TestSignalboxDoneOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	timeout := 50 * time.Millisecond
	signalbox := notifier.NewSignalbox()

	// Test.
	go func() {
		time.Sleep(5 * time.Millisecond)
		signalbox.Notify(notifier.Stopping)
	}()

	select {
	case <-signalbox.Done(notifier.Working):
		assert.True(true)
	case <-time.After(timeout):
		assert.Fail("timeout")
	}
}

// TestBundle tests notification of multiple signalboxes via a bundle.
func TestBundle(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	timeout := 10 * time.Millisecond
	signalboxA := notifier.NewSignalbox()
	signalboxB := notifier.NewSignalbox()
	signalboxC := notifier.NewSignalbox()
	var wg sync.WaitGroup
	wg.Add(3)
	asserter := func(signalbox *notifier.Signalbox) {
		defer wg.Done()
		assert.NoError(signalbox.Wait(notifier.Stopped, timeout))
	}
	bundle := notifier.NewBundle()
	bundle.Add(signalboxA, signalboxB, signalboxC)

	// Test.
	go asserter(signalboxA)
	go asserter(signalboxB)
	go asserter(signalboxC)

	bundle.Notify(notifier.Stopped)

	wg.Wait()

	assert.NoError(bundle.Wait(notifier.Stopped, timeout))
}

// EOF
