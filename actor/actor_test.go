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
// TESTS
//--------------------

// TestPureOK is simply starting and stopping an Actor.
func TestPureOK(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	done := false
	act, err := actor.Go(actor.WithFinalizer(func(err error) error {
		done = true
		return err
	}))
	assert.OK(err)
	assert.NotNil(act)

	assert.OK(act.Stop())
	assert.OK(act.Err())
	assert.OK(done)
}

// TestPureKill is simply starting and killing an Actor.
func TestPureKill(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	done := false
	act, err := actor.Go(actor.WithFinalizer(func(err error) error {
		done = true
		return err
	}))
	assert.OK(err)
	assert.NotNil(act)

	assert.ErrorMatch(act.Kill(errors.New("killed")), "killed")
	assert.ErrorMatch(act.Err(), "killed")
	assert.OK(done)
}

// TestPureError is simply starting and stopping an Actor.
// Returning the stop error.
func TestPureError(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	act, err := actor.Go(actor.WithFinalizer(func(err error) error {
		return errors.New("damn")
	}))
	assert.OK(err)
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
	assert.OK(err)
	assert.NotNil(act)

	cancel()
	assert.OK(act.Err())
}

// TestSync tests synchronous calls.
func TestSync(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	act, err := actor.Go()
	assert.OK(err)

	counter := 0

	for i := 0; i < 5; i++ {
		assert.OK(act.DoSync(func() {
			counter++
		}))
	}

	assert.Equal(counter, 5)
	assert.OK(act.Stop())

	assert.ErrorMatch(act.DoSync(func() {
		counter++
	}), ".*timeout.*")
}

// TestTimeout tests timout error of a synchronous Action.
func TestTimeout(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	act, err := actor.Go()
	assert.OK(err)

	// Scenario: Timeout is shorter than needed time, so error
	// is returned.
	err = act.DoSyncTimeout(func() {
		time.Sleep(5 * time.Second)
	}, 500*time.Millisecond)

	assert.ErrorMatch(err, ".*timeout.*")
	assert.OK(act.Stop())
}

// TestAsyncWithQueueCap tests running multiple calls asynchronously.
func TestAsyncWithQueueCap(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	act, err := actor.Go(actor.WithQueueCap(128))
	assert.OK(err)

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
		assert.OK(act.DoAsync(func() {
			time.Sleep(5 * time.Millisecond)
			sigs <- struct{}{}
		}))
	}
	enqueued := time.Since(start)

	// Expect signal done to be sent about one second later.
	<-done
	duration := time.Since(start)

	assert.OK((duration - 640*time.Millisecond) > enqueued)
	assert.OK(act.Stop())
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
	assert.OK(err)

	err = act.DoSyncTimeout(func() {
		counter++
		// Will crash on first call.
		print(counter / (counter - 1))
	}, time.Second)
	assert.ErrorMatch(err, ".*timeout.*")
	<-done
	assert.OK(recovered)
	err = act.DoSync(func() {
		counter++
	})
	assert.OK(err)
	assert.Equal(counter, 2)
	assert.OK(act.Stop())
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
	assert.OK(err)

	err = act.DoSyncTimeout(func() {
		counter++
		// Will crash on first call.
		print(counter / (counter - 1))
	}, time.Second)
	assert.ErrorMatch(err, "ouch")
	<-done
	assert.OK(recovered)
	assert.ErrorMatch(act.DoSync(func() {
		counter++
	}), "ouch")
	assert.ErrorMatch(act.Stop(), "ouch")
}

// EOF
