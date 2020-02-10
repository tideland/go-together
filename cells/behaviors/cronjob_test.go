// Tideland Go Together - Cells - Behaviors - Unit Tests
//
// Copyright (C) 2010-2020 Frank Mueller / Tideland / Oldenburg / Germany
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
	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestCronjobBehavior tests the cronjob behavior.
func TestCronjobBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	cronjob := func(emitter mesh.Emitter) {
		assert.NoError(emitter.Broadcast(event.New("action")))
	}
	plant := mesh.NewTestPlant(assert, behaviors.NewCronjobBehavior("cb", 100*time.Millisecond, cronjob), 1)
	defer plant.Stop()

	time.Sleep(525 * time.Millisecond)

	plant.AssertLength(0, 5)
	plant.AssertAll(0, func(evt *event.Event) bool {
		return evt.Topic() == "action"
	})
}

// EOF
