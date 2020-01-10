// Tideland Go Together - Limiter - Unit Tests
//
// Copyright (C) 2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package limiter_test

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/limiter"
)

//--------------------
// TESTS
//--------------------

// TestLimitOK tests the limiting of a number of function calls
// in multiple goroutines without an error.
func TestLimitOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	runs := 100
	startedC := make(chan struct{}, 1)
	stoppedC := make(chan struct{}, 1)
	resultC := make(chan int, 1)
	go func() {
		max := 0
		act := 0
		count := 0
		for count < runs {
			select {
			case <-startedC:
				act++
				if act > max {
					max = act
				}
			case <-stoppedC:
				act--
			}
			count++
		}
		resultC <- max
	}()
	job := func() error {
		startedC <- struct{}{}
		time.Sleep(50 * time.Millisecond)
		stoppedC <- struct{}{}
		return nil
	}
	l := limiter.New(10)
	ctx := context.Background()

	// Test.
	for i := 0; i < runs; i++ {
		go func() {
			err := l.Do(ctx, job)
			assert.NoError(err)
		}()
	}

	assert.True(<-resultC <= 11)
}

// TestLimitError tests the returning of en error by an
// executed function.
func TestLimitError(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	job := func() error {
		time.Sleep(25 * time.Millisecond)
		return errors.New("ouch")
	}
	l := limiter.New(5)
	ctx := context.Background()

	var wg sync.WaitGroup
	wg.Add(25)

	// Test.
	for i := 0; i < 25; i++ {
		go func() {
			err := l.Do(ctx, job)
			assert.ErrorMatch(err, "ouch")
			wg.Done()
		}()
	}

	wg.Wait()
}

// EOF
