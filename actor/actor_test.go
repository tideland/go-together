// Tideland Go Together - Actor - Unit Tests
//
// Copyright (C) 2017-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package actor_test

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"errors"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/actor"
)

//--------------------
// CONSTANTS
//--------------------

const (
	timeout = 5 * time.Second
)

//--------------------
// TESTS
//--------------------

// TestPureOK is simply starting and stopping an Actor.
func TestPureGoOK(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	done := false
	act, err := actor.Go(actor.WithFinalizer(func(err error) error {
		done = true
		return err
	}))
	assert.NoError(err)
	assert.NotNil(act)

	assert.NoError(act.Stop())
	assert.NoError(act.Err())
	assert.True(done)
}

// TestPureError is simply starting and stopping an Actor.
// Returning the stop error.
func TestPureError(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	act, err := actor.Go(actor.WithFinalizer(func(err error) error {
		return errors.New("damn")
	}))
	assert.NoError(err)
	assert.NotNil(act)

	assert.ErrorMatch(act.Stop(), "damn")
	assert.ErrorMatch(act.Err(), "damn")
}

// TestWithContext is simply starting and stopping an Actor
// with a context.
func TestWithContext(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	ctx, cancel := context.WithCancel(context.Background())
	act, err := actor.Go(actor.WithContext(ctx))
	assert.NoError(err)
	assert.NotNil(act)

	cancel()
	assert.NoError(act.Err())
}

// TestSync tests synchronous calls.
func TestSync(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	act, err := actor.Go()
	assert.NoError(err)
	defer act.Stop()

	counter := 0

	for i := 0; i < 5; i++ {
		err := act.DoSync(func() {
			counter++
		})
		assert.Nil(err)
	}

	assert.Equal(counter, 5)
}

// TestTimeout tests timout error of a synchronous Action.
func TestTimeout(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	act, err := actor.Go()
	assert.NoError(err)
	defer act.Stop()

	// Scenario: Timeout is shorter than needed time, so error
	// is returned.
	err = act.DoSyncTimeout(func() {
		time.Sleep(time.Second)
	}, 500*time.Millisecond)

	assert.ErrorMatch(err, ".*timed out.*")
}

// TestAsyncWithQueueCap tests running multiple calls asynchronously.
func TestAsyncWithQueueCap(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	act, err := actor.Go(actor.WithQueueCap(128))
	assert.NoError(err)
	defer act.Stop()

	sigs := make(chan struct{}, 1)
	done := make(chan struct{}, 1)

	// Start background func waiting for the signals of
	// the asynchrounous calls.
	go func() {
		count := 0
		for range sigs {
			count++
			if count == 128 {
				break
			}
		}
		close(done)
	}()

	// Now start asynchrounous calls.
	start := time.Now()
	for i := 0; i < 128; i++ {
		act.DoAsync(func() {
			time.Sleep(5 * time.Millisecond)
			sigs <- struct{}{}
		})
	}
	enqueued := time.Since(start)

	// Expect signal done to be sent about one second later.
	<-done
	duration := time.Since(start)

	assert.True((duration - 640*time.Millisecond) > enqueued)
}

// TestRecoveryOK tests handling panics successfully.
func TestRecoveryOK(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	counter := 0
	recovered := false
	done := make(chan struct{})
	recoverer := func(reason interface{}) error {
		recovered = true
		close(done)
		return nil
	}
	act, err := actor.Go(actor.WithRecoverer(recoverer))
	assert.NoError(err)
	defer act.Stop()

	err = act.DoSyncTimeout(func() {
		counter++
		// Will crash on first call.
		print(counter / (counter - 1))
	}, time.Second)
	assert.ErrorMatch(err, ".*timed out.*")
	<-done
	assert.True(recovered)
	err = act.DoSync(func() {
		counter++
	})
	assert.NoError(err)
	assert.Equal(counter, 2)
}

// TestRecoveryError tests handling panics with error.
func TestRecoveryError(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	counter := 0
	recovered := false
	done := make(chan struct{})
	recoverer := func(reason interface{}) error {
		recovered = true
		close(done)
		return errors.New("ouch")
	}
	act, err := actor.Go(actor.WithRecoverer(recoverer))
	assert.NoError(err)
	defer act.Stop()

	err = act.DoSyncTimeout(func() {
		counter++
		// Will crash on first call.
		print(counter / (counter - 1))
	}, time.Second)
	assert.ErrorMatch(err, "ouch")
	<-done
	assert.True(recovered)
	err = act.DoSync(func() {
		counter++
	})
	assert.ErrorMatch(err, "ouch")
}

// EOF
