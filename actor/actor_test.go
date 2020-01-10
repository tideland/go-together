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
	"tideland.dev/go/together/notifier"
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
	signalbox := notifier.NewSignalbox()
	act, err := actor.Go(actor.WithSignalbox(signalbox))
	assert.NoError(err)
	assert.NotNil(act)

	err = signalbox.Wait(notifier.Working, timeout)
	assert.NoError(err)
	assert.NoError(act.Stop(nil))
	assert.NoError(act.Err())
}

// TestPureError is simply starting and stopping an Actor.
// Returning the stop error.
func TestPureGo(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	signalbox := notifier.NewSignalbox()
	act, err := actor.Go(actor.WithSignalbox(signalbox))
	assert.NoError(err)
	assert.NotNil(act)

	err = signalbox.Wait(notifier.Working, timeout)
	assert.NoError(err)
	assert.ErrorMatch(act.Stop(errors.New("damn")), "damn")
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
	defer act.Stop(nil)

	counter := 0

	for i := 0; i < 5; i++ {
		err := act.DoSync(func() error {
			counter++
			return nil
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
	defer act.Stop(nil)

	// Scenario: Timeout is shorter than needed time, so error
	// is returned.
	err = act.DoSyncTimeout(func() error {
		time.Sleep(time.Second)
		return nil
	}, 500*time.Millisecond)

	assert.ErrorMatch(err, ".*timed out.*")
}

// TestAsyncWithQueueCap tests running multiple calls asynchronously.
func TestAsyncWithQueueCap(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	act, err := actor.Go(actor.WithQueueCap(100))
	assert.NoError(err)
	defer act.Stop(nil)

	assert.Equal(act.QueueCap(), 100)

	sigC := make(chan bool, 1)
	doneC := make(chan bool, 1)

	// Start background func waiting for the signals of
	// the asynchrounous calls.
	go func() {
		count := 0
		for range sigC {
			count++
			if count == 100 {
				break
			}
		}
		doneC <- true
	}()

	// Now start asynchrounous calls.
	start := time.Now()
	for i := 0; i < 100; i++ {
		act.DoAsync(func() error {
			time.Sleep(5 * time.Millisecond)
			sigC <- true
			return nil
		})
	}
	enqueued := time.Since(start)

	// Expect signal done to be sent about one second later.
	<-doneC
	done := time.Since(start)

	assert.True((done - 500*time.Millisecond) > enqueued)
}

// TestRecovery tests handling panics.
func TestRecovery(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	counter := 0
	recovered := false
	doneC := make(chan struct{})
	recoverer := func(reason interface{}) error {
		recovered = true
		close(doneC)
		return nil
	}
	act, err := actor.Go(actor.WithRecoverer(recoverer))
	assert.NoError(err)
	defer act.Stop(nil)

	err = act.DoSyncTimeout(func() error {
		counter++
		// Will crash on first call.
		print(counter / (counter - 1))
		return nil
	}, time.Second)
	assert.ErrorMatch(err, ".*timed out.*")
	<-doneC
	assert.True(recovered)
	err = act.DoSync(func() error {
		counter++
		return nil
	})
	assert.NoError(err)
	assert.Equal(counter, 2)
}

// EOF
