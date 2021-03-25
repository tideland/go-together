// Tideland Go Together - Cells - Behaviors
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors_test // import "tideland.dev/go/together/cells/behaviors"

//--------------------
// IMPORTS
//--------------------

import (
	"strconv"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"

	"tideland.dev/go/together/cells/behaviors"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestCallbackBehavior tests the callback behavior.
func TestCallbackBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	count := 50
	countA := 0
	countB := 0
	callbackA := func(evt mesh.Event, out mesh.Emitter) error {
		return out.Emit("a")
	}
	callbackB := func(evt mesh.Event, out mesh.Emitter) error {
		return out.Emit("b")
	}
	behavior := behaviors.NewCallbackBehavior(callbackA, callbackB)
	tester := func(evt mesh.Event) bool {
		switch evt.Topic() {
		case "a":
			countA++
		case "b":
			countB++
		}
		return countA == countB && countA == count
	}
	tb := mesh.NewTestbed(behavior, tester)

	// Run the tests and check if the number of calls
	// of both callbacks.
	for i := 0; i < count; i++ {
		topic := strconv.Itoa(i)
		tb.Emit(topic)
	}

	err := tb.Wait(time.Second)
	assert.NoError(err)
}

// EOF
