// Tideland Go Together - Loop - Examples
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
	"fmt"
	"time"

	"tideland.dev/go/together/loop"
)

//--------------------
// EXAMPLES
//--------------------

// ExampleSimple shows the minimum usage of loop without an external context
// and no error.
func ExampleSimple() {
	done := make(chan struct{})
	// A loop worker is just any function or method getting a loop.Terminator
	// as argument and returning an error. It runs the business logic.
	worker := func(lt loop.Terminator) error {
		// Inside the worker you run a typical for-select-loop.
		for {
			select {
			case <-lt.Done():
				// loop.Terminator.Done() tells the loop to stop working.
				// In case of an internal error it can be returned here.
				close(done)
				return nil
			case <-time.Tick(time.Minute):
				// Just for linter.
			}
		}
	}
	// loop.Go() runs the worker function concurrently.
	l, err := loop.Go(worker)
	if err != nil {
		fmt.Printf("error in loop.Go: %v\n", err)
		return
	}
	// loop.Loop.Stop() tells the worker function via loop.Terminator
	// to stop working.
	err = l.Stop()
	if err != nil {
		fmt.Printf("error returned by loop.Stop: %v\n", err)
		return
	}
	// So now done has been closed.
	<-done
}

// ExampleError shows running a goroutine returning an error immediately.
func ExampleError() {
	done := make(chan struct{})
	worker := func(lt loop.Terminator) error {
		defer close(done)
		return fmt.Errorf("ouch")
	}
	// loop.Go() runs the worker function concurrently.
	l, err := loop.Go(worker)
	if err != nil {
		fmt.Printf("error in loop.Go: %v\n", err)
		// Output: error in loop.Go: ouch
		return
	}
	// loop.Loop.Stop() tells the worker function via loop.Terminator
	// to stop working.
	err = l.Stop()
	if err != nil {
		fmt.Printf("error returned by loop.Stop: %v\n", err)
		// Output: error returned by loop.Stop: ouch
		return
	}
	// So now done has been closed.
	<-done
}

// EOF
