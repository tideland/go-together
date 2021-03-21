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
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"

	"tideland.dev/go/together/cells/behaviors"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestBroadcasterBehavior tests the broadcaster behavior. It simple receives some
// events and checks if those have been emitted again.
func TestBroadcasterBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	behavior := behaviors.NewBroadcasterBehavior()
	topics := make(map[string]bool)
	tester := func(evt *mesh.Event) bool {
		if evt.Topic() == "done" {
			return true
		}
		topics[evt.Topic()] = true
		return false
	}
	tb := mesh.NewTestbed(behavior, tester)

	tb.Emit("one")
	tb.Emit("two")
	tb.Emit("three")
	tb.Emit("done")

	err := tb.Wait(time.Second)
	assert.NoError(err)

	assert.True(topics["one"])
	assert.True(topics["two"])
	assert.True(topics["three"])
	assert.False(topics["done"])
}

// EOF
