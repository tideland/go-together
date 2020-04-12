// Tideland Go Together - Fuse - Unit Tests
//
// Copyright (C) 2017-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package fuse_test

//--------------------
// IMPORTS
//--------------------

import (
	"sync"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/fuse"
)

//--------------------
// TESTS
//--------------------

// TestSignalWaitOK tests the correct notification and waiting.
func TestSignalWaitOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	timeout := 10 * time.Millisecond
	signal := fuse.NewSignal()
	var wg sync.WaitGroup
	wg.Add(2)
	asserter := func() {
		defer wg.Done()
		assert.OK(signal.Wait(fuse.Starting, timeout))
		assert.OK(signal.Wait(fuse.Starting, timeout))
		assert.OK(signal.Wait(fuse.Ready, timeout))
		assert.OK(signal.Wait(fuse.Working, timeout))
		assert.OK(signal.Wait(fuse.Working, timeout))
		assert.OK(signal.Wait(fuse.Stopping, timeout))
		assert.OK(signal.Wait(fuse.Stopped, timeout))
	}

	// Test.
	go asserter()
	go asserter()
	go func() {
		states := []fuse.Status{
			fuse.Starting,
			fuse.Ready,
			fuse.Working,
			fuse.Working,
			fuse.Ready,
			fuse.Stopping,
			fuse.Stopped,
		}
		for _, s := range states {
			signal.Notify(s)
		}
	}()

	wg.Wait()
}

// TestSignalWaitDirectOK tests the direct setting of a high status
// also raises lower signals.
func TestSignalWaitDirectOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	timeout := 10 * time.Millisecond
	signal := fuse.NewSignal()

	// Test.
	go func() {
		time.Sleep(5 * time.Millisecond)
		signal.Notify(fuse.Stopped)
	}()

	assert.OK(signal.Wait(fuse.Working, timeout))
}

// TestSignalWaitTimeout tests the timeout when waiting for
// a signal.
func TestSignalWaitTimeout(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	timeout := 10 * time.Millisecond
	signal := fuse.NewSignal()

	// Test.
	go func() {
		time.Sleep(5 * time.Millisecond)
		signal.Notify(fuse.Working)
	}()

	assert.OK(signal.Wait(fuse.Ready, timeout))
	assert.OK(signal.Wait(fuse.Working, timeout))
	assert.ErrorContains(signal.Wait(fuse.Stopped, timeout), "waiting signal for stopped: timeout")
}

// TestSignalDoneOK tests the correct notification and signalling.
func TestSignalDoneOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	timeout := 50 * time.Millisecond
	signal := fuse.NewSignal()

	// Test.
	go func() {
		time.Sleep(5 * time.Millisecond)
		signal.Notify(fuse.Stopping)
	}()

	select {
	case <-signal.Done(fuse.Working):
		assert.True(true)
	case <-time.After(timeout):
		assert.Fail("timeout")
	}
}

// EOF
