// Tideland Go Together - Cells - Behaviors - Unit Tests
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
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/generators"

	"tideland.dev/go/together/cells/behaviors"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestCollectorBehavior tests the collector behavior.
func TestCollectorBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	processor := func(r mesh.EventSinkReader) (mesh.Event, error) {
		l := r.Len()
		return mesh.NewEvent("length", l)
	}
	behavior := behaviors.NewCollectorBehavior(10, processor)
	tester := func(evt mesh.Event) bool {
		switch evt.Topic() {
		case "length":
			var l int
			err := evt.Payload(&l)
			assert.NoError(err)
			assert.Equal(l, 10)
			return false
		case behaviors.TopicResetted:
			return true
		}
		return false
	}
	tb := mesh.NewTestbed(behavior, tester)

	// Run the tests and check length and resetting.
	for _, topic := range generator.Words(25) {
		tb.Emit(topic)
	}

	tb.Emit(behaviors.TopicProcess)
	tb.Emit(behaviors.TopicReset)

	err := tb.Wait(time.Second)
	assert.NoError(err)
}

// EOF
